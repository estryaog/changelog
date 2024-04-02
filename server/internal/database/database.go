package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
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

func (s *ServiceImpl) IsKeyValueExist(table, key, value string) (bool, error) {
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s = $1)", table, key)
	var exists bool
	err := s.db.QueryRow(query, value).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
