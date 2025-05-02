package main

import (
	"keep_going/controllers"
	"keep_going/databases"
	"keep_going/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	databases.ConnectDatabase()
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello from Gin! love"})
	})

	r.POST("/signup", controllers.SignUp)
	r.POST("/signin", controllers.SignIn)
	r.POST("/signout", controllers.SignOut)
	r.POST("/forget-password", controllers.ForgetPassword)
	r.POST("/reset-password", controllers.ResetPassword)

	r.GET("/habits", middlewares.Authenticate(), controllers.HabitList)
	r.POST("/habits", middlewares.Authenticate(), controllers.AddNewHabit)

	r.Run()
}
