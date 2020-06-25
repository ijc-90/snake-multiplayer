package commons


type Map struct {
	Snakes [2]Snake
	GameId int
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