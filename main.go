package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/samuelhegner/best-things/matchupManager"
	"github.com/samuelhegner/best-things/types"
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
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/categories", func(ctx *gin.Context) {
		cat := mm.GetCategories()
		ctx.JSON(http.StatusOK, cat)
	})

	r.GET("/matchup", func(ctx *gin.Context) {
		cat := ctx.Query("category")

		if cat == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "provide a category in params"})
			return
		}

		matchup, err := mm.GetMatchup(cat)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "issue getting matchup: " + err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, matchup)
	})

	r.POST("/matchup", func(ctx *gin.Context) {
		var json types.MatchupSubmit

		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		success, err := mm.SubmitMatchupResponse(json.Guid, json.Winner, json.Category)

		if err != nil || !success {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "issue submitting matchup: " + err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"success": "true"})
	})

	r.GET("/leaderboard", func(ctx *gin.Context) {
		cat := ctx.Query("category")

		if cat == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "provide a category in params"})
			return
		}

		res, err := mm.GetCategoryBoards(cat)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "issue getting board: " + err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, res)
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
