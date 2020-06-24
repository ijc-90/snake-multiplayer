package commons


type Map struct {
	SnakePosition Point
	SnakeDirection int
	Width int
	Height int
	FruitPosition Point
}

type Point struct {
	X int
	Y int
}