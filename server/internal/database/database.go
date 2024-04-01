package database

import (
	"database/sql"
	"fmt"
	"os"
)

type Service interface {
	GetDB() *sql.DB
}

type ServiceImpl struct {
	db *sql.DB
}

var (
	database = os.Getenv("PSQL_DATABASE")
	password = os.Getenv("PSQL_PASSWORD")
	username = os.Getenv("PSQL_USERNAME")
	port     = os.Getenv("PSQL_PORT")
	host     = os.Getenv("PSQL_HOST")
)

func New() *ServiceImpl {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		panic(err)
	}

	return &ServiceImpl{db: db}
}

func (s *ServiceImpl) GetDB() *sql.DB {
	return s.db
}