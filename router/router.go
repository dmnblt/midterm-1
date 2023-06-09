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

	products := router.Group("/products")
	{
		products.POST("/add", controllers.AddProduct)
		products.GET("/user-product", controllers.GetProductsByUserId)
		products.GET("/find-filter-product", controllers.GetProductsBetweenPrices)
		products.PUT("/user-product", controllers.UpdateProductOwner)
	}
	comments := router.Group("/comments")
	{
		comments.POST("/add", controllers.AddComment)
		comments.GET("/product-comment", controllers.GetCommentsByProductID)
		comments.GET("/comment-desc", controllers.GetCommentsByLikesDesc)
		comments.GET("/like-comment", controllers.LikeComment)
	}
	return router
}
