package util

import (
	"math/rand"
	"time"
)

// No need to seed the global random generator in Go 1.20+

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	if min >= max {
		return min
	}
	return min + rand.Int63n(max-min+1)
}


func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	length := len(letters)
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(length)]
	}
	return string(b)
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "INR"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}