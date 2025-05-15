package util

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateRandomString generates a cryptographically secure random string of the specified length.
// It uses crypto/rand for secure random number generation and returns the result as a base64-encoded string.
// The actual length of the returned string might be slightly longer than the requested length due to base64 encoding.
func GenerateRandomString(length int) string {
	if length <= 0 {
		return ""
	}

	// Calculate the number of random bytes needed
	// Base64 encoding: 3 bytes -> 4 characters
	// Therefore, to get at least n characters, we need ceil(n * 3/4) bytes
	bytes := make([]byte, (length*3+3)/4)

	// Read random bytes
	_, err := rand.Read(bytes)
	if err != nil {
		// In case of error, return empty string
		// In production, you might want to handle this differently
		return ""
	}

	// Encode to base64 and trim to desired length
	encoded := base64.URLEncoding.EncodeToString(bytes)
	// Trim to desired length and remove any padding
	result := encoded[:length]

	return result
}