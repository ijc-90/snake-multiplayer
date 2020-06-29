package main

import (
	"bufio"
	"context"
	"fmt"
	commons "github.com/ijc-90/snake-multiplayer/commons"
	"github.com/ijc-90/snake-multiplayer/communication/game_communicator"
	"github.com/ijc-90/snake-multiplayer/communication/matchmaking_communicator"

	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"time"
)

const (
	gameBoardAddress = "localhost:50051"
	matchMakingAddress = "localhost:50052"
)

var gameOver = false

func main() {
	var snakeNumber int32
	var  gameId int32

	// Set up a connection to the matchMaking.
	matchMakingConn, err := grpc.Dial(matchMakingAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer matchMakingConn.Close()
	matchMakingServer := matchmaking.NewMatchMakingCommunicatorClient(matchMakingConn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*200)
	defer cancel()

	//Request game to matchmaking
	fmt.Printf("Requesting game to matchmaking\n")
	matchMakingResponse, err := matchMakingServer.GetGame(ctx, &matchmaking.MatchMakingRequest{})
	fmt.Printf("matchmaking response %v\n", matchMakingResponse)
	if err != nil {
		log.Fatalf("Failed to connect to game : %v", err)
	}




	snakeNumber = matchMakingResponse.PlayerId
	gameId = matchMakingResponse.GameId
	fmt.Printf("Matchmaking returned game %v \n", int(gameId))

	// Set up a connection to the game.
	gameConnection, err := grpc.Dial(gameBoardAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer gameConnection.Close()

	gameServer := game_communicator.NewGameCommunicatorClient(gameConnection)

	fmt.Printf("Starting stream communication with game server\n")
	stream, err := gameServer.SetDirectionsAndUpdateGame(context.Background())

	waitc := make(chan struct{})



	// Constantly fetch and draw game state
	go func() {
		for ! gameOver {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a game state : %v", err)
			}

			//Convert message to game drawable game state
			messageMap := in.GameState

			var aMap commons.Map
			var fruitPosition commons.Point
			fruitPosition = commons.Point{
				X: int(messageMap.FruitPosition.X),
				Y: int(messageMap.FruitPosition.Y),
			}

			aMap = commons.Map{
				Snakes: [2]commons.Snake{
					commons.Snake{
						Id: int(messageMap.SnakeOne.Id),
						Position: commons.Point{ X: int(messageMap.SnakeOne.Position.X), Y: int(messageMap.SnakeOne.Position.Y) },
						Direction: int(messageMap.SnakeOne.Direction),
						Won: messageMap.SnakeOne.Won,
						Lost: messageMap.SnakeOne.Lost,
					},
					commons.Snake{
						Id: int(messageMap.SnakeTwo.Id),
						Position: commons.Point{ X: int(messageMap.SnakeTwo.Position.X), Y: int(messageMap.SnakeTwo.Position.Y) },
						Direction: int(messageMap.SnakeTwo.Direction),
						Won: messageMap.SnakeTwo.Won,
						Lost: messageMap.SnakeTwo.Lost,
					},
				},
				FruitPosition: fruitPosition,
				Width:         int(messageMap.Width),
				Height:        int(messageMap.Height),
				GameOver: messageMap.GameOver,
			}
			gameOver = aMap.GameOver

			DrawMap(aMap, int(snakeNumber))
		}
	}()


	reader := bufio.NewReader(os.Stdin)
	for !gameOver {
		char, _, err := reader.ReadRune()
		if err == nil && !gameOver{
			if value, found := commons.Directions[char]; found {
				log.Println("match! %v %v", char, value )
				direction := &game_communicator.DirectionRequest{SnakeNumber: snakeNumber, SnakeDirection: int32(value), GameId: gameId}
				if err := stream.Send(direction); err != nil {
					log.Fatalf("Failed to send. error: %v", err)
				}
			}
		}

	}
	stream.CloseSend()
}