package server

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/estryaog/changelog/internal/database"
)

type Server struct {
	port           int
	db             *database.ServiceImpl
	UserStore      database.UserStore
	ChangelogStore database.ChangelogStore
}

func NewServer(db *database.ServiceImpl) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	newServer := &Server{
		port:           port,
		db:             db,
		UserStore:      database.NewPostgresUserStore(db),
		ChangelogStore: database.NewPostgresChangelogStore(db),
	}

	server := &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
