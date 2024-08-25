package main

import (
	"fmt"
	"github.com/drunkleen/rasta/config"
	_ "github.com/drunkleen/rasta/docs/swagger"
	newsletterroute "github.com/drunkleen/rasta/internal/route/newsletter"
	userroute "github.com/drunkleen/rasta/internal/route/user"
	"github.com/drunkleen/rasta/pkg/database"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Rasta API
// @version 1.0
// @description API for Rasta

// @BasePath /api/v1
func main() {
	config.Init()
	database.InitDB()
	fmt.Printf("\nEnvironment Variables:%+v\n\n", config.GetEnvVars())

	r := gin.Default()
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := r.Group("/api/v1")

	userroute.RegisterUserRoutes(api)
	newsletterroute.RegisterUserRoutes(api)

	if r.Run(":"+config.GetServerPort()) != nil {
		return
	}
}
