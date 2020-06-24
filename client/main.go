package main

import (
	"bufio"
	"context"
	"fmt"
	commons "github.com/ijc-90/snake-multiplayer/commons"
	pb "github.com/ijc-90/snake-multiplayer/gamecommunicator"
	bla "github.com/paulrademacher/climenu"
	"io"
	"log"
	"os"

	//"time"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
)

func main() {
	var aMap commons.Map
	var snakePosition, fruitPosition commons.Point

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGameCommunicatorClient(conn)


	stream, err := c.SetDirectionsAndUpdateGame(context.Background())
	waitc := make(chan struct{})

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

			DrawMap(aMap)
		}
	}()

	/*
	// disable input buffering
	exec.Command("/bin/stty", "-F", "/dev/tty", "-icanon", "min", "1")
	//exec.Command("/bin/stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()

	// do not display entered characters on the screen
	exec.Command("/bin/stty", "-F", "/dev/tty", "-echo").Run()
	var b []byte = make([]byte, 1)
	_,_ = os.Stdin.Read(b)
	fmt.Println("I got the byte", b, "("+string(b)+")")

	*/

	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	fmt.Printf("buffio char %s \n", char)

	for {
		char, _, _ := bla.GetChar()
		fmt.Printf("getchar char %s \n", char)
		if err == nil{
			if value, found := commons.Directions[rune(char)]; found {
				fmt.Printf("match! %v %v\n", char, value )
				direction := &pb.DirectionRequest{SnakeNumber: 1, SnakeDirection: int32(value)}
				if err := stream.Send(direction); err != nil {
					log.Fatalf("Failed to send. error: %v", err)
				}
			}
		}

	}
	stream.CloseSend()
}