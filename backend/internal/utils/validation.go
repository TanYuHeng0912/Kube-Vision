package utils

import (
	"regexp"
)

// ValidateContainerID validates a container ID to prevent path traversal and injection
func ValidateContainerID(id string) bool {
	if id == "" {
		return false
	}

	// Container IDs are typically 12-64 characters, alphanumeric with hyphens
	// This prevents path traversal and injection attacks
	matched, err := regexp.MatchString("^[a-zA-Z0-9_-]{12,64}$", id)
	if err != nil {
		return false
	}

	return matched
}

// SanitizeString removes potentially dangerous characters from user input
func SanitizeString(s string) string {
	// Remove null bytes and control characters
	re := regexp.MustCompile(`[\x00-\x1F\x7F]`)
	return re.ReplaceAllString(s, "")
}

