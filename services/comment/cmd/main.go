package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/benjohns1/basic-blog-ms/services/comment/internal/postgres"
	pb "github.com/benjohns1/basic-blog-ms/services/comment/proto"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

type server struct {
	db *sql.DB
}

// Comment blog post comment
type Comment struct {
	ID          int       `json:"id"`
	Body        string    `json:"body"`
	CreatedTime time.Time `json:"createdTime"`
	Commenter   string    `json:"commenter"`
}

func (s *server) New(ctx context.Context, in *pb.NewCommand) (*pb.NewResponse, error) {

	var id int
	err := s.db.QueryRow(`INSERT INTO comment.comment (body, created_time, post_id, commenter_id) VALUES ($1, $2, $3, $4) RETURNING id`, in.Body, time.Now(), in.PostId, in.Commenter).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &pb.NewResponse{Id: int32(id)}, nil
}

func (s *server) List(ctx context.Context, in *pb.ListQuery) (*pb.ListResponse, error) {

	rows, err := s.db.Query(`SELECT id, body, created_time, commenter_id FROM comment.comment WHERE post_id = $1 ORDER BY created_time DESC`, in.PostId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := pb.ListResponse{}

	for rows.Next() {
		var parse struct {
			createdTime *string
		}

		comment := pb.ViewResponse{}
		err := rows.Scan(&comment.Id, &comment.Body, &parse.createdTime, &comment.Commenter)
		if err != nil {
			return nil, err
		}
		if t, err := time.Parse(time.RFC3339Nano, *parse.createdTime); err == nil {
			comment.CreatedTime = t.Unix()
		}
		comments.Comments = append(comments.Comments, &comment)
	}
	return &comments, nil
}

func main() {
	// environment configs
	port := 3003
	if v, exists := os.LookupEnv("PORT"); exists {
		p, err := strconv.Atoi(v)
		if err == nil {
			port = p
		}
	}

	// db connection
	dbconn := postgres.DBConn{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "asdf1234",
		DBName:   "postgres",
		AppName:  "postService",
	}
	if v, exists := os.LookupEnv("DB_HOST"); exists {
		dbconn.Host = v
	}
	if v, exists := os.LookupEnv("DB_PASSWORD"); exists {
		dbconn.Password = v
	}
	if v, exists := os.LookupEnv("DB_PORT"); exists {
		port, err := strconv.Atoi(v)
		if err == nil {
			dbconn.Port = port
		}
	}

	db, err := postgres.Setup(dbconn)
	if err != nil {
		panic(err)
	}

	// start service
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCommentServer(s, &server{db})
	log.Printf("starting comment service on port %d\n", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
