package main

import (
	"os"
	"strconv"

	"github.com/benjohns1/basic-blog/services/api/internal/api"
	"github.com/benjohns1/basic-blog/services/api/internal/postgres"
	_ "github.com/lib/pq"
)

// Config db and api configs
type Config struct {
	DBConn  postgres.DBConn
	APIPort int
}

func main() {

	cfg := loadConfigs()

	db, err := postgres.Setup(cfg.DBConn)
	if err != nil {
		panic(err)
	}

	api.Run(cfg.APIPort, db)
}

func loadConfigs() Config {
	// environment configs
	apiPort := 3000
	if v, exists := os.LookupEnv("API_PORT"); exists {
		port, err := strconv.Atoi(v)
		if err == nil {
			apiPort = port
		}
	}

	// db connection
	dbconn := postgres.DBConn{
		Host:     "localhost",
		Password: "asdf1234",
		Port:     5432,
	}
	if v, exists := os.LookupEnv("DB_HOST"); exists {
		dbconn.Host = v
	}
	if v, exists := os.LookupEnv("DB_PASSWORD"); exists {
		dbconn.Password = v
	}
	if v, exists := os.LookupEnv("DB_PORT"); exists {
		port, err := strconv.Atoi(v)
		if err == nil {
			dbconn.Port = port
		}
	}

	return Config{dbconn, apiPort}
}
