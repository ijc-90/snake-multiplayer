package main

import (
   "context"
   "fmt"
   commons "github.com/ijc-90/snake-multiplayer/commons"
   pb "github.com/ijc-90/snake-multiplayer/gamecommunicator"
   "google.golang.org/grpc"
   "io"
   "log"
   "net"
   "time"
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

   // report game state every X seconds
   go func(mapPointer *commons.Map) error{
      for{
         fmt.Println("Sending game state")
         // Random move

         x := aMap.SnakePosition.X
         y := aMap.SnakePosition.Y
         switch aMap.SnakeDirection {
         case 1: //UP
            y = commons.Max(0, y - 1)

         case 2: //LEFT
            x = commons.Max(0, x - 1)
         case 3: //DOWN
            y = commons.Min(aMap.Height -1, y + 1)
         case 4: //RIGHT
            x = commons.Min(aMap.Width -1, x + 1)
         }

         snakePosition := commons.Point{X: x, Y: y}

         //fruitPosition := commons.Point{X: rand.Intn(commons.Width), Y:rand.Intn(commons.Height)}
         fruitPosition := commons.Point{X: aMap.FruitPosition.X, Y: aMap.FruitPosition.Y}

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
         time.Sleep(500 * time.Millisecond)
      }
   }(aMap)

   for {
      fmt.Println("start receiving")
      in, err := stream.Recv()
      fmt.Println("in: ", in.GetSnakeDirection(), in.GetSnakeNumber())
      aMap.SnakeDirection = int(in.GetSnakeDirection())
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