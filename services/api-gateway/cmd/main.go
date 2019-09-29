package main

import (
	"log"
	"os"
	"strconv"

	"github.com/benjohns1/basic-blog-ms/services/api-gateway/internal/api"
	authpb "github.com/benjohns1/basic-blog-ms/services/authentication/proto"
	commentpb "github.com/benjohns1/basic-blog-ms/services/comment/proto"
	postpb "github.com/benjohns1/basic-blog-ms/services/post/proto"
	"google.golang.org/grpc"
)

func main() {

	services := &api.Services{}

	// Authentication grpc service
	authconn, err := grpc.Dial(envStr("AUTHENTICATION_ADDR", "localhost:3001"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("error setting up authentication service connection: %v", err)
	}
	defer authconn.Close()
	services.Authentication = api.AuthenticationService{
		Authentication: authpb.NewAuthenticationClient(authconn),
	}

	// Comment grpc service
	commentconn, err := grpc.Dial(envStr("COMMENT_ADDR", "localhost:3003"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("error setting up comment service connection: %v", err)
	}
	defer commentconn.Close()
	commentClient := commentpb.NewCommentClient(commentconn)
	services.Comment = api.CommentService{
		Comment: commentClient,
	}

	// Post grpc service
	postconn, err := grpc.Dial(envStr("POST_ADDR", "localhost:3002"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("error setting up post service connection: %v", err)
	}
	defer postconn.Close()
	services.Post = api.PostService{
		Authenticate: services.Authentication.Authenticate,
		Post:         postpb.NewPostClient(postconn),
		Comment:      commentClient,
	}

	// Start API Gateway
	api.Start(api.Config{
		APIPort:  envInt("API_PORT", 3000),
		Services: services,
	})
}

func envInt(envKey string, defaultValue int) int {
	if v, exists := os.LookupEnv(envKey); exists {
		value, err := strconv.Atoi(v)
		if err == nil {
			return value
		}
	}
	return defaultValue
}

func envStr(envKey string, defaultValue string) string {
	if v, exists := os.LookupEnv(envKey); exists {
		return v
	}
	return defaultValue
}
