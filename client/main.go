package main

import (
	"bufio"
	"context"
	"fmt"
	commons "github.com/ijc-90/snake-multiplayer/commons"
	pb "github.com/ijc-90/snake-multiplayer/gamecommunicator"
	"io"
	"log"
	"os"
	"time"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
)

func main() {
	var aMap commons.Map
	var fruitPosition commons.Point
	var snakeNumber, gameId int32

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGameCommunicatorClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*200)
	defer cancel()

	gameResponse, err := c.ConnectToGame(ctx, &pb.GameRequest{})
	if err != nil {
		log.Fatalf("Failed to connect to game : %v", err)
	}


	snakeNumber = gameResponse.PlayerId
	gameId = gameResponse.GameId
	fmt.Printf("game id %v \n", gameId)


	stream, err := c.SetDirectionsAndUpdateGame(context.Background())
	waitc := make(chan struct{})

	//Send first request when ready
	direction := &pb.DirectionRequest{SnakeNumber: snakeNumber, SnakeDirection: int32(1),GameId: gameId}
	if err := stream.Send(direction); err != nil {
		log.Fatalf("Failed to send. error: %v", err)
	}


	// Constantly fetch and draw game state
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a note : %v", err)
			}

			//Convert message to game drawable game state
			messageMap := in.GameState
			fruitPosition = commons.Point{
				X: int(messageMap.FruitPosition.X),
				Y: int(messageMap.FruitPosition.Y),
			}

			println("Received new game state")
			fmt.Printf(	"snake numberId %v position %v,%v direction %v",messageMap.SnakeOne.Id,messageMap.SnakeOne.Position.X,messageMap.SnakeOne.Position.Y, messageMap.SnakeOne.Direction)
			fmt.Printf(	"snake numberId %v position %v,%v direction %v",messageMap.SnakeTwo.Id,messageMap.SnakeTwo.Position.X,messageMap.SnakeTwo.Position.Y, messageMap.SnakeTwo.Direction)
			aMap = commons.Map{
				Snakes: [2]commons.Snake{
					commons.Snake{
						Id: int(messageMap.SnakeOne.Id),
						Position: commons.Point{ X: int(messageMap.SnakeOne.Position.X), Y: int(messageMap.SnakeOne.Position.Y) },
						Direction: int(messageMap.SnakeOne.Direction),
					},
					commons.Snake{
						Id: int(messageMap.SnakeTwo.Id),
						Position: commons.Point{ X: int(messageMap.SnakeTwo.Position.X), Y: int(messageMap.SnakeTwo.Position.Y) },
						Direction: int(messageMap.SnakeTwo.Direction),
					},
				},
				FruitPosition: fruitPosition,
				Width:         int(messageMap.Width),
				Height:        int(messageMap.Height),
			}

			DrawMap(aMap)
		}
	}()


	reader := bufio.NewReader(os.Stdin)
	for {
		char, _, err := reader.ReadRune()
		if err == nil{
			if value, found := commons.Directions[char]; found {
				log.Println("match! %v %v", char, value )
				direction := &pb.DirectionRequest{SnakeNumber: snakeNumber, SnakeDirection: int32(value)}
				if err := stream.Send(direction); err != nil {
					log.Fatalf("Failed to send. error: %v", err)
				}
			}
		}

	}
	stream.CloseSend()
}