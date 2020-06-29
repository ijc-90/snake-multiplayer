package main

import (
	"context"
	"fmt"
	commons "github.com/ijc-90/snake-multiplayer/commons"
	"github.com/ijc-90/snake-multiplayer/communication/game_communicator"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"time"
)


type gameBoard struct{
	gameMap commons.Map
	streams []game_communicator.GameCommunicator_SetDirectionsAndUpdateGameServer
}

var games map[int]*gameBoard
var waitingForOpponent bool
var currentGameNumber int


const (
	port = ":50051"
)

type server struct {
	game_communicator.UnimplementedGameCommunicatorServer
}

func tick( mapPointer *commons.Map){
	var newSnakes [2]commons.Snake

	var oldSnakePositions [2]commons.Point
	oldSnakePositions[0] = commons.Point{X: mapPointer.Snakes[0].Position.X,Y: mapPointer.Snakes[0].Position.Y}
	oldSnakePositions[1] = commons.Point{X: mapPointer.Snakes[1].Position.X,Y: mapPointer.Snakes[1].Position.Y}

	for index, snake := range mapPointer.Snakes{
		x := snake.Position.X
		y := snake.Position.Y
		switch snake.Direction {
		case 1: //UP
			y = commons.Max(0, y - 1)
		case 2: //LEFT
			x = commons.Max(0, x - 1)
		case 3: //DOWN
			y = commons.Min(mapPointer.Height -1, y + 1)
		case 4: //RIGHT
			x = commons.Min(mapPointer.Width -1, x + 1)
		}

		newSnakes[index] = commons.Snake{
			Position: commons.Point{X: x, Y: y},
			Direction: snake.Direction,
			Id: snake.Id,
		}
	}
	mapPointer.Snakes = newSnakes

	//Detect collision
	if mapPointer.Snakes[0].Position.IsEqual(mapPointer.Snakes[1].Position) ||
	 (mapPointer.Snakes[0].Position.IsEqual(oldSnakePositions[1]) && mapPointer.Snakes[1].Position.IsEqual(oldSnakePositions[0]) ) {
		mapPointer.Snakes[0].Won = false
		mapPointer.Snakes[0].Lost = true
		mapPointer.Snakes[1].Won = false
		mapPointer.Snakes[1].Lost = true
		mapPointer.GameOver = true
	}


	fruitPosition := commons.Point{X: mapPointer.FruitPosition.X, Y: mapPointer.FruitPosition.Y}

	mapPointer.FruitPosition = fruitPosition
}


func (s *server) CreateMatch(ctx context.Context, in *game_communicator.MatchRequest) (*game_communicator.MatchResponse, error) {
	fmt.Printf("Creating match, currentGameNumber %v \n", currentGameNumber)
	error := internalCreateMatch(currentGameNumber)
	response := &game_communicator.MatchResponse{
		GameId: int32(currentGameNumber),
	}
	game, found := games[currentGameNumber]
	fmt.Printf("game %v  found %v \n", game,found)
	currentGameNumber ++
	return response,error
}

func internalCreateMatch(gameNumber int) error{
	newMap := commons.Map{
		Snakes:[2]commons.Snake{
			commons.Snake{
				Id: 1,
				Position: commons.Point{X: 0, Y: 0},
				Direction: 3,
			},
			commons.Snake{
				Id: 2,
				Position: commons.Point{X: commons.Width-1, Y: commons.Height-1},
				Direction: 3,
			},

		},
		Width: commons.Width,
		Height: commons.Height,
		FruitPosition: commons.Point{X:0, Y:1},
		GameId: gameNumber,
		GameOver: false,
	}
	newGameBoard := &gameBoard{
		gameMap: newMap,
	}
	games[gameNumber] = newGameBoard
	return nil
}


func (s *server) SetDirectionsAndUpdateGame(stream game_communicator.GameCommunicator_SetDirectionsAndUpdateGameServer) error {
	var gameBoardPointer *gameBoard

	fmt.Println("Receive first message to start game")
	in, err := stream.Recv()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}
	fmt.Println("Game: %v", in.GameId)

	gameBoardPointer, found := games[int(in.GameId)]
	if !found{
		log.Fatalf("Error retreiving game")
	}
	gameBoardPointer.gameMap.Snakes[in.GetSnakeNumber()-1].Direction = int(in.GetSnakeDirection())
	gameBoardPointer.streams = append(gameBoardPointer.streams, stream)

	//Both players connected we start game and notify them
	if len(gameBoardPointer.streams) == 2{
		// report game state every X seconds
		go func(gameBoardPointer *gameBoard) error{
			for{
				//fmt.Println("Ticking on map with gameid ", gameBoardPointer.gameMap.GameId)
				tick(&gameBoardPointer.gameMap)


				err2 := notifyGameState(gameBoardPointer)
				if err2 != nil {
					return err2
				}

				if (gameBoardPointer.gameMap.GameOver){
					return nil
				}
				time.Sleep(commons.TickInterval * time.Millisecond)
			}
		}(gameBoardPointer)
	}

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		gameBoardPointer.gameMap.Snakes[in.GetSnakeNumber()-1].Direction = int(in.GetSnakeDirection())
	}

}

func notifyGameState(gameBoard *gameBoard) error {
	mapPointer := gameBoard.gameMap
	fruitPositionRequest := &game_communicator.Point{X: int32(mapPointer.FruitPosition.X), Y: int32(mapPointer.FruitPosition.Y)}

	aMapRequest := game_communicator.Map{
		SnakeOne: &game_communicator.Snake{
			Id:        int32(mapPointer.Snakes[0].Id),
			Position:  &game_communicator.Point{X: int32(mapPointer.Snakes[0].Position.X), Y: int32(mapPointer.Snakes[0].Position.Y)},
			Direction: int32(mapPointer.Snakes[0].Direction),
			Won: mapPointer.Snakes[0].Won,
			Lost: mapPointer.Snakes[0].Lost,
		},
		SnakeTwo: &game_communicator.Snake{
			Id:        int32(mapPointer.Snakes[1].Id),
			Position:  &game_communicator.Point{X: int32(mapPointer.Snakes[1].Position.X), Y: int32(mapPointer.Snakes[1].Position.Y)},
			Direction: int32(mapPointer.Snakes[1].Direction),
			Won: mapPointer.Snakes[1].Won,
			Lost: mapPointer.Snakes[1].Lost,
		},
		FruitPosition: fruitPositionRequest,
		Height:        commons.Height,
		Width:         commons.Width,
		GameOver: mapPointer.GameOver,
	}
	newGameState := &game_communicator.GameStateRequest{GameState: &aMapRequest}
	for _, stream := range gameBoard.streams {
		if err := stream.Send(newGameState); err != nil {
			log.Printf("Error notifying game state, %v\n", err)
		}
	}

	return nil
}

func (s *server) SetDirection(ctx context.Context, in *game_communicator.DirectionRequest) (*game_communicator.DirectionResponse, error) {
	log.Printf("Received: %d , %d", in.GetSnakeNumber(), in.GetSnakeDirection())
	return &game_communicator.DirectionResponse{Received: 1 }, nil
}

func main() {
	currentGameNumber = 0
	waitingForOpponent = false
	games = make(map[int]*gameBoard)

	/*//Creating harcoded game 1
	newMap := commons.Map{
		Snakes:[2]commons.Snake{
			commons.Snake{
				Id: 1,
				Position: commons.Point{X: 0, Y: 0},
				Direction: 3,
			},
			commons.Snake{
				Id: 2,
				Position: commons.Point{X: commons.Width-1, Y: commons.Height-1},
				Direction: 3,
			},

		},
		Width: commons.Width,
		Height: commons.Height,
		FruitPosition: commons.Point{X:0, Y:1},
		GameId: 1,
	}
	newGameBoard := &gameBoard{
		gameMap: newMap,
	}
	games[1] = newGameBoard

	 */



	fmt.Println("Starting game server")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	game_communicator.RegisterGameCommunicatorServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}