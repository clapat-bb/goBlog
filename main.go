// @title GoBlog API文档
// @version 1.1
// @description 展示接口功能
// @host localhost:8080
// @BasePath /api
package main

import (
	"goblog/config"
	"goblog/database"
	"goblog/pkg/cache"
	"goblog/routes"

	_ "goblog/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	config.LoadConfig()

	database.ConnectDatabase()

	cache.InitRedis("localhost", "6379", "", 0)

	r := gin.Default()

	routes.SetupRoutes(r)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
