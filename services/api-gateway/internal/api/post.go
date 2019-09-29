package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	cpb "github.com/benjohns1/basic-blog-ms/services/comment/proto"
	pb "github.com/benjohns1/basic-blog-ms/services/post/proto"
)

// PostService wraps the blog post service client
type PostService struct {
	Authenticate func(r *http.Request) (bool, error)
	Post         pb.PostClient
	Comment      cpb.CommentClient
}

// PostFull blog post with a post body
type PostFull struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	CreatedTime time.Time `json:"createdTime"`
	Body        string    `json:"body"`
	Author      string    `json:"author"`
	Deleted     bool      `json:"deleted"`
	//Comments    []Comment `json:"comments"`
}

// PostUpdate blog post with updatable fields
type PostUpdate struct {
	Title       *string    `json:"title"`
	CreatedTime *time.Time `json:"createdTime"`
	Body        *string    `json:"body"`
	Author      *string    `json:"author"`
	Deleted     *bool      `json:"deleted"`
}

// Post blog post without the post body
type Post struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	CreatedTime time.Time `json:"createdTime"`
	Author      string    `json:"author"`
	Deleted     bool      `json:"deleted"`
}

// PostIDResponse responses with a post ID
type PostIDResponse struct {
	ID int `json:"id"`
}

// PostsHandler handles posts requests
func (s *PostService) PostsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Parse API request
		authenticated, err := s.Authenticate(r)
		if err != nil {
			writeError(w, err, http.StatusUnauthorized)
			return
		}

		// Make service RPC call
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		pbResp, err := s.Post.List(ctx, &pb.ListQuery{IncludeDeleted: authenticated})
		if err != nil {
			writeError(w, fmt.Errorf("internal service error: %v", err), http.StatusInternalServerError)
			return
		}

		// Build API response
		apiResp, err := json.Marshal(&pbResp.Posts)
		if err != nil {
			writeError(w, fmt.Errorf("error processing response %v", err), http.StatusInternalServerError)
			return
		}
		writeResponse(w, r, apiResp, 200)
		return
	case "POST":
		authenticated, err := s.Authenticate(r)
		if err != nil {
			writeError(w, err, http.StatusUnauthorized)
			return
		}
		if !authenticated {
			writeError(w, fmt.Errorf("not authenticated"), http.StatusUnauthorized)
			return
		}

		// Parse API request
		postData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
		r.Body.Close()

		post := PostFull{}
		if err := json.Unmarshal(postData, &post); err != nil {
			writeError(w, fmt.Errorf("could not parse request: %v", err), http.StatusBadRequest)
			return
		}

		// Make service RPC call
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		resp, err := s.Post.New(ctx, &pb.NewCommand{Title: post.Title, Body: post.Body})
		if err != nil {
			writeError(w, fmt.Errorf("internal service error: %v", err), http.StatusInternalServerError)
			return
		}

		postID := resp.GetId()
		if postID <= 0 {
			writeError(w, fmt.Errorf("error creating new post"), http.StatusInternalServerError)
			return
		}

		// Build API response
		postIDResponse, err := json.Marshal(&PostIDResponse{ID: int(postID)})
		if err != nil {
			writeError(w, fmt.Errorf("error processing new post response %v", err), http.StatusInternalServerError)
			return
		}
		writeResponse(w, r, postIDResponse, 201)
	default:
		writeError(w, fmt.Errorf("invalid method"), http.StatusBadRequest)
	}
}

// PostHandler handles post requests
func (s *PostService) PostHandler(w http.ResponseWriter, r *http.Request, postIDStr string) {
	// Parse post ID
	postIDInt, err := strconv.ParseInt(postIDStr, 10, 32)
	if err != nil {
		writeError(w, fmt.Errorf("invalid id: %v", err), http.StatusBadRequest)
		return
	}
	postID := int32(postIDInt)

	// Determine if user authenticated
	authenticated, err := s.Authenticate(r)
	if err != nil {
		writeError(w, err, http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case "GET":
		// Make service RPC calls for blog post and comments
		postCtx, postCancel := context.WithTimeout(context.Background(), time.Second)
		defer postCancel()
		commentCtx, commentCancel := context.WithTimeout(context.Background(), time.Second)
		defer commentCancel()

		postChan := make(chan pb.ViewResponse)
		commentChan := make(chan cpb.ListResponse)
		errChan := make(chan error)

		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			resp, err := s.Post.View(postCtx, &pb.ViewQuery{Id: postID, IncludeDeleted: authenticated})
			if err != nil {
				errChan <- err
				return
			}
			postChan <- *resp
		}()
		go func() {
			defer wg.Done()
			resp, err := s.Comment.List(commentCtx, &cpb.ListQuery{PostId: postID})
			if err != nil {
				errChan <- err
				return
			}
			commentChan <- *resp
		}()

		// Wait for both post and comment query to return
		var post pb.ViewResponse
		var comments cpb.ListResponse
		for i := 0; i < 2; i++ {
			select {
			case post = <-postChan:
				continue
			case comments = <-commentChan:
				continue
			case err := <-errChan:
				writeError(w, fmt.Errorf("error retrieving post or its comments: %v", err), http.StatusInternalServerError)
				return
			}
		}
		wg.Wait()

		// Parse responses into JSON structs
		type commentJSON struct {
			Body      string `json:"body,omitempty"`
			Commenter string `json:"commenter,omitempty"`
		}
		commentsJSON := []*commentJSON{}
		for _, comment := range comments.Comments {
			commentsJSON = append(commentsJSON, &commentJSON{Body: comment.Body, Commenter: comment.Commenter})
		}

		type viewResponseJSON struct {
			ID          int            `json:"id"`
			Title       string         `json:"title,omitempty"`
			Body        string         `json:"body,omitempty"`
			Author      string         `json:"author,omitempty"`
			CreatedTime time.Time      `json:"createdTime,omitempty"`
			Deleted     bool           `json:"deleted"`
			Comments    []*commentJSON `json:"comments"`
		}

		// Build API response
		apiResp, err := json.Marshal(&viewResponseJSON{
			ID:          int(post.Id),
			Title:       post.Title,
			Body:        post.Body,
			Author:      post.Author,
			CreatedTime: time.Unix(post.CreatedTime, 0),
			Deleted:     post.Deleted,
			Comments:    commentsJSON,
		})
		if err != nil {
			writeError(w, fmt.Errorf("error processing response %v", err), http.StatusInternalServerError)
			return
		}
		writeResponse(w, r, apiResp, 200)
		return
	case "POST":
		if !authenticated {
			writeError(w, fmt.Errorf("not authenticated"), http.StatusUnauthorized)
			return
		}

		// Parse request data
		rawPostJSON, err := ioutil.ReadAll(r.Body)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
		r.Body.Close()

		var editPost struct {
			ID      int     `json:"id"`
			Title   *string `json:"title"`
			Body    *string `json:"body"`
			Deleted *bool   `json:"deleted"`
		}
		if err := json.Unmarshal(rawPostJSON, &editPost); err != nil {
			writeError(w, fmt.Errorf("error parsing request: %v", err), http.StatusBadRequest)
			return
		}

		// Edit body or title
		if editPost.Title != nil || editPost.Body != nil {
			// Make edit service RPC call
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			resp, err := s.Post.Edit(ctx, &pb.EditCommand{Id: postID, Title: *editPost.Title, Body: *editPost.Body})
			if err != nil {
				writeError(w, fmt.Errorf("internal service error: %v", err), http.StatusInternalServerError)
				return
			}
			if !resp.GetSuccess() {
				writeError(w, fmt.Errorf("unable to update post"), http.StatusInternalServerError)
				return
			}
		}

		// Restore
		if editPost.Deleted != nil && *editPost.Deleted == false {
			// Make restore service RPC call
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			resp, err := s.Post.Restore(ctx, &pb.RestoreCommand{Id: postID})
			if err != nil {
				writeError(w, fmt.Errorf("internal service error: %v", err), http.StatusInternalServerError)
				return
			}
			if !resp.GetSuccess() {
				writeError(w, fmt.Errorf("unable to restore post"), http.StatusInternalServerError)
				return
			}
		}

		writeEmpty(w, r)
		return
	case "DELETE":
		if !authenticated {
			writeError(w, fmt.Errorf("not authenticated"), http.StatusUnauthorized)
			return
		}

		// Make service RPC call
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		resp, err := s.Post.Delete(ctx, &pb.DeleteCommand{Id: postID})
		if err != nil {
			writeError(w, fmt.Errorf("internal service error: %v", err), http.StatusInternalServerError)
			return
		}
		if !resp.GetSuccess() {
			writeError(w, fmt.Errorf("unable to delete post"), http.StatusInternalServerError)
			return
		}
		writeEmpty(w, r)
		return
	}

	writeError(w, fmt.Errorf("invalid method"), http.StatusBadRequest)
}
