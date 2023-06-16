package main

import "math"

func CheckDualIntersect(a1, b1, a2, b2 float64) bool {
	return a2 <= a1 && a1 <= b2 ||
		a2 <= b1 && b1 <= b2 ||
		a1 <= a2 && b2 <= b1
}

func RoundFloat64ToInt(x float64) int {
	return int(math.Round(x))
}
