package urlGenerator

import (
	"math/rand"
)

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"
)

func Generate(size uint) string {
	res := make([]byte, size)
	for i := range size {
		res[i] = letters[rand.Intn(len(letters))]
	}
	return string(res)
}
