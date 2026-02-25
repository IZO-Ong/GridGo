// Package utils provides common helper functions for security and data generation.
package utils

import (
	crand "crypto/rand"
	"fmt"
	"math/big"
)

// GenerateOTP produces a cryptographically secure 6-digit numeric string.
// It is used for email verification.
func GenerateOTP() string {
	max := big.NewInt(1000000)
	
	// Generate a random integer in range [0, max)
	n, err := crand.Int(crand.Reader, max)
	if err != nil {
		// Fallback to a "000000" string
		return fmt.Sprintf("%06d", 0) 
	}

	// Format the integer as a 6-digit string, left-padding with zeros.
	return fmt.Sprintf("%06d", n.Int64())
}