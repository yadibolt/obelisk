package handlers

import (
	"log"
	"net/http"
)

func HandleTest(writer http.ResponseWriter, request *http.Request) {
    log.Println("Test handler executed")
}