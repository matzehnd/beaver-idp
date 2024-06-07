package config

import (
	"database/sql"
)

type Config struct {
	DB *sql.DB
}

func LoadConfig() *Config {
	// Hier kann man die Konfiguration laden (z.B. aus einer Datei oder Umgebungsvariablen)
	db, err := sql.Open("postgres", "user=postgres dbname=yourdb sslmode=disable")
	if err != nil {
		panic(err)
	}
	return &Config{
		DB: db,
	}
}
