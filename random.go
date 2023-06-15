package main

import "math/rand"

func GetRandomNumberBetween(a int, b int) int {
	return a + rand.Intn(1+b-a)
}
