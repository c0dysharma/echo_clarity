package helpers

import (
	"regexp"
)

func ExtractJSONString(text string) (string, error) {
	// Define the regular expression pattern (same as before).
	pattern := regexp.MustCompile("(?s)`json(.*?)`")

	// Find the first match.
	match := pattern.FindStringSubmatch(text)

	// Extract the captured group (the JSON string) if a match is found.
	if len(match) > 1 {
					return match[1], nil
	}

	// Return an empty string and no error if no match is found.
	return "", nil
}