package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// DBConn DB connection details
type DBConn struct {
	host     string
	password string
	port     int
}

func main() {
	// db connection
	dbconn := DBConn{
		host:     "localhost",
		password: "asdf1234",
		port:     5432,
	}
	if v, exists := os.LookupEnv("DB_HOST"); exists {
		dbconn.host = v
	}
	if v, exists := os.LookupEnv("DB_PASSWORD"); exists {
		dbconn.password = v
	}
	if v, exists := os.LookupEnv("DB_PORT"); exists {
		port, err := strconv.Atoi(v)
		if err == nil {
			dbconn.port = port
		}
	}

	db, err := sql.Open("postgres", fmt.Sprintf("host=%v port=%v user='postgres' password='%v' dbname=postgres application_name=blog sslmode=disable", dbconn.host, dbconn.port, dbconn.password))
	if err != nil {
		fmt.Println(err)
		return
	}

	// initial db setup
	var isSetup bool
	err = db.QueryRow(`SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'post');`).Scan(&isSetup)
	if err != nil {
		maxAttempts := 100
		for attempts := 1; attempts <= maxAttempts; attempts++ {
			fmt.Printf("Attempt %v/%v: %v", attempts, maxAttempts, err)
			time.Sleep(time.Duration(5 * time.Second))
			err = db.Ping()
			if err == nil {
				break
			}
		}
	}
	if !isSetup {
		// create tables
		fmt.Println("performing first-time DB setup")
		db.Exec(`
			CREATE TABLE post (
				id SERIAL PRIMARY KEY,
				title character varying(64),
				body character varying(1024),
				created_time TIMESTAMPTZ,
				author_id character varying(64),
				deleted boolean
			);
			CREATE TABLE comment (
				id SERIAL PRIMARY KEY,
				body character varying(1024),
				created_time TIMESTAMPTZ,
				post_id integer REFERENCES post(id) ON DELETE CASCADE ON UPDATE CASCADE,
				commenter_id character varying(64)
			);
			SET timezone = 'GMT'
		`)

		// dummy seed data
		var id int
		db.QueryRow(`INSERT INTO post (title, body, created_time, author_id) VALUES ($1, $2, $3, $4) RETURNING id`, "Clean Architecture", "<p>Post body html</p>", time.Now(), "Robert C. Martin").Scan(&id)
		db.QueryRow(`INSERT INTO comment (body, created_time, post_id, commenter_id) VALUES ($1, $2, $3, $4)`, "Pulsara is awesome!", time.Now(), id, "commenter name")
		db.QueryRow(`INSERT INTO post (title, body, created_time, author_id, deleted) VALUES ($1, $2, $3, $4, $5)`, "Implementing Domain Driven Design", "<p>Post body html</p>", time.Now(), "Vaughn Vernon", true)
	}

	// api
	baseURL := "/api/v1/"

	http.HandleFunc(baseURL, func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%v: %v\n", r.Method, r.URL)

		pieces := strings.Split(r.URL.Path[len(baseURL):], "/")

		if len(pieces) <= 0 {
			w.Write(nil)
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
				w.WriteHeader(404)
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("static app request for %v\n", r.RemoteAddr)
		http.ServeFile(w, r, "index.html")
	})
	fmt.Println("starting server")
	http.ListenAndServe(":8080", nil)
}

func writeJSONHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	w.Header().Add("access-control-allow-origin", "*")
}

func writeResponse(w http.ResponseWriter, r *http.Request, handler func() ([]byte, error)) {
	writeJSONHeaders(w, r)

	o, err := handler()
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	if len(o) > 0 {
		w.Write(o)
		return
	}
	w.WriteHeader(204)
	return
}

// ErrorResponse JSON wrapper for error string
type ErrorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, err error, errorCode int) {
	w.WriteHeader(errorCode)
	errorBytes, err := json.Marshal(ErrorResponse{err.Error()})
	if err != nil {
		w.Write([]byte("Internal error"))
	}
	w.Write(errorBytes)
}

// Credentials user credentials
type Credentials struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func authenticateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	credentialData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	creds := Credentials{}
	if err := json.Unmarshal(credentialData, &creds); err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	resp, err := getToken(creds)
	if err != nil {
		w.WriteHeader(401)
	}
	w.Write(resp)
}

func getToken(creds Credentials) ([]byte, error) {
	if creds.User == "bobross" && creds.Password == "painter" {
		return []byte("notSoSecureToken!"), nil
	}
	return nil, fmt.Errorf("Invalid credentials")
}

func authenticate(r *http.Request) bool {
	token := r.Header.Get("Authorization")
	return "notSoSecureToken!" == token
}

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
