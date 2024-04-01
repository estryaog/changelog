package main

import (
	"github.com/estryaog/changelog/internal/database"
	"github.com/estryaog/changelog/internal/server"
)

func main() {
	db := database.New()
	server := server.NewServer(db)

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}
