package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Run configures and starts the API server
func Run(apiPort int, db *sql.DB) {
	// api
	baseURL := "/api/v1/"

	http.HandleFunc(baseURL, func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%v: %v\n", r.Method, r.URL)

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
			writeError(w, fmt.Errorf("Invalid request"), http.StatusBadRequest)
			return
		}

		switch pieces[0] {
		case "post":
			switch {
			case len(pieces) == 1 || pieces[1] == "":
				postsHandler(w, r, db)
			case len(pieces) == 2 || pieces[2] == "":
				postHandler(w, r, pieces[1], db)
			case (len(pieces) == 3 || pieces[3] == "") && pieces[2] == "comment":
				commentHandler(w, r, pieces[1], db)
			default:
				writeError(w, fmt.Errorf("Invalid request"), http.StatusBadRequest)
				return
			}
		case "authenticate":
			if len(pieces) == 1 {
				authenticateHandler(w, r)
			} else {
				writeError(w, fmt.Errorf("Not authenticated"), http.StatusUnauthorized)
				return
			}
		default:
			writeError(w, fmt.Errorf("Unknown resource"), http.StatusForbidden)
			return
		}
	})
	fmt.Printf("starting server on port %d\n", apiPort)
	http.ListenAndServe(fmt.Sprintf(":%d", apiPort), nil)
}

func writeResponse(w http.ResponseWriter, r *http.Request, handler func() ([]byte, error)) {

	o, err := handler()
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	if len(o) > 0 {
		w.Header().Add("access-control-allow-origin", "*")
		w.Header().Add("content-type", "application/json")
		w.Write(o)
		return
	}
	w.Header().Add("access-control-allow-origin", "*")
	w.WriteHeader(204)
	return
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
		w.Write([]byte("Internal error"))
	}
	w.Write(errorBytes)
}
