package main


type Map struct {
	snakePosition Point
	width int
	height int
	fruitPosition Point
}

type Point struct {
	x int
	y int
}