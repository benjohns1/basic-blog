package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	pb "github.com/benjohns1/basic-blog-ms/services/comment/proto"
)

// CommentService wraps the blog comment service client dependencies
type CommentService struct {
	Comment pb.CommentClient
}

// CommentIDResponse responses with a comment ID
type CommentIDResponse struct {
	ID int `json:"id"`
}

// CommentsHandler handles comment requests
func (s *CommentService) CommentsHandler(w http.ResponseWriter, r *http.Request, postIDStr string) {

	// Parse post ID
	postIDInt, err := strconv.ParseInt(postIDStr, 10, 32)
	if err != nil {
		writeError(w, fmt.Errorf("invalid id: %v", err), http.StatusBadRequest)
		return
	}
	postID := int32(postIDInt)

	switch r.Method {
	case "POST":
		// Parse request data
		rawCommentJSON, err := ioutil.ReadAll(r.Body)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
		r.Body.Close()

		var newComment struct {
			Body      string `json:"body"`
			Commenter string `json:"commenter"`
		}
		if err := json.Unmarshal(rawCommentJSON, &newComment); err != nil {
			writeError(w, fmt.Errorf("error parsing request: %v", err), http.StatusBadRequest)
			return
		}

		// Make service RPC call
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		resp, err := s.Comment.New(ctx, &pb.NewCommand{PostId: postID, Body: newComment.Body, Commenter: newComment.Commenter})
		if err != nil {
			writeError(w, fmt.Errorf("internal service error: %v", err), http.StatusInternalServerError)
			return
		}

		commentID := resp.GetId()
		if commentID <= 0 {
			writeError(w, fmt.Errorf("error creating new comment"), http.StatusInternalServerError)
			return
		}

		// Build API response
		commentIDResponse, err := json.Marshal(&CommentIDResponse{ID: int(commentID)})
		if err != nil {
			writeError(w, fmt.Errorf("error processing new comment response %v", err), http.StatusInternalServerError)
			return
		}
		writeResponse(w, r, commentIDResponse, 201)
		return
	default:
		writeError(w, fmt.Errorf("invalid method"), http.StatusBadRequest)
	}
}
