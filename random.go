package main

import (
	"math/rand"
	"time"
)

func GetRandomNumberBetween(a int, b int) int {
	return a + rand.Intn(1+b-a)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
