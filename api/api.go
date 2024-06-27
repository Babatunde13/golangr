package api

import (
	"bkoiki950/go-store/api/config"
	"bkoiki950/go-store/api/database"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Run () {
	config.GetConfig()
	PORT := config.GetEnv("PORT")
	DB_URI := config.GetEnv("DB_URI")
	database.ConnectDatabase(DB_URI)
	dbClient, err := database.GetActiveClient(); if err != nil {
		fmt.Println(err)
	}

	defer database.CloseDatabase(dbClient)
	if PORT == "" {
		PORT = "8080"
	}
	
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to our API",
		})
	})

	Router(r)

	r.Run(fmt.Sprintf(":%s", PORT))
}
