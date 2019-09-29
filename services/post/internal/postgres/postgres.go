package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// DBConn DB connection details
type DBConn struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	AppName  string
}

const (
	maxRetryAttempts = 100
	schema           = "post"
	postTableName    = "post"
)

// Setup sets up db connection and seeds dummy data
func Setup(dbconn DBConn) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%v port=%v user='%v' password='%v' dbname=%v application_name=%v sslmode=disable", dbconn.Host, dbconn.Port, dbconn.User, dbconn.Password, dbconn.DBName, dbconn.AppName))
	if err != nil {
		return nil, err
	}
	// retry if DB isn't available yet
	for attempts := 1; attempts <= maxRetryAttempts; attempts++ {
		err := db.Ping()
		if err == nil {
			log.Printf("db connection attempt %v/%v succeeded\n", attempts, maxRetryAttempts)
			break
		}
		log.Printf("db connection attempt %v/%v: %v\n", attempts, maxRetryAttempts, err)
		time.Sleep(time.Duration(5 * time.Second))
	}

	// initial db setup
	query := fmt.Sprintf(`SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = '%v' AND table_name = '%v');`, schema, postTableName)
	var isSetup bool
	err = db.QueryRow(query).Scan(&isSetup)
	if err != nil {
		return nil, err
	}
	if !isSetup {
		// create tables
		log.Println("performing first-time DB setup")
		db.Exec(fmt.Sprintf(`
			CREATE SCHEMA %v
				CREATE TABLE %v (
					id SERIAL PRIMARY KEY,
					title character varying(64),
					body character varying(1024),
					created_time TIMESTAMPTZ,
					author_id character varying(64),
					deleted boolean
				);
			SET timezone = 'GMT';
		`, schema, postTableName))

		// dummy seed data
		var id int
		db.QueryRow(`INSERT INTO post.post (title, body, created_time, author_id) VALUES ($1, $2, $3, $4) RETURNING id`, "Clean Architecture", "<p>Post body html</p>", time.Now(), "Robert C. Martin").Scan(&id)
		db.QueryRow(`INSERT INTO post.post (title, body, created_time, author_id, deleted) VALUES ($1, $2, $3, $4, $5)`, "Implementing Domain Driven Design", "<p>Post body html</p>", time.Now(), "Vaughn Vernon", true)
	}
	return db, nil
}
