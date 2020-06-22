package main

import (
   "fmt"
   "context"
   "io"
   "log"
   "net"
   "google.golang.org/grpc"
   pb "github.com/ijc-90/snake-multiplayer/gamecommunicator"
)


const (
   port = ":50051"
)

type server struct {
   pb.UnimplementedGameCommunicatorServer
}

func (s *server) SetDirectionsAndUpdateGame(stream pb.GameCommunicator_SetDirectionsAndUpdateGameServer) error {
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


      fmt.Printf("Sending game state")
      newGameState := &pb.GameStateRequest{GameState: 1}
      if err := stream.Send(newGameState); err != nil {
         return err
      }
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