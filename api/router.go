package api

import (
	"user/api/handler"
	"user/api/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @title User
// @version 1.0
// @description API Gateway
// BasePath: /
func Router(hand *handler.Handler) *gin.Engine {
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	auth := router.Group("/auth")
	{
		auth.POST("/register")
		auth.POST("/login")
		auth.POST("/forgot-password")
		auth.POST("/reset-password")
	}

	user := router.Group("/user")
	user.Use(middleware.Check)
	{
		user.POST("/logout")
		user.GET("/profile")
		user.PUT("/profile")
		user.POST("/change-password")
		user.POST("/photo")

	}
	return router
}
