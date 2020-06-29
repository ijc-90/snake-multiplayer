package commons

import "math"

type Map struct {
	GameId int
	Snakes [2]Snake
	Width int
	Height int
	FruitPosition Point
	GameOver bool
}

type Point struct {
	X int
	Y int
}

func (p Point) IsEqual(otherPoint Point) bool {
	return p.X == otherPoint.X && p.Y == otherPoint.Y
}
func (p Point) Distance(otherPoint Point) int {
	x := p.X - otherPoint.X
	y := p.Y - otherPoint.Y
	return int(math.Abs(float64(x)) + math.Abs(float64(y)))
}

type Snake struct{
	Id int
	Position Point
	Direction int
	Won bool
	Lost bool
}