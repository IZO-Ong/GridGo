package utils

import (
	crand "crypto/rand"
	"fmt"
	"math/big"
)


func GenerateOTP() string {
	max := big.NewInt(1000000)
	
	n, err := crand.Int(crand.Reader, max)
	if err != nil {
		return fmt.Sprintf("%06d", 0) 
	}
	return fmt.Sprintf("%06d", n.Int64())
}