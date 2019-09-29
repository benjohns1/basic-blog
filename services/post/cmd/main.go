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

	"github.com/benjohns1/basic-blog-ms/services/post/internal/postgres"
	pb "github.com/benjohns1/basic-blog-ms/services/post/proto"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

type server struct {
	db *sql.DB
}

func (s *server) New(ctx context.Context, in *pb.NewCommand) (*pb.NewResponse, error) {
	var id int
	err := s.db.QueryRow(`INSERT INTO post.post (title, body, created_time, author_id) VALUES ($1, $2, $3, $4) RETURNING id`, in.Title, in.Body, time.Now(), "author").Scan(&id)
	if err != nil {
		return &pb.NewResponse{Id: 0}, err
	}

	return &pb.NewResponse{Id: int32(id)}, nil
}

func (s *server) Delete(ctx context.Context, in *pb.DeleteCommand) (*pb.DeleteResponse, error) {
	ID := in.GetId()
	result, err := s.db.Exec("UPDATE post.post SET deleted = TRUE WHERE id = $1", ID)
	if err != nil {
		return &pb.DeleteResponse{Success: false}, err
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected != 1 || err != nil {
		return &pb.DeleteResponse{Success: false}, fmt.Errorf("unable to update post db entry %v", ID)
	}
	return &pb.DeleteResponse{Success: true}, nil
}

func (s *server) Restore(ctx context.Context, in *pb.RestoreCommand) (*pb.RestoreResponse, error) {
	ID := in.GetId()
	result, err := s.db.Exec("UPDATE post.post SET deleted = FALSE WHERE id = $1", ID)
	if err != nil {
		return &pb.RestoreResponse{Success: false}, err
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected != 1 || err != nil {
		return &pb.RestoreResponse{Success: false}, fmt.Errorf("unable to update post db entry %v", ID)
	}
	return &pb.RestoreResponse{Success: true}, nil
}

func (s *server) Edit(ctx context.Context, in *pb.EditCommand) (*pb.EditResponse, error) {
	ID := in.GetId()
	result, err := s.db.Exec("UPDATE post.post SET title = $1, body = $2 WHERE id = $3", in.GetTitle(), in.GetBody(), ID)
	if err != nil {
		return &pb.EditResponse{Success: false}, err
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected != 1 || err != nil {
		return &pb.EditResponse{Success: false}, fmt.Errorf("unable to update post db entry %v", ID)
	}
	return &pb.EditResponse{Success: true}, nil
}

func (s *server) List(ctx context.Context, in *pb.ListQuery) (*pb.ListResponse, error) {

	where := ""
	if !in.GetIncludeDeleted() {
		where = " WHERE deleted IS NOT TRUE"
	}
	query := fmt.Sprintf(`SELECT id, title, created_time, author_id, deleted FROM post.post%v ORDER BY created_time DESC`, where)
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := pb.ListResponse{}
	for rows.Next() {
		var parse struct {
			createdTime *string
			deleted     *bool
		}

		post := pb.ViewResponse{}
		err := rows.Scan(&post.Id, &post.Title, &parse.createdTime, &post.Author, &parse.deleted)
		if err != nil {
			return nil, err
		}
		if t, err := time.Parse(time.RFC3339Nano, *parse.createdTime); err == nil {
			post.CreatedTime = t.Unix()
		}
		if parse.deleted != nil {
			post.Deleted = *parse.deleted
		}
		posts.Posts = append(posts.Posts, &post)
	}

	return &posts, nil
}

func (s *server) View(ctx context.Context, in *pb.ViewQuery) (*pb.ViewResponse, error) {

	andWhere := ""
	if !in.GetIncludeDeleted() {
		andWhere = " AND deleted IS NOT TRUE"
	}

	query := fmt.Sprintf(`SELECT id, title, body, created_time, author_id, deleted FROM post.post WHERE id = $1%v`, andWhere)
	row := s.db.QueryRow(query, in.GetId())
	var parse struct {
		createdTime *string
		deleted     *bool
	}
	post := pb.ViewResponse{}
	err := row.Scan(&post.Id, &post.Title, &post.Body, &parse.createdTime, &post.Author, &parse.deleted)
	if err != nil {
		return nil, err
	}
	if t, err := time.Parse(time.RFC3339Nano, *parse.createdTime); err == nil {
		post.CreatedTime = t.Unix()
	}
	if parse.deleted != nil {
		post.Deleted = *parse.deleted
	}
	/*
		if comments, err := scanComments(postID, db); err == nil {
			post.Comments = comments
		} else {
			return nil, err
		}*/

	return &post, nil
}

func main() {
	// environment configs
	port := 3002
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
	pb.RegisterPostServer(s, &server{db})
	log.Printf("starting post service on port %d\n", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
