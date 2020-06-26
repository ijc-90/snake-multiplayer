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


type gameBoard struct{
   gameMap commons.Map
   streams []pb.GameCommunicator_SetDirectionsAndUpdateGameServer
}

var games map[int]*gameBoard
var waitingForOpponent bool
var currentGameNumber int

const (
   port = ":50051"
)

type server struct {
   pb.UnimplementedGameCommunicatorServer
}

func tick( mapPointer *commons.Map){
   var newSnakes [2]commons.Snake

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


   fruitPosition := commons.Point{X: mapPointer.FruitPosition.X, Y: mapPointer.FruitPosition.Y}

   mapPointer.FruitPosition = fruitPosition
}

func (s *server) ConnectToGame(ctx context.Context, in *pb.GameRequest) (*pb.GameResponse, error) {
   if waitingForOpponent{
      currentGameBoard, found := games[currentGameNumber]
      if !found{
         log.Fatalf("Game Not Found")
      }
      waitingForOpponent = false
      fmt.Printf("Agame id connected %d, pointer %v", currentGameNumber, currentGameBoard)
      fmt.Println()
      currentGameNumber ++

      return &pb.GameResponse{
         GameId: int32(currentGameBoard.gameMap.GameId),
         PlayerId: 2,
      }, nil
   }else{
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
         GameId: currentGameNumber,
      }
      newGameBoard := &gameBoard{
         gameMap: newMap,
         //streams: make([]pb.GameCommunicator_SetDirectionsAndUpdateGameServer, 0, 2),
      }
      games[currentGameNumber] = newGameBoard
      waitingForOpponent = true


      return &pb.GameResponse{
         GameId: int32(newMap.GameId),
         PlayerId: 1,
      }, nil

   }
}

func (s *server) SetDirectionsAndUpdateGame(stream pb.GameCommunicator_SetDirectionsAndUpdateGameServer) error {
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


            //Convert to request
            err2 := notifyGameState(gameBoardPointer)
            if err2 != nil {
               return err2
            }
            time.Sleep(500 * time.Millisecond)
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
   fruitPositionRequest := &pb.Point{X: int32(mapPointer.FruitPosition.X), Y: int32(mapPointer.FruitPosition.Y)}

   aMapRequest := pb.Map{
      SnakeOne: &pb.Snake{
         Id:        int32(mapPointer.Snakes[0].Id),
         Position:  &pb.Point{X: int32(mapPointer.Snakes[0].Position.X), Y: int32(mapPointer.Snakes[0].Position.Y)},
         Direction: int32(mapPointer.Snakes[0].Direction),
      },
      SnakeTwo: &pb.Snake{
         Id:        int32(mapPointer.Snakes[1].Id),
         Position:  &pb.Point{X: int32(mapPointer.Snakes[1].Position.X), Y: int32(mapPointer.Snakes[1].Position.Y)},
         Direction: int32(mapPointer.Snakes[1].Direction),
      },
      FruitPosition: fruitPositionRequest,
      Height:        commons.Height,
      Width:         commons.Width}
   newGameState := &pb.GameStateRequest{GameState: &aMapRequest}
   for _, stream := range gameBoard.streams {
      if err := stream.Send(newGameState); err != nil {
         log.Printf("Error notifying game state, %v", err)
      }
   }

   return nil
}

func (s *server) SetDirection(ctx context.Context, in *pb.DirectionRequest) (*pb.DirectionResponse, error) {
   log.Printf("Received: %d , %d", in.GetSnakeNumber(), in.GetSnakeDirection())
   return &pb.DirectionResponse{Received: 1 }, nil
}

func main() {
   currentGameNumber = 0
   waitingForOpponent = false
   games = make(map[int]*gameBoard)
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