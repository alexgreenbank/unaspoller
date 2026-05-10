package main

import (
	"fmt"
	"os"
	"time"
)

func (u *UNAS) doProbeKnownURLs(fname string) (error, int, int, int, int) {
	knownURLs := []string{
		"/proxy/drive/api/v2/systems/device-info", // TODO - move to types.go ...
		"/proxy/drive/api/v1/systems/performance/file-operations",
		"/proxy/drive/api/v1/systems/storage?type=detail",
		"/proxy/drive/api/v2/systems/disk-stats",
		"/proxy/drive/api/v2/systems/network-io",
		"/proxy/users/drive/api/v1/systems/identity",
		"/proxy/users/drive/api/v1/systems/info",
		"/proxy/users/drive/api/v2/drives",
		"/proxy/users/drive/api/v2/groups",
		"/proxy/users/drive/api/v2/storage",
	}

	u.c.log.Infof("performing probe of %d known URLs and then exiting", len(knownURLs))

	// Open file to write to
	tmpFile, err := os.Create(fname)
	if err != nil {
		u.c.log.Errorf("Can't open [%s] to write probe output to, using stdout: %s", err)
	} else {
		u.c.log.Infof("Output will be saved in [%s]", fname)
		u.probeFile = tmpFile
		fmt.Fprintf(u.probeFile, "PROBE:version=[%s]\n", version)
		defer u.probeFile.Close()
	}

	// Keep track of how many problems we had
	nosOK := 0
	nosHTTPError := 0
	nosUnmarshalError := 0
	nosValidationError := 0

	for i, baseURL := range knownURLs {
		if i > 0 {
			u.c.log.Debugf("Sleeping for %s between polls", u.c.durBetweenProbes)
			time.Sleep(u.c.durBetweenProbes)
		}

		// We may need to append or otherwise modify the URL
		probeURL := baseURL

		if baseURL == "/proxy/drive/api/v2/systems/disk-stats" {
			// Add on a param for
			// ?start=1777565467&end=1777652767&interval=900
			now := time.Now().Unix()
			probeURL += fmt.Sprintf("?start=%d&end=%d&interval=900", now-87300, now)
		}

		u.c.log.Debugf("Probing URL %d of %d", i+1, len(knownURLs))
		body, err := u.doGetRequest(probeURL)
		if err != nil {
			nosHTTPError++
			u.c.log.Errorf("doRequest: %w", err)
			// Don't add the body here
			if u.probeFile != nil {
				fmt.Fprintf(u.probeFile, "PROBE:ERROR:[%s]: err=[%s]\n", probeURL, err)
			} else {
				fmt.Printf("PROBE:ERROR:[%s]: err=[%s]\n", probeURL, err)
			}
		} else {
			u.c.log.Debugf("Probe of URL successful")
			nosOK++
			// do something with body
			// TODO - validate JSON (i.e. it parses) as expected (if known)
			// TODO - validate individual values (if known)
			// scrub possible sensitive info
			scrubbed := scrub(baseURL, body)
			// Just print it for now
			if u.probeFile != nil {
				fmt.Fprintf(u.probeFile, "PROBE:[%s]: resp=[", probeURL)
				u.probeFile.Write(scrubbed)
				u.probeFile.WriteString("]\n")
			} else {
				fmt.Printf("PROBE:[%s]: resp=[%s]\n", probeURL, string(scrubbed))
			}
		}
	}
	fmt.Fprintf(u.probeFile, "PROBE:DONE\n")
	return nil, nosOK, nosHTTPError, nosUnmarshalError, nosValidationError
}

func (u *UNAS) Probe(fname string) {
	// Check for single shot probe mode
	err, nosOK, nosHTTP, nosJSON, nosValidation := u.doProbeKnownURLs(u.c.flagProbeFile)
	if err != nil {
		u.c.log.Errorf("probe mode: %w", err)
	}
	u.c.log.Debugf("Probe mode summary. OK=%d errHTTP=%d errJSON=%d errValidation=%d", nosOK, nosHTTP, nosJSON, nosValidation)
	if nosHTTP+nosJSON+nosValidation > 0 {
		if u.probeFile != nil {
			u.c.log.Info("Probe mode finished. There were some problems, but this is OK as the project needs more example data")
			u.c.log.Infof("Review the data in file [%s] and scrub any data given possible security concerns", u.c.flagProbeFile)
		} else {
			u.c.log.Info("Probe mode finished. There were some problems, but this is OK as the project needs more example data")
			u.c.log.Info("Review this data and scrub any data given possible security concerns")
		}
		u.c.log.Info("Please read the details at https://github.com/alexgreenbank/unaspoller/CONTRIBUTING.md")
		u.c.log.Info("and consider uploading your scrubbed data to help unaspoller better!")
		u.c.log.Info("Thank you!")
	} else {
		u.c.log.Info("Probe mode finished. Everything looks good!")
		u.c.log.Info("Your data may still be useful so please consider reading https://github.com/alexgreenbank/unaspoller/CONTRIBUTING.md")
	}
}
