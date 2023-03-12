package router

import (
	"go-midterm/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	users := router.Group("/user")
	{
		users.POST("/registration", controllers.RegisterUser)
		users.POST("/auth", controllers.AuthorizeUser)
		users.GET("/by-name", controllers.SearchUsersByName)
	}

	return router
}
