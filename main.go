package main

import (
	"log"
	"meme/config"
	"net/http"

	"github.com/gorilla/mux"
)


func main() {
	// Connect to MongoDB
	config.ConnectDB()

	router := mux.NewRouter()

	// Register routes
	// routes.RegisterRoutes(router)

	// Start server
	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", router)
}