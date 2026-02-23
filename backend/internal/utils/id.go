package utils

import (
	crand "crypto/rand"
	"fmt"
	"math/big"
)

func GenerateMazeID() string {
	const digits = "0123456789"
	const letters = "ABCDEFGHIJKLNOPQRSTUVWXYZ" 

	digitPart := make([]byte, 6)
	for i := 0; i < 6; i++ {
		num, _ := crand.Int(crand.Reader, big.NewInt(int64(len(digits))))
		digitPart[i] = digits[num.Int64()]
	}

	letIdx, _ := crand.Int(crand.Reader, big.NewInt(int64(len(letters))))
	lastLetter := letters[letIdx.Int64()]

	return fmt.Sprintf("M-%s-%c", string(digitPart), lastLetter)
}