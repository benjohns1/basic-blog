package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	baseURL := "/api/v1/"

	http.HandleFunc(baseURL, func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%v: %v\n", r.Method, r.URL)

		pieces := strings.Split(r.URL.Path[len(baseURL):], "/")

		if len(pieces) <= 0 {
			w.Write(nil)
			return
		}

		if pieces[0] == "post" {
			if len(pieces) == 1 || pieces[1] == "" {
				postsHandler(w, r)
			} else {
				postHandler(w, r, pieces[1])
			}
			return
		}
	})
	fmt.Println("starting server")
	http.ListenAndServe(":8080", nil)
}

func writeResponse(w http.ResponseWriter, r *http.Request, handler func() ([]byte, error)) {
	w.Header().Add("content-type", "application/json")
	w.Header().Add("access-control-allow-origin", "*")

	o, err := handler()
	if err != nil {
		writeError(w, err)
		return
	}
	w.Write(o)
}

func writeError(w http.ResponseWriter, err error) {
	w.Write([]byte(err.Error()))
}

// PostBody blog post with a post body
type PostBody struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	CreatedTime time.Time `json:"createdTime"`
	Body        string    `json:"body"`
	Author      string    `json:"author"`
}

// Post blog post without the post body
type Post struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	CreatedTime time.Time `json:"createdTime"`
	Author      string    `json:"author"`
}

func postsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		writeResponse(w, r, getPosts)
		return
	case "POST":
		postData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			writeError(w, err)
			return
		}
		r.Body.Close()
		writeResponse(w, r, newPost(postData))
		return
	}
}

func postHandler(w http.ResponseWriter, r *http.Request, postID string) {
	switch r.Method {
	case "GET":
		ID, err := strconv.Atoi(postID)
		if err != nil {
			writeError(w, err)
			return
		}
		writeResponse(w, r, getPost(ID))
		return
	case "POST":
		//writeResponse(w, r, newPost)
	}
}

func getPost(postID int) func() ([]byte, error) {
	return func() ([]byte, error) {
		post := PostBody{ID: 1, Title: "Hello World", CreatedTime: time.Now(), Body: "post body!", Author: "Robert C. Martin"}
		return json.Marshal(post)
	}
}

func getPosts() ([]byte, error) {
	posts := []Post{
		Post{ID: 1, Title: "Hello World", CreatedTime: time.Now(), Author: "Robert C. Martin"},
		Post{ID: 2, Title: "Implementing Domain Driven Design", CreatedTime: time.Now(), Author: "Vaughn Vernon"},
	}
	return json.Marshal(posts)
}

func newPost(postData []byte) func() ([]byte, error) {
	return func() ([]byte, error) {
		post := PostBody{}
		if err := json.Unmarshal(postData, &post); err != nil {
			return nil, err
		}
		post.ID = 123
		return json.Marshal(post)
	}
}
