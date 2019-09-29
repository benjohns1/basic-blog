package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Config configures API
type Config struct {
	APIPort  int
	Services *Services
}

// Services concrete service implementations
type Services struct {
	Authentication AuthenticationService
	Post           PostService
	Comment        CommentService
}

const baseURL = "/api/v1/"

// Start configures and starts the API server
func Start(cfg Config) {
	// api
	http.HandleFunc(baseURL, func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v: %v\n", r.Method, r.URL)

		if r.Method == "OPTIONS" {
			// Allow CORS for localhost (NOT FOR PRODUCTION!)
			w.Header().Add("access-control-allow-origin", "http://localhost:8080")
			w.Header().Add("access-control-allow-methods", "POST, GET, OPTIONS, DELETE")
			w.Header().Add("access-control-allow-headers", "content-type, accept, authorization")
			w.Header().Add("access-control-max-age", "1728000")
			w.WriteHeader(200)
			return
		}

		pieces := strings.Split(r.URL.Path[len(baseURL):], "/")

		if len(pieces) <= 0 {
			writeError(w, fmt.Errorf("invalid request"), http.StatusBadRequest)
			return
		}

		switch pieces[0] {
		case "post":
			switch {
			case len(pieces) == 1 || pieces[1] == "":
				cfg.Services.Post.PostsHandler(w, r)
				return
			case len(pieces) == 2 || pieces[2] == "":
				cfg.Services.Post.PostHandler(w, r, pieces[1])
				return
			case (len(pieces) == 3 || pieces[3] == "") && pieces[2] == "comment":
				cfg.Services.Comment.CommentsHandler(w, r, pieces[1])
				return
			default:
				writeError(w, fmt.Errorf("invalid request"), http.StatusBadRequest)
				return
			}
		case "authenticate":
			if len(pieces) == 1 {
				cfg.Services.Authentication.LoginHandler(w, r)
				return
			}
			writeError(w, fmt.Errorf("invalid request"), http.StatusBadRequest)
			return
		default:
			writeError(w, fmt.Errorf("unknown resource"), http.StatusBadRequest)
			return
		}
	})
	log.Printf("starting API gateway on port %d\n", cfg.APIPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.APIPort), nil)
	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}

func writeEmpty(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("access-control-allow-origin", "*")
	w.WriteHeader(204)
}

func writeResponse(w http.ResponseWriter, r *http.Request, response []byte, status int) {
	w.Header().Add("access-control-allow-origin", "*")
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

// ErrorResponse JSON wrapper for error string
type ErrorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, err error, errorCode int) {
	w.Header().Add("access-control-allow-origin", "*")
	w.Header().Add("content-type", "application/json")

	w.WriteHeader(errorCode)
	errorBytes, err := json.Marshal(ErrorResponse{err.Error()})
	if err != nil {
		w.Write([]byte("internal error"))
	}
	w.Write(errorBytes)
}
