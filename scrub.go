package main

import (
	"strings"
)

// Perform some scrubbing of data
//
// Must work on raw []byte in case we need to log some JSON that does not parse
// into the expected structure
//
// Scrubbing is URL dependent. We want to scrub the `name` value from a drive entry
// in /proxy/users/drive/api/v2/drives but not the `name` value from /proxy/users/drive/api/v1/systems/info

func scrubKeyString(data, key, replace string) string {
	curr := data
	start := 0
	for start < len(curr) {
		i := strings.Index(curr[start:], "\""+key+"\"")
		if i == -1 {
			break
		}
		// Find the start of the value
		valStart := start + i + 1 + len(key) + 1 + 1 // " <key> " :
		// We should be at the " now
		valEnd := valStart + 1
		for valEnd < len(curr) && curr[valEnd] != '"' {
			valEnd++
		}
		// valEnd should now point to the closing " for the value
		if valEnd == len(curr) {
			// ARGH
			// TODO - how to report error?
			return curr
		}
		curr = curr[0:valStart] + "\"" + replace + curr[valEnd:]
		// Move start far enough along we don't hit the same key again
		start += i + len(key)
	}
	return curr
}

func scrub(url string, data []byte) []byte {
	curr := string(data)
	switch url {
	case "/proxy/drive/api/v2/systems/device-info":
		curr = scrubKeyString(curr, "address", "1.2.3.4")
		curr = scrubKeyString(curr, "mac", "01:23:45:67:89:ab")
		curr = scrubKeyString(curr, "name", "scrubbed")
	case "/proxy/users/drive/api/v1/systems/info":
		curr = scrubKeyString(curr, "deviceId", "scrubbed")
		curr = scrubKeyString(curr, "name", "scrubbed")
		curr = scrubKeyString(curr, "cdn", "scrubbed")
		curr = scrubKeyString(curr, "guid", "scrubbed")
		curr = scrubKeyString(curr, "id", "scrubbed")
	case "/proxy/users/drive/api/v2/drives":
		curr = scrubKeyString(curr, "id", "scrubbed")
		curr = scrubKeyString(curr, "name", "scrubbed")
	default:
		// Nothing to scrub (as far as I know)
		// /proxy/drive/api/v2/systems/disk-stats?start=1777850684&end=1777937984&interval=900
		// /proxy/users/drive/api/v1/systems/identity
		// /proxy/drive/api/v1/systems/performance/file-operations
		// /proxy/drive/api/v2/systems/network-io
		// /proxy/drive/api/v1/systems/storage?type=detail
		// /proxy/users/drive/api/v2/storage
		// Need more detail
		// /proxy/users/drive/api/v2/groups
	}
	return []byte(curr)
}
