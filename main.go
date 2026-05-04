package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	version = "v0.0.1"
)

func (u *UNAS) mainPollLoop() error {
	for true {
		u.c.log.Debugf("Starting poll")

		pollURLs := []string{
			"/proxy/drive/api/v2/storage",
			"/proxy/drive/api/v2/systems/device-info",
		}
		for _, url := range pollURLs {
			u.c.log.Debugf("Polling %s\n", url)
			if err := u.doDriveAPIDef(url); err != nil {
				fmt.Printf("Blahgle: %s\n", err)
			}
		}

		u.c.log.Debugf("Sleeping for %s between polls", u.c.durBetweenPolls)
		time.Sleep(u.c.durBetweenPolls)
	}
	return nil
}

func main() {
	conf := &config{}
	processFlags(conf)

	setupLogging(conf)
	conf.log.Infof("unaspoller version %s", version)

	// Setup object with login details from flags/env
	u := NewUNAS(conf.flagTarget, conf)
	u.SetUsername(u.c.flagUsername)
	// Set password, command line overrules env var
	pw := os.Getenv("UNAS_PW")
	if u.c.flagPassword != "" {
		// FIXME - this means you can't set a blank password
		// But it is unlikely that will be the case
		pw = u.c.flagPassword
	}
	u.SetPassword(pw)
	u.c.flagPassword = "" // Scrub this

	if err := u.registerAPIDefs(); err != nil {
		u.c.log.Errorf("Failed to register all API definitions: %w", err)
		return
	}

	// Register metrics and start HTTP server
	reg := prometheus.NewRegistry()
	u.m = u.newMetrics(reg)
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	// FIXME - hackity hack
	go func() {
		http.ListenAndServe(fmt.Sprintf(":%d", u.c.flagMetricPort), nil)
	}()

	// Now we can login
	err := u.LoginUNAS()
	if err != nil {
		u.c.log.Errorf("Failed to login: %w", err)
		return
	}
	u.c.log.Debugf("Logged in as '%s'", u.username)

	if u.c.optProbe {
		// Do a single shot of a probe of the device and then exit
		u.Probe(u.c.flagProbeFile)
		u.c.log.Info("Exiting after probe mode")
		return
	}

	// Main loop
	if err := u.mainPollLoop(); err != nil {
		u.c.log.Errorf("main poll loop error: %s\n", err)
	}
	u.c.log.Info("Exiting")
}
