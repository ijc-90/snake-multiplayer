package main

import (
	"context"
	"log"
	"io"
	//"time"
	"google.golang.org/grpc"
	pb "github.com/ijc-90/snake-multiplayer/gamecommunicator"
	commons "github.com/ijc-90/snake-multiplayer/commons"
)

const (
	address     = "localhost:50051"
)

func main() {
	var aMap commons.Map
	var snakePosition, fruitPosition commons.Point
	/*snakePosition = commons.Point{X:12,Y:4}
	fruitPosition = commons.Point{X:3,Y:3}
	aMap = commons.Map{
		SnakePosition: snakePosition,
		FruitPosition: fruitPosition,
		Width: commons.Width,
		Height: commons.Height}

	DrawMap(aMap)*/
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGameCommunicatorClient(conn)

	// Contact the server and print out its response.
	/*
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	log.Printf("Send one Snake Direction on a direct call")
	r, err := c.SetDirection(ctx, &pb.DirectionRequest{SnakeNumber: 1, SnakeDirection: 1})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Direction sent, Response code: %s", r.GetReceived())*/

	stream, err := c.SetDirectionsAndUpdateGame(context.Background())
	waitc := make(chan struct{})

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
			snakePosition = commons.Point{
				X: int(messageMap.SnakePosition.X),
				Y: int(messageMap.SnakePosition.Y),
			}
			fruitPosition = commons.Point{
				X: int(messageMap.FruitPosition.X),
				Y: int(messageMap.FruitPosition.Y),
			}
			aMap = commons.Map{
				SnakePosition: snakePosition,
				FruitPosition: fruitPosition,
				Width:         int(messageMap.Width),
				Height:        int(messageMap.Height),
			}

			//Draw map
			DrawMap(aMap)
		}
	}()

	for {
		direction := &pb.DirectionRequest{SnakeNumber: 1, SnakeDirection:1}
		if err := stream.Send(direction); err != nil {
			log.Fatalf("Failed to send. error: %v", err)
		}

	}
	stream.CloseSend()
}