package util

import (
	"math/rand"
	"time"
)

/**
 * @see https://www.calhoun.io/creating-random-strings-in-go/
 */

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(length int) string {
	return RandomStringWithCharset(length, charset)
}

func RandomStringWithCharset(length int, charset string) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
