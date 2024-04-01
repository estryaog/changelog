package database

import (
	"database/sql"
	"fmt"
	"os"

	store "github.com/estryaog/changelog/database"
)

type Service interface {
}

type service struct {
	db *sql.DB
}

var (
	database = os.Getenv("PSQL_DATABASE")
	password = os.Getenv("PSQL_PASSWORD")
	username = os.Getenv("PSQL_USERNAME")
	port     = os.Getenv("PSQL_PORT")
	host     = os.Getenv("PSQL_HOST")
)

func New() Service {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		panic(err)
	}

	store.NewPostgresUserStore(db)

	s := &service{db: db}
	return s
}
