package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.Handler()
	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	server.ListenAndServe()
}
