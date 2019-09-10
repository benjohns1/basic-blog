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
}

func main() {
	// db connection
	dbconn := DBConn{
		host:     "localhost",
		password: "asdf1234",
	}
	if v, exists := os.LookupEnv("DB_HOST"); exists {
		dbconn.host = v
	}
	if v, exists := os.LookupEnv("DB_PASSWORD"); exists {
		dbconn.password = v
	}
	db, err := sql.Open("postgres", fmt.Sprintf("host=%v port=5432 user='postgres' password='%v' dbname=postgres application_name=blog sslmode=disable", dbconn.host, dbconn.password))
	if err != nil {
		fmt.Println(err)
		return
	}

	// initial db setup
	var isSetup bool
	err = db.QueryRow(`SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'post');`).Scan(&isSetup)
	if err != nil {
		for attempts := 0; attempts < 100; attempts++ {
			fmt.Println(err)
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
				w.WriteHeader(404)
			}
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
		writeError(w, err)
		return
	}
	w.Write(o)
}

func writeError(w http.ResponseWriter, err error) {
	w.Write([]byte(err.Error()))
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
		writeError(w, err)
		return
	}
	r.Body.Close()

	creds := Credentials{}
	if err := json.Unmarshal(credentialData, &creds); err != nil {
		writeError(w, err)
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
	Comments    []Comment `json:"comments"`
}

// Post blog post without the post body
type Post struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	CreatedTime time.Time `json:"createdTime"`
	Author      string    `json:"author"`
}

func postsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "GET":
		writeResponse(w, r, getPosts(db))
		return
	case "POST":
		if !authenticate(r) {
			w.WriteHeader(401)
			return
		}

		postData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			writeError(w, err)
			return
		}
		r.Body.Close()
		writeResponse(w, r, newPost(postData, db))
		return
	}
}

func postHandler(w http.ResponseWriter, r *http.Request, postID string, db *sql.DB) {
	switch r.Method {
	case "GET":
		ID, err := strconv.Atoi(postID)
		if err != nil {
			writeError(w, err)
			return
		}
		writeResponse(w, r, getPost(ID, db))
		return
	}
}

func commentHandler(w http.ResponseWriter, r *http.Request, postID string, db *sql.DB) {
	switch r.Method {
	case "POST":
		ID, err := strconv.Atoi(postID)
		if err != nil {
			writeError(w, err)
			return
		}
		commentData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			writeError(w, err)
			return
		}
		r.Body.Close()
		writeResponse(w, r, newComment(ID, commentData, db))
	}
}

func getPost(postID int, db *sql.DB) func() ([]byte, error) {
	return func() ([]byte, error) {
		row := db.QueryRow(`SELECT id, title, body, created_time, author_id FROM post WHERE id = $1 AND deleted IS NOT TRUE`, postID)
		var parse struct {
			createdTime *string
		}
		post := PostFull{}
		err := row.Scan(&post.ID, &post.Title, &post.Body, &parse.createdTime, &post.Author)
		if err != nil {
			return nil, err
		}
		if t, err := time.Parse(time.RFC3339Nano, *parse.createdTime); err != nil {
			post.CreatedTime = t
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

func getPosts(db *sql.DB) func() ([]byte, error) {
	return func() ([]byte, error) {

		rows, err := db.Query(`SELECT id, title, created_time, author_id FROM post WHERE deleted IS NOT TRUE ORDER BY created_time DESC`)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		posts := []Post{}
		for rows.Next() {
			var parse struct {
				createdTime *string
			}

			post := Post{}
			err := rows.Scan(&post.ID, &post.Title, &parse.createdTime, &post.Author)
			if err != nil {
				return nil, err
			}
			if t, err := time.Parse(time.RFC3339Nano, *parse.createdTime); err == nil {
				post.CreatedTime = t
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
