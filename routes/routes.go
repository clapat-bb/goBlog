package routes

import (
	"goblog/controllers"
	"goblog/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")

	api.POST("/register", controllers.Register)
	api.POST(("/login"), controllers.Login)
	api.GET("/tags/:name/posts", controllers.GetPostByTag)

	posts := api.Group("/posts")
	{
		posts.GET("", controllers.GetPosts)
		posts.POST("", middlewares.JWTAuthMiddleware(), controllers.CreatePost)
		posts.GET("/:id", controllers.GetPostByID)
		posts.PUT("/:id", middlewares.JWTAuthMiddleware(), controllers.UpdataPost)
		posts.DELETE("/:id", middlewares.JWTAuthMiddleware(), controllers.DeletePost)
	}

	comments := api.Group("/comments")
	{
		comments.POST("", middlewares.JWTAuthMiddleware(), controllers.CreateComment)
		comments.GET("", controllers.GetCommentsByPostID)
	}

	likes := api.Group("/likes")
	{
		likes.POST("", middlewares.JWTAuthMiddleware(), controllers.Like)
		likes.GET("/count", controllers.GetLikeCount)
		likes.GET("/check", middlewares.JWTAuthMiddleware(), controllers.CheckIfLiked)
	}
}
