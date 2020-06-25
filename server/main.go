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
      //SnakePosition: commons.Point{X:0, Y:0},
      //SnakeDirection: 1,
      Width: commons.Width,
      Height: commons.Height,
      FruitPosition: commons.Point{X:0, Y:1},
   }

   // report game state every X seconds
   go func(mapPointer *commons.Map) error{
      for{
         //fmt.Println("Sending game state")
         // Random move
         var newSnakes [2]commons.Snake

         for index, snake := range mapPointer.Snakes{
            //fmt.Printf(	"snake number %v position %v,%v direction %v",snake.Id,snake.Position.X,snake.Position.Y, snake.Direction)
            //fmt.Println()

            x := snake.Position.X
            y := snake.Position.Y
            switch snake.Direction {
            case 1: //UP
               y = commons.Max(0, y - 1)
            case 2: //LEFT
               x = commons.Max(0, x - 1)
            case 3: //DOWN
               y = commons.Min(aMap.Height -1, y + 1)
            case 4: //RIGHT
               x = commons.Min(aMap.Width -1, x + 1)
            }

            newSnakes[index] = commons.Snake{
               Position: commons.Point{X: x, Y: y},
               Direction: snake.Direction,
               Id: snake.Id,
            }
         }
         //fmt.Println("snakes before tick %v",mapPointer.Snakes)
         mapPointer.Snakes = newSnakes
         //fmt.Println("snakes after tick %v",mapPointer.Snakes)


         //fruitPosition := commons.Point{X: rand.Intn(commons.Width), Y:rand.Intn(commons.Height)}
         fruitPosition := commons.Point{X: aMap.FruitPosition.X, Y: aMap.FruitPosition.Y}

         mapPointer.FruitPosition = fruitPosition


         //Convert to request
         fruitPositionRequest := &pb.Point{X: int32(mapPointer.FruitPosition.X), Y: int32(mapPointer.FruitPosition.Y)}


         aMapRequest := pb.Map{
            SnakeOne: &pb.Snake{
               Id: int32(mapPointer.Snakes[0].Id),
               Position: &pb.Point{X: int32(mapPointer.Snakes[0].Position.X), Y: int32(mapPointer.Snakes[0].Position.Y)} ,
               Direction: int32(mapPointer.Snakes[0].Direction),
            },
            SnakeTwo: &pb.Snake{
               Id: int32(mapPointer.Snakes[1].Id),
               Position: &pb.Point{X: int32(mapPointer.Snakes[1].Position.X), Y: int32(mapPointer.Snakes[1].Position.Y)} ,
               Direction: int32(mapPointer.Snakes[1].Direction),
            },
            FruitPosition:fruitPositionRequest,
            Height: commons.Height,
            Width: commons.Width}
         newGameState := &pb.GameStateRequest{GameState: &aMapRequest}
         fmt.Printf(	"snake numberone %v position %v,%v direction %v",aMapRequest.SnakeOne.Id,aMapRequest.SnakeOne.Position.X,aMapRequest.SnakeOne.Position.Y, aMapRequest.SnakeOne.Direction)
         fmt.Println()
         fmt.Printf(	"snake numbertwo %v position %v,%v direction %v",aMapRequest.SnakeTwo.Id,aMapRequest.SnakeTwo.Position.X,aMapRequest.SnakeTwo.Position.Y, aMapRequest.SnakeTwo.Direction)
         fmt.Println()
         fmt.Println()
         //fmt.Println("%v", aMapRequest)
         //fmt.Println("%v", mapPointer.Snakes[1])
         //fmt.Println("snaketwo position %v", aMapRequest.SnakeTwo.Position)
         if err := stream.Send(newGameState); err != nil {
            return err
         }
         time.Sleep(200 * time.Millisecond)
      }
   }(aMap)

   fmt.Println("start receiving")
   for {
      in, err := stream.Recv()
      fmt.Println("in: ", in.GetSnakeDirection(), in.GetSnakeNumber())
      aMap.Snakes[in.GetSnakeNumber()-1].Direction = int(in.GetSnakeDirection())
      if err == io.EOF {
         return nil
      }
      if err != nil {
         return err
      }
      fmt.Println("Received... snake: %v , direction: %v", in.GetSnakeNumber(), in.GetSnakeDirection())

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