package main

import (
	"net/http"

	"github.com/yadibolt/obelisk/handlers"
)

func main() {
	var s IServer = &Server{}
	endpoints := map[string]func(writer http.ResponseWriter, request *http.Request){
		"test": handlers.HandleTest,
	}
	server := s.NewServer(endpoints)
	server.Serve()
}