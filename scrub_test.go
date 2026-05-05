package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/require"
)

func TestScrubKeyString(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		key      string
		replace  string
		expected string
	}{
		{"empty", "", "", "", ""},
		{"no match", "{\"simplejson\":123.0}", "address", "foo", "{\"simplejson\":123.0}"},
		{"partial match", "{\"simplejson\":123.0}", "json", "foo", "{\"simplejson\":123.0}"},
		{"single match", "{\"address\":\"192.168.1.1\"}", "address", "1.1.1.1", "{\"address\":\"1.1.1.1\"}"},
		{"truncated", "{\"address\":\"192.168.1.1", "address", "1.1.1.1", "{\"address\":\"192.168.1.1"},
		{
			"multiple match",
			"{\"objs\":[{\"name\":\"one\",\"address\":\"192.168.1.1\"},{\"name\":\"two\",\"address\":\"10.0.0.1\"}]}",
			"address",
			"1.1.1.1",
			"{\"objs\":[{\"name\":\"one\",\"address\":\"1.1.1.1\"},{\"name\":\"two\",\"address\":\"1.1.1.1\"}]}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := scrubKeyString(tt.data, tt.key, tt.replace)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestScrub(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		input    string
		expected string
	}{
		{"empty", "/unknown", "", ""},
		{
			"device-info",
			"/proxy/drive/api/v2/systems/device-info",
			"{\"usbs\":null,\"version\":\"4.1.16\",\"name\":\"secret\",\"mac\":\"de:ad:be:ef:01:02\",\"startupTime\":\"2026-04-24T22:20:27Z\"}",
			"{\"usbs\":null,\"version\":\"4.1.16\",\"name\":\"scrubbed\",\"mac\":\"01:23:45:67:89:ab\",\"startupTime\":\"2026-04-24T22:20:27Z\"}",
		},
		{
			"systems info",
			"/proxy/users/drive/api/v1/systems/info",
			"{\"id\":\"secretid\",\"name\":\"secretname\",\"idname\":\"this is fine\"}",
			"{\"id\":\"scrubbed\",\"name\":\"scrubbed\",\"idname\":\"this is fine\"}",
		},
		{
			"systems info",
			"/proxy/users/drive/api/v2/drives",
			"{\"id\":\"secretid\",\"name\":\"secretname\",\"idname\":\"this is fine\"}",
			"{\"id\":\"scrubbed\",\"name\":\"scrubbed\",\"idname\":\"this is fine\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBytes := scrub(tt.url, []byte(tt.input))
			gotStr := string(gotBytes)
			assert.Equal(t, tt.expected, gotStr)
		})

	}
}
