package main

import "math"

const DOUBLE_PI = math.Pi * 2

func CheckDualIntersect(a1, b1, a2, b2 float64) bool {
	return a2 <= a1 && a1 <= b2 ||
		a2 <= b1 && b1 <= b2 ||
		a1 <= a2 && b2 <= b1
}

func RoundFloat64ToInt(x float64) int {
	return int(math.Round(x))
}

func UnwindAngle(angle float64) float64 {
	var i = angle / DOUBLE_PI
	angle -= float64(int64(i)) * DOUBLE_PI
	return angle
}
