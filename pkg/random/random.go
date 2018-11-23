package random

import (
	"math/rand"
	"time"
)

const UpperCaseBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const LowerCaseBytes = "abcdefghijklmnopqrstuvwxyz"
const Digits = "0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func String(n int, alpha string) string {

	b := make([]byte, n)
	for i := range b {
		b[i] = alpha[rand.Int63()%int64(len(alpha))]
	}
	return string(b)
}
