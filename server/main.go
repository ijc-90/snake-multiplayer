package main

import (
   "fmt"
   "context"
   "math/rand"
   "io"
   "log"
   "time"
   "net"
   "google.golang.org/grpc"
   pb "github.com/ijc-90/snake-multiplayer/gamecommunicator"
   commons "github.com/ijc-90/snake-multiplayer/commons"
)


const (
   port = ":50051"
)

type server struct {
   pb.UnimplementedGameCommunicatorServer
}

func (s *server) SetDirectionsAndUpdateGame(stream pb.GameCommunicator_SetDirectionsAndUpdateGameServer) error {
   var aMap *commons.Map
   aMap = &commons.Map{
      SnakePosition: commons.Point{X:0, Y:0},
      SnakeDirection: 1,
      Width: commons.Width,
      Height: commons.Height,
      FruitPosition: commons.Point{X:0, Y:1},
   }
   go func(mapPointer *commons.Map) error{
      for{
         fmt.Println("Sending game state")
         // Random move
         snakePosition := commons.Point{X: rand.Intn(commons.Width), Y:rand.Intn(commons.Height)}
         fruitPosition := commons.Point{X: rand.Intn(commons.Width), Y:rand.Intn(commons.Height)}

         mapPointer.SnakePosition = snakePosition
         mapPointer.FruitPosition = fruitPosition


         //Convert to request
         snakePositionRequest := &pb.Point{X: int32(mapPointer.SnakePosition.X), Y: int32(mapPointer.SnakePosition.Y)}
         fruitPositionRequest := &pb.Point{X: int32(mapPointer.FruitPosition.X), Y: int32(mapPointer.FruitPosition.Y)}


         aMapRequest := pb.Map{
            SnakePosition:snakePositionRequest,
            FruitPosition:fruitPositionRequest,
            SnakeDirection:int32(mapPointer.SnakeDirection),
            Height: commons.Height,
            Width: commons.Width}
         newGameState := &pb.GameStateRequest{GameState: &aMapRequest}
         if err := stream.Send(newGameState); err != nil {
            return err
         }
         time.Sleep(1000 * time.Millisecond)
      }
   }(aMap)
   for {
      fmt.Println("start receiving")
      in, err := stream.Recv()
      fmt.Println("in: ", in.GetSnakeDirection(), in.GetSnakeNumber())
      if err == io.EOF {
         return nil
      }
      if err != nil {
         return err
      }
      fmt.Printf("Received... snake: %d , direction: %d", in.GetSnakeNumber(), in.GetSnakeDirection())

   }
}
func (s *server) SetDirection(ctx context.Context, in *pb.DirectionRequest) (*pb.DirectionResponse, error) {
   log.Printf("Received: %d , %d", in.GetSnakeNumber(), in.GetSnakeDirection())
   return &pb.DirectionResponse{Received: 1 }, nil
}

func main() {
   fmt.Println("Starting server")
   lis, err := net.Listen("tcp", port)
   if err != nil {
      log.Fatalf("failed to listen: %v", err)
   }
   s := grpc.NewServer()
   pb.RegisterGameCommunicatorServer(s, &server{})
   if err := s.Serve(lis); err != nil {
      log.Fatalf("failed to serve: %v", err)
   }
}