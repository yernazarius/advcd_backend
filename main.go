package main

import (
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "github.com/marcojulian/go-jwt/controllers"
    "github.com/marcojulian/go-jwt/initializers"
    "github.com/marcojulian/go-jwt/middleware"
)

func init() {
	initializers.LoanEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Authorization", "Content-Type", "Accept"}
	r.Use(cors.New(config))


	r.Static("/uploads", "./uploads")

	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.POST("/add-word", middleware.RequireAuth, controllers.AddWord)
	r.GET("/words", controllers.GetWords) 
	r.DELETE("/words/:id", middleware.RequireAuth, controllers.DeleteWord)
	r.PATCH("/words/:id", middleware.RequireAuth, controllers.UpdateWord) 

	r.Run()
}
