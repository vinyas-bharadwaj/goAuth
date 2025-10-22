package utils

import (
	"crypto/rand"
	"fmt"
)

func GenerateOTP() string {
	b := make([]byte, 3)
	rand.Read(b)
	otp := fmt.Sprintf("%06d", (int(b[0])<<16|int(b[1])<<8|int(b[2]))%1000000)
	return otp
}