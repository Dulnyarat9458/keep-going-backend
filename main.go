package main

import (
	"keep_going/controllers"
	database "keep_going/databases"

	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDatabase()
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello from Gin! love"})
	})

	r.GET("/users", controllers.GetUsers)
	r.POST("/signup", controllers.SignUp)

	r.Run()
}
