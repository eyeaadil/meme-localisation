package routes

import (
	"fmt"
	"meme/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine) {
	// Define all user-related routes

	fmt.Println("Registering routes...aaaaaaaaaaaaaa")
	userGroup := router.Group("/api/users")
	{
		userGroup.POST("/register", controllers.Register)
		userGroup.POST("/login", controllers.Login)
		userGroup.POST("/logout", controllers.Logout)
		userGroup.GET("/profile", controllers.Profile)
	}
}
