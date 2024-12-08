package main

import (
	"log"
	"meme/config"
	"meme/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.DefaultWriter = os.Stdout
	log.SetOutput(os.Stderr) // Force Go's log to use stderr

	// Connect to MongoDB
	config.ConnectDB()

	// Set Gin to release mode (use gin.DebugMode during development)
	gin.SetMode(gin.ReleaseMode)

	// Initialize Gin router
	router := gin.Default()

	// Add Logger and Recovery middleware
	router.Use(gin.Logger())    // Logs all requests
	router.Use(gin.Recovery()) // Recovers from panics

	// Register user routes
	routes.RegisterUserRoutes(router)
	// routes.RegisterMemeRoutes(router)

	// Define the server address
	serverAddress := ":8080"

	// Start the server and log the success message only if no error occurs
	err := router.Run(serverAddress)
	if err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	} else {
		log.Printf("🚀 Server is running at http://localhost%s", serverAddress)
	}

	
}
