package main

import (
	"regexp"
	"strings"
)

func parseEmuListing(input string) []string {
	re := regexp.MustCompile(`\b[0-9A-Fa-f]{4}:\s+((?:[0-9A-Fa-f]{2}\s+)+)`)

	lines := strings.Split(input, "\n")

	var hexValues []string
	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			hexValuesInLine := strings.Fields(matches[1])
			hexValues = append(hexValues, hexValuesInLine...)
		}
	}

	return hexValues
}
