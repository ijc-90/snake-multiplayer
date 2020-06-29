package main

import (
	"context"
	"fmt"
	"github.com/ijc-90/snake-multiplayer/communication/game_communicator"
	"github.com/ijc-90/snake-multiplayer/communication/matchmaking_communicator"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

var waitingForOpponent bool
var currentGameId int32
const (
	port = ":50052"
)
const (
	gameBoardAddress = "localhost:50051"
)

type message struct{

}

type pendingMatch struct{
	userId int
	gameId int
}

type server struct {
	matchmaking.UnimplementedMatchMakingCommunicatorServer
}


func (s *server) GetGame(ctx context.Context, request *matchmaking.MatchMakingRequest) (*matchmaking.MatchMakingResponse, error){
	fmt.Printf("client connected requesting match\n")

	if waitingForOpponent{
		waitingForOpponent = false
		return &matchmaking.MatchMakingResponse{
			GameId: currentGameId,
			PlayerId: 2,
		}, nil
	}else{
		waitingForOpponent = true


		createdGame, err := requestNewGameCreationToGameServer()

		fmt.Printf(" created game %p ... %v \n", createdGame)
		currentGameId = createdGame.GameId

		matchMakingResponse := matchmaking.MatchMakingResponse{
			GameId: createdGame.GameId,
			PlayerId: 1,
		}
		return &matchMakingResponse, err
	}
}

func requestNewGameCreationToGameServer() (*game_communicator.MatchResponse, error) {
	gameConnection, err := grpc.Dial(gameBoardAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer gameConnection.Close()

	gameServer := game_communicator.NewGameCommunicatorClient(gameConnection)

	newCtx, cancel := context.WithTimeout(context.Background(), time.Second*200)
	defer cancel()

	var matchRequest game_communicator.MatchRequest
	matchRequest = game_communicator.MatchRequest{}

	createdGame, err := gameServer.CreateMatch(newCtx, &matchRequest)

	return createdGame, err
}


func main() {
	waitingForOpponent = false
	fmt.Println("Starting MatchmakingServer")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	matchmaking.RegisterMatchMakingCommunicatorServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}


}