package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func commentHandler(w http.ResponseWriter, r *http.Request, postID string, db *sql.DB) {
	switch r.Method {
	case "POST":
		ID, err := strconv.Atoi(postID)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
		commentData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
		r.Body.Close()
		writeResponse(w, r, newComment(ID, commentData, db))
	default:
		writeError(w, fmt.Errorf("Invalid method"), http.StatusBadRequest)
	}
}

// Comment blog post comment
type Comment struct {
	ID          int       `json:"id"`
	Body        string    `json:"body"`
	CreatedTime time.Time `json:"createdTime"`
	Commenter   string    `json:"commenter"`
}

func scanComments(postID int, db *sql.DB) ([]Comment, error) {

	rows, err := db.Query(`SELECT id, body, created_time, commenter_id FROM comment WHERE post_id = $1 ORDER BY created_time DESC`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}
	for rows.Next() {
		var parse struct {
			createdTime *string
		}

		comment := Comment{}
		err := rows.Scan(&comment.ID, &comment.Body, &parse.createdTime, &comment.Commenter)
		if err != nil {
			return nil, err
		}
		if t, err := time.Parse(time.RFC3339Nano, *parse.createdTime); err == nil {
			comment.CreatedTime = t
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func newComment(postID int, commentData []byte, db *sql.DB) func() ([]byte, error) {
	return func() ([]byte, error) {
		comment := Comment{}
		if err := json.Unmarshal(commentData, &comment); err != nil {
			return nil, err
		}

		var id int
		err := db.QueryRow(`INSERT INTO comment (body, created_time, post_id, commenter_id) VALUES ($1, $2, $3, $4) RETURNING id`, comment.Body, time.Now(), postID, comment.Commenter).Scan(&id)
		if err != nil {
			return nil, err
		}

		comment.ID = id
		return json.Marshal(comment)
	}
}
