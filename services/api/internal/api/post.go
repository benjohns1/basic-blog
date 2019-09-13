package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// PostFull blog post with a post body
type PostFull struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	CreatedTime time.Time `json:"createdTime"`
	Body        string    `json:"body"`
	Author      string    `json:"author"`
	Deleted     bool      `json:"deleted"`
	Comments    []Comment `json:"comments"`
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

func postsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "GET":
		authenticated := authenticate(r)
		writeResponse(w, r, getPosts(db, authenticated))
	case "POST":
		if !authenticate(r) {
			writeError(w, fmt.Errorf("Not authenticated"), http.StatusUnauthorized)
			return
		}

		postData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
		r.Body.Close()
		writeResponse(w, r, newPost(postData, db))
	default:
		writeError(w, fmt.Errorf("Invalid method"), http.StatusBadRequest)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request, postID string, db *sql.DB) {
	ID, err := strconv.Atoi(postID)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	authenticated := authenticate(r)

	switch r.Method {
	case "GET":
		writeResponse(w, r, getPost(ID, db, authenticated))
	case "DELETE":
		if !authenticate(r) {
			writeError(w, fmt.Errorf("Not authenticated"), http.StatusUnauthorized)
			return
		}
		writeResponse(w, r, deletePost(ID, db))
	case "POST":
		if !authenticate(r) {
			writeError(w, fmt.Errorf("Not authenticated"), http.StatusUnauthorized)
			return
		}
		postData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
		r.Body.Close()
		writeResponse(w, r, updatePost(ID, postData, db))
	default:
		writeError(w, fmt.Errorf("Invalid method"), http.StatusBadRequest)
	}
}

func addFieldToUpdate(fields *[]string, values *[]interface{}, name string, value interface{}) {
	*fields = append(*fields, fmt.Sprintf("%s = $%d", name, len(*values)+1))
	*values = append(*values, value)
}

func updatePost(ID int, postData []byte, db *sql.DB) func() ([]byte, error) {
	return func() ([]byte, error) {
		post := PostUpdate{}
		if err := json.Unmarshal(postData, &post); err != nil {
			return nil, err
		}

		// Only update fields that were sent (e.g. not zero-value after unmarshaling)
		var fields []string
		var values []interface{}
		values = append(values, ID)
		if post.Title != nil {
			addFieldToUpdate(&fields, &values, "title", *post.Title)
		}
		if post.CreatedTime != nil {
			addFieldToUpdate(&fields, &values, "created_time", *post.CreatedTime)
		}
		if post.Body != nil {
			addFieldToUpdate(&fields, &values, "body", *post.Body)
		}
		if post.Author != nil {
			addFieldToUpdate(&fields, &values, "author_id", *post.Author)
		}
		if post.Deleted != nil {
			addFieldToUpdate(&fields, &values, "deleted", *post.Deleted)
		}

		query := fmt.Sprintf("UPDATE post SET %v WHERE id = $1", strings.Join(fields, ", "))
		result, err := db.Exec(query, values...)
		if err != nil {
			return nil, err
		}
		rowsAffected, err := result.RowsAffected()
		if rowsAffected != 1 || err != nil {
			return nil, fmt.Errorf("Unable to update post %v", ID)
		}
		return []byte{}, nil
	}
}

func deletePost(postID int, db *sql.DB) func() ([]byte, error) {
	return func() ([]byte, error) {
		result, err := db.Exec("UPDATE post SET deleted = TRUE WHERE id = $1", postID)
		if err != nil {
			return nil, err
		}
		rowsAffected, err := result.RowsAffected()
		if rowsAffected != 1 || err != nil {
			return nil, fmt.Errorf("Unable to delete post %v", postID)
		}
		return []byte{}, nil
	}
}

func getPost(postID int, db *sql.DB, authenticated bool) func() ([]byte, error) {
	return func() ([]byte, error) {
		authCondition := ""
		if !authenticated {
			authCondition = " AND deleted IS NOT TRUE"
		}
		query := fmt.Sprintf(`SELECT id, title, body, created_time, author_id, deleted FROM post WHERE id = $1%v`, authCondition)
		row := db.QueryRow(query, postID)
		var parse struct {
			createdTime *string
			deleted     *bool
		}
		post := PostFull{}
		err := row.Scan(&post.ID, &post.Title, &post.Body, &parse.createdTime, &post.Author, &parse.deleted)
		if err != nil {
			return nil, err
		}
		if t, err := time.Parse(time.RFC3339Nano, *parse.createdTime); err == nil {
			post.CreatedTime = t
		}
		if parse.deleted != nil {
			post.Deleted = *parse.deleted
		}
		if comments, err := scanComments(postID, db); err == nil {
			post.Comments = comments
		} else {
			return nil, err
		}

		return json.Marshal(post)
	}
}

func getPosts(db *sql.DB, authenticated bool) func() ([]byte, error) {
	return func() ([]byte, error) {

		where := ""
		if !authenticated {
			where = " WHERE deleted IS NOT TRUE"
		}
		query := fmt.Sprintf(`SELECT id, title, created_time, author_id, deleted FROM post%v ORDER BY created_time DESC`, where)
		rows, err := db.Query(query)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		posts := []Post{}
		for rows.Next() {
			var parse struct {
				createdTime *string
				deleted     *bool
			}

			post := Post{}
			err := rows.Scan(&post.ID, &post.Title, &parse.createdTime, &post.Author, &parse.deleted)
			if err != nil {
				return nil, err
			}
			if t, err := time.Parse(time.RFC3339Nano, *parse.createdTime); err == nil {
				post.CreatedTime = t
			}
			if parse.deleted != nil {
				post.Deleted = *parse.deleted
			}
			posts = append(posts, post)
		}

		return json.Marshal(posts)
	}
}

func newPost(postData []byte, db *sql.DB) func() ([]byte, error) {
	return func() ([]byte, error) {
		post := PostFull{}
		if err := json.Unmarshal(postData, &post); err != nil {
			return nil, err
		}

		var id int
		err := db.QueryRow(`INSERT INTO post (title, body, created_time, author_id) VALUES ($1, $2, $3, $4) RETURNING id`, post.Title, post.Body, time.Now(), "author").Scan(&id)
		if err != nil {
			return nil, err
		}

		post.ID = id
		return json.Marshal(post)
	}
}
