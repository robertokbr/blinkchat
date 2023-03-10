package utils

import (
	"math/rand"
	"time"
)

func Rand(n int) int {
	s1 := rand.NewSource(time.Now().UnixMicro())
	r1 := rand.New(s1)

	return r1.Intn(n)
}
