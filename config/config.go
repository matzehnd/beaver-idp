package config

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Config struct {
	DB *sql.DB
}

func LoadConfig(conStr string) *Config {
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		panic(err)
	}
	return &Config{
		DB: db,
	}
}

func InitDb(config Config) {
	config.DB.Exec("CREATE TABLE IF NOT EXISTS events (id SERIAL PRIMARY KEY, type VARCHAR(255), data JSON);")
}
