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
	schema           = "comment"
	commentTableName = "comment"
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
	query := fmt.Sprintf(`SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = '%v' AND table_name = '%v');`, schema, commentTableName)
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
					body character varying(1024),
					created_time TIMESTAMPTZ,
					post_id integer,
					commenter_id character varying(64)
				);
			SET timezone = 'GMT';
		`, schema, commentTableName))

		db.QueryRow(`INSERT INTO comment.comment (body, created_time, post_id, commenter_id) VALUES ($1, $2, $3, $4)`, "Pulsara is awesome!", time.Now(), 1, "Ben Johns")
	}
	return db, nil
}
