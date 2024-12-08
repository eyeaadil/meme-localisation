package routes

import (
	"meme/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterMemeRoutes(router *gin.Engine) {
	memeGroup := router.Group("/api/memes")
	{

		println("adilllllllllll")
		memeGroup.POST("/", controllers.UploadMemeWithTextExtraction) // Upload meme with text extraction
		memeGroup.GET("/:id", controllers.GetMeme)                   // Get a single meme by ID
		memeGroup.GET("/", controllers.ListMemes)                    // List memes with filtering and pagination
		memeGroup.PUT("/:id", controllers.UpdateMeme)                // Update a meme by ID
		memeGroup.DELETE("/:id", controllers.DeleteMeme)             // Delete a meme by ID
	}
}
