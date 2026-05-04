package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/op/go-logging"
)

type config struct {
	// Global config
	optDebug         bool // Global debug flag
	optLogRawJSON    bool // Log the raw JSON in debug even if it is fine
	optVerifySSL     bool
	optProbe         bool // Probe mode
	durBetweenPolls  time.Duration
	durBetweenProbes time.Duration
	nos429retries    int
	durAfter429      time.Duration
	flagTarget       string
	flagUsername     string
	flagPassword     string
	flagMetricPort   int
	flagProbeFile    string
	flagMetricPrefix string
	// global objs
	log *logging.Logger
}

func printVersion() {
	fmt.Printf("unaspoller %s\n", version)
}

// Process command line flags
func processFlags(c *config) {
	var flagDurBetweenPolls, flagDurBetweenProbes, flagDurAfter429 string
	var optVersion bool

	flag.BoolVar(&optVersion, "version", false, "print version information")
	flag.StringVar(&c.flagUsername, "username", "unifipoller", "username to connect to UNAS device")
	flag.StringVar(&c.flagPassword, "password", "", "password to connect to UNAS device")
	flag.StringVar(&c.flagTarget, "target", "192.168.1.1", "IP address of UNAS device")
	flag.BoolVar(&c.optProbe, "probe", false, "probe known URLs once dumping info and stopping")
	flag.StringVar(&c.flagProbeFile, "probefile", "unaspoller-probe.txt", "Filename for -probe mode output")
	flag.BoolVar(&c.optDebug, "debug", false, "enable debug mode")
	flag.BoolVar(&c.optVerifySSL, "verifyssl", true, "verify SSL certificates")
	flag.BoolVar(&c.optLogRawJSON, "lograwjson", false, "log the raw JSON received even if it is good")
	flag.StringVar(&flagDurBetweenPolls, "pollinterval", "15s", "Sleep duration between polls")
	flag.StringVar(&flagDurBetweenProbes, "probeinterval", "2s", "Sleep duration between polls in probe mode")
	flag.IntVar(&c.flagMetricPort, "metricport", 8090, "Port to listen on for /metrics endpoint")
	flag.IntVar(&c.nos429retries, "max429retries", 5, "Maximum number of retries if 429 received")
	flag.StringVar(&flagDurAfter429, "retrydelay", "2s", "Time to wait until retry if 429 received")
	flag.StringVar(&c.flagMetricPrefix, "metrixprefix", "unas", "Prefix for metrics")
	// TODO - metrics host

	flag.Parse()

	if tDur, err := time.ParseDuration(flagDurBetweenPolls); err == nil {
		c.durBetweenPolls = tDur
	} else {
		fmt.Printf("ERROR: failed to parse -pollinterval value of [%s]: %s\n", flagDurBetweenPolls, err)
		os.Exit(1)
	}

	if tDur, err := time.ParseDuration(flagDurBetweenProbes); err == nil {
		c.durBetweenProbes = tDur
	} else {
		fmt.Printf("ERROR: failed to parse -probeinterval value of [%s]: %s\n", flagDurBetweenProbes, err)
		os.Exit(1)
	}

	if tDur, err := time.ParseDuration(flagDurAfter429); err == nil {
		c.durAfter429 = tDur
	} else {
		fmt.Printf("ERROR: failed to parse -retrydelay value of [%s]: %s\n", flagDurAfter429, err)
		os.Exit(2)
	}

	if optVersion {
		printVersion()
		os.Exit(0)
	}
}
