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
		users.GET("/between-dates", controllers.FindUsersBetweenDates)
		users.PUT("/:id/email", controllers.UpdateUserEmail)
		users.POST("/change/pass", controllers.UpdateUserPassword)
	}

	news := router.Group("/news")
	{
		news.POST("/add", controllers.AddNews)
	}
	// add comment router

	return router
}
