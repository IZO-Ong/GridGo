// Package utils provides helper functions for ID generation and security.
package utils

import (
	crand "crypto/rand"
	"fmt"
	"math/big"
)

// GenerateMazeID creates an ID in the format M-123456-A.
func GenerateMazeID() string {
	digits := "0123456789"
	letters := "ABCDEFGHIJKLNOPQRSTUVWXYZ"

	digitPart := generateRandomString(digits, 6)
	letterPart := generateRandomString(letters, 1)

	return fmt.Sprintf("M-%s-%s", digitPart, letterPart)
}

// GeneratePostID creates an ID with a 'P-' prefix followed by 12 random characters.
func GeneratePostID() string {
	charset := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	return "P-" + generateRandomString(charset, 12)
}

// GenerateCommentID creates an ID with a 'C-' prefix followed by 12 random characters.
func GenerateCommentID() string {
	charset := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	return "C-" + generateRandomString(charset, 12)
}

// generateRandomString is a helper that picks 'n' random characters from the provided charset
// using a cryptographically secure random number generator.
func generateRandomString(charset string, n int) string {
	result := make([]byte, n)
	for i := range result {
		num, _ := crand.Int(crand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}
	return string(result)
}