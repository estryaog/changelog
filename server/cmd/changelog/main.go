package main

import (
	"github.com/estryaog/changelog/internal/server"
)

func main() {
	server := server.NewServer()

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}
