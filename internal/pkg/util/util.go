package util

import (
	"math/rand"
	"time"
)

func GenerateRandomNumber(n int) string {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	charsets := []rune("0123456789")
	letters := make([]rune, n)
	for i := range letters {
		letters[i] = charsets[r.Intn(len(charsets))]
	}
	return string(letters)
}