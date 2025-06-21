package main

import (
	"fmt"
	"keep_going/controllers"
	"keep_going/databases"
	"keep_going/middlewares"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	databases.ConnectDatabase()
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("FRONTEND_WEB")},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/", func(c *gin.Context) {
		fmt.Println(os.Getenv("FRONTEND_WEB"))
		c.JSON(200, gin.H{"message": "Hello from Gin! love"})
	})

	r.POST("/signup", controllers.SignUp)
	r.POST("/signin", controllers.SignIn)
	r.POST("/signout", controllers.SignOut)
	r.POST("/forget-password", controllers.ForgetPassword)
	r.POST("/reset-password", controllers.ResetPassword)

	r.GET("/habits", middlewares.Authenticate(), controllers.HabitList)
	r.GET("/habits/:id", middlewares.Authenticate(), controllers.HabitDetail)
	r.PATCH("/habits/:id", middlewares.Authenticate(), controllers.HabitEdit)
	r.POST("/habits/:id/reset", middlewares.Authenticate(), controllers.HabitReset)
	r.DELETE("/habits/:id", middlewares.Authenticate(), controllers.HabitDelete)
	r.POST("/habits", middlewares.Authenticate(), controllers.AddNewHabit)

	r.Run()
}
