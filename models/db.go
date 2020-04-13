package models

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
)

var db *sqlx.DB

func ConnectToDB() error {
	dbUrl := os.Getenv("DATABASE_URL")

	if dbUrl == "" {
		dbUrl = "postgres://uajeskdbvzrjul:4d5807b139549f65484a97b575e0839319dc9ef1945799c0e03ce5a8c0877b2b@ec2-174-129-227-80.compute-1.amazonaws.com:5432/defo6k3unodbg0"
	}

	conn, err := sqlx.Open("postgres", dbUrl)
	if err != nil {
		return err
	}
	db = conn
	return nil
}

func GetDB() *sqlx.DB {
	return db
}

