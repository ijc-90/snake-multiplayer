package main

import (
   "fmt"
   "context"
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