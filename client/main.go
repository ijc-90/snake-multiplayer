package main

import (
	"context"
	"log"
	"time"

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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SetDirection(ctx, &pb.DirectionRequest{SnakeNumber: 1, SnakeDirection: 1})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Direction sent, Response code: %s", r.GetReceived())
}