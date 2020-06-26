package commons


type Map struct {
	GameId int
	Snakes [2]Snake
	Width int
	Height int
	FruitPosition Point
}

type Point struct {
	X int
	Y int
}

type Snake struct{
	Id int
	Position Point
	Direction int
}