package main

func Add(a, b int) int {
	return a + b
}

func Sub(a, b int) int {
	return a - b
}

func Test1(x, y int, s string) (int, string) {
	return x + y, s + " tested"
}
