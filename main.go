package main

import (
	"gogam/api"
	"gogam/database"
	"github.com/gin-gonic/gin"
)

func main() {
	db := database.InitDB()
	defer db.Close()

	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(func(c *gin.Context){
		c.Set("db",db)
		c.Set("hmacSecret",[]byte("asdassasdsdsdswew"))
	})
	 
	api.ApplyRoutes(r)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() 
}