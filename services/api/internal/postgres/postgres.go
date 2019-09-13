package postgres

import (
	"database/sql"
	"fmt"
	"time"
)

// DBConn DB connection details
type DBConn struct {
	Host     string
	Password string
	Port     int
}

// Setup sets up db connection and seeds dummy data
func Setup(dbconn DBConn) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%v port=%v user='postgres' password='%v' dbname=postgres application_name=blog sslmode=disable", dbconn.Host, dbconn.Port, dbconn.Password))
	if err != nil {
		return nil, err
	}

	// initial db setup
	var isSetup bool
	err = db.QueryRow(`SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'post');`).Scan(&isSetup)
	if err != nil {
		maxAttempts := 100
		for attempts := 1; attempts <= maxAttempts; attempts++ {
			fmt.Printf("Attempt %v/%v: %v\n", attempts, maxAttempts, err)
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
	return db, nil
}
