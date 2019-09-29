package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	pb "github.com/benjohns1/basic-blog-ms/services/authentication/proto"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Login(ctx context.Context, in *pb.LoginCommand) (*pb.LoginResponse, error) {
	if in.GetUsername() == "bobross" && in.GetPassword() == "painter" {
		return &pb.LoginResponse{Success: true, Token: "notSoSecureToken!"}, nil
	}
	return &pb.LoginResponse{Success: false}, nil
}

func (s *server) Authenticate(ctx context.Context, in *pb.AuthenticateQuery) (*pb.AuthenticateResponse, error) {
	if in.GetToken() == "notSoSecureToken!" {
		return &pb.AuthenticateResponse{Success: true}, nil
	}
	return &pb.AuthenticateResponse{Success: false}, nil
}

func main() {
	// environment configs
	port := 3001
	if v, exists := os.LookupEnv("PORT"); exists {
		p, err := strconv.Atoi(v)
		if err == nil {
			port = p
		}
	}

	// start service
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAuthenticationServer(s, &server{})
	log.Printf("starting authentication service on port %d\n", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
