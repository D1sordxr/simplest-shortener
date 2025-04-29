package pkg

import "math/rand"

type Generator interface {
	GenerateRandomString(newLen int) string
}

type SharedGenerator struct{}

const (
	newLength = 8
	charset   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func (SharedGenerator) GenerateRandomString(newLen int) string {
	if newLen <= 0 {
		newLen = newLength
	}

	code := make([]byte, newLen)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}
