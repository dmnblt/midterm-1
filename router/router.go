package router

import (
	"github.com/gin-gonic/gin"
	"go-midterm/controllers"
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
