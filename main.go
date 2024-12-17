package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"SIDIMASBE/controllers/authcontroller"
	"SIDIMASBE/controllers/productcontroller"
	_ "SIDIMASBE/docs"
	"SIDIMASBE/middlewares"
	"SIDIMASBE/models"
)

// @title GOjawet API
// @version 1.0
// @description This is the API server for the GOjawet application.
// @contact.name Elysia
// @contact.email loveyouelysia@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api
// @schemes http
// @SecurityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	models.ConnectDatabase()
	r := gin.Default()

	r.POST("/login", authcontroller.Login)
	r.POST("/register", authcontroller.Register)
	r.GET("/logout", authcontroller.Logout)

	api := r.Group("/api")
	api.Use(middlewares.JWTVerif())
	api.GET("/products", productcontroller.Getproduct)
	api.POST("/addpr", productcontroller.Addproduct)

	// Serve the Swagger UI at /swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}

