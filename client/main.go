package main

import (
	"context"
	"log"
	//"time"
	"io"

	"google.golang.org/grpc"
	pb "github.com/ijc-90/snake-multiplayer/gamecommunicator"
)

const (
	address     = "localhost:50051"
)

func main() {
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

	log.Printf("Streaming")
	stream, err := c.SetDirectionsAndUpdateGame(context.Background())
	waitc := make(chan struct{})
	log.Printf("B")
	for {
		log.Printf("Sending one")
		direction := &pb.DirectionRequest{SnakeNumber: 1, SnakeDirection:1}
		if err := stream.Send(direction); err != nil {
			log.Fatalf("Failed to send a note: %v", err)
		}

		log.Printf("Receiving one")
		in, err := stream.Recv()
		if err == io.EOF {
			// read done.
			close(waitc)
			return
		}
		if err != nil {
			log.Fatalf("Failed to receive a note : %v", err)
		}
		log.Printf("Got message, gamestate is %d", in.GameState)


	}
	stream.CloseSend()
}