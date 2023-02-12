package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/samuelhegner/best-things/matchupManager"
)

func init() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	mm := matchupManager.NewMatchupManager()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/categories", func(ctx *gin.Context) {
		cat := mm.GetCategories()
		ctx.JSON(200, cat)
	})

	r.GET("/matchup", func(ctx *gin.Context) {
		fmt.Println("lil test")
		fmt.Println(ctx.Query("category"))
		fmt.Println("lil test 2")

		//m, err := mm.GetMatchup()
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
