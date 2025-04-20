package main

import (
	"github.com/vmamchur/vacancy-board/internal/server"
)

func main() {
	server := server.NewServer()

	server.ListenAndServe()
}
