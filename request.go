package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var err401 = fmt.Errorf("401")

func (u *UNAS) LoginUNAS() error {
	params := fmt.Sprintf(`{"username":"%s","password":"%s"}`, u.username, u.password)
	var req *http.Request
	var err error

	apiPath := "/api/auth/login"

	req, err = http.NewRequest(http.MethodPost, u.baseURL+apiPath, bytes.NewBufferString(params))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	dur429loginSleep := time.Second * 60

	for attempt := range u.c.nos429retries {
		// Try to do the request...
		u.c.log.Debugf("Attempt %d/%d to login", attempt+1, u.c.nos429retries)
		resp, err := u.doRequestOnce(req)
		if err != nil {
			return fmt.Errorf("making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests {
			u.c.log.Debugf("Got 429 when trying to login, sleeping for %s", dur429loginSleep)

			// check for `Retry-After` header...
			// It doesn't return one at the moment, but if it does we should probably
			// add code to honour it
			if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
				// TODO - just write the code to handle it
				u.c.log.Debugf("Got 429 during login and Retry-After=[%s]", retryAfter)
			}
			time.Sleep(dur429loginSleep)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("login failure: target=%s user:=%s: (status: %s)", u.target, u.username, resp.Status)
		}

		// Successful login
		u.c.log.Infof("login successful")

		// FIXME - use a cookie jar
		cookies := resp.Cookies()
		u.cookies = make([]*http.Cookie, len(cookies))
		copy(u.cookies, cookies)

		if rCsrf := resp.Header.Get("x-csrf-token"); rCsrf != "" {
			u.csrf = rCsrf
		}

		if rCsrf := resp.Header.Get("x-updated-csrf-token"); rCsrf != "" {
			u.csrf = rCsrf
		}
		return nil
	}
	// If we get here we've tried too many times and got 429s, just fail
	u.c.log.Errorf("too many 429s. Failed to login")
	return fmt.Errorf("login failure: too many 429s")
}

func (u *UNAS) doRequestOnce(req *http.Request) (*http.Response, error) {
	// TODO - logging?
	resp, err := u.client.Do(req)
	return resp, err
}

// Perform an HTTP request using the cookies/headers
// need to handle 401 (StatusUnauthorized) with a relogin
//   - UNAS TOKEN cookie just expires after a short period (2h-ish?)
//
// also handles 429 with retries
// - Handle various status codes and 429/retries
//
// We don't use this function for login in case we end up in a 401/429 endless loop
func (u *UNAS) doRequest(req *http.Request) ([]byte, error) {
	// Add in the headers and cookies
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("X-CSRF-Token", u.csrf)
	for _, c := range u.cookies {
		req.AddCookie(c)
	}

	for attempt := range u.c.nos429retries {
		// Try to do the request...
		u.c.log.Debugf("Attempt %d/%d to perform %s %s", attempt+1, u.c.nos429retries, req.Method, req.URL)

		resp, err := u.doRequestOnce(req)
		defer resp.Body.Close()
		if err != nil {
			return []byte{}, fmt.Errorf("doRequestOnce: %w", err)
		}

		// Happy path
		if resp.StatusCode == http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return body, fmt.Errorf("reading response: %w", err)
			}
			return body, nil
		}

		// Check for 429
		if resp.StatusCode == http.StatusTooManyRequests {
			// TODO - check what 429 behaviour is
			//	n requests in last 30 seconds
			// 	or bucketed so that it releases randomly...
			// TODO - test once own device has finished repairing pool
			u.c.log.Debugf("Got 429 when doing url=%s sleeping for %s", req.URL, u.c.durAfter429)

			// check for `Retry-After` header...
			// It doesn't return one at the moment, but if it does we should probably
			// add code to honour it
			if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
				// TODO - just write the code to handle it
				u.c.log.Debugf("Got 429 and Retry-After=[%s]", retryAfter)
			}
			time.Sleep(u.c.durAfter429)
			continue
		}
		if resp.StatusCode == http.StatusUnauthorized {
			u.c.log.Debug("Got 401. Will bail and attempt to login again")
			// Bail with appropriate error and retry once in calling func `doGetRequest()`
			return []byte{}, err401
		}

		// Some other status code here, we just fail without trying to retry
		u.c.log.Debugf("Got some unexpected statusCode (%d). Failed to perform request", resp.StatusCode)
		return []byte{}, fmt.Errorf("request failed (statusCode=%d)", resp.StatusCode)
	}

	u.c.log.Debugf("Too many 429s. Failed to perform request")
	return []byte{}, fmt.Errorf("request failed too many 429s")
}

// Perform an HTTP GET using the cookies/headers
func (u *UNAS) doGetRequest(url string) ([]byte, error) {
	var req *http.Request
	var err error

	req, err = http.NewRequest(http.MethodGet, u.baseURL+url, nil)
	if err != nil {
		return []byte{}, fmt.Errorf("creating request: %w", err)
	}

	body, err := u.doRequest(req)
	if err == nil {
		return body, err
	} else if errors.Is(err, err401) {
		// We got a 401, login and try again
		err := u.LoginUNAS()
		if err != nil {
			u.c.log.Errorf("got a 401 so retried login but that failed: %w", err)
			return []byte{}, err
		}
		// We logged in fine, retry the request
		body, err = u.doRequest(req)
		if err == nil {
			return body, err
		}
	}
	// Not sure what the error could be - just return it
	return []byte{}, err
}
