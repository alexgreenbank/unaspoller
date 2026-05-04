package main

import (
	"os"

	"github.com/op/go-logging"
)

func setupLogging(c *config) {
	c.log = logging.MustGetLogger("unaspoller")
	logBackend := logging.NewLogBackend(os.Stdout, "", 0)

	format := logging.MustStringFormatter(`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`)

	// Set the formatter
	logging.SetFormatter(format)

	// Set the backend to be used.
	logging.SetBackend(logBackend)

	if c.optDebug {
		logging.SetLevel(logging.DEBUG, "")
	} else {
		logging.SetLevel(logging.INFO, "")
	}
}
