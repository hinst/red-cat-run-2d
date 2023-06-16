package main

func Check4Intersect(a1, b1, a2, b2 float64) bool {
	return a2 <= a1 && a1 <= b2 ||
		a2 <= b1 && b1 <= b2 ||
		a1 <= a2 && b2 <= b1
}
