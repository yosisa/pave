package template

import "math"

func Add(a, b float64) float64 {
	return a + b
}

func Sub(a, b float64) float64 {
	return a - b
}

func Mul(a, b float64) float64 {
	return a * b
}

func Div(a, b float64) float64 {
	return a / b
}

func Major(a int) int {
	return a/2 + 1
}

func init() {
	funcmap["add"] = Add
	funcmap["sub"] = Sub
	funcmap["mul"] = Mul
	funcmap["div"] = Div
	funcmap["major"] = Major
	funcmap["ceil"] = math.Ceil
	funcmap["floor"] = math.Floor
}
