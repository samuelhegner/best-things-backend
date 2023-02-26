package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/gin-gonic/gin"
	"github.com/samuelhegner/best-things/matchupManager"
	"github.com/samuelhegner/best-things/types"
)

func main() {

	secretName := "production_best-things-api"
	region := "eu-west-1"

	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		// For a list of exceptions thrown, see
		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
		log.Fatal(err.Error())
	}

	// Decrypts secret using the associated KMS key.
	secretString := *result.SecretString
	var secretMap map[string]string
	json.Unmarshal([]byte(secretString), &secretMap)

	for key, val := range secretMap {
		os.Setenv(key, val)
	}

	mm := matchupManager.NewMatchupManager()

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

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

		res, err := mm.GetLeaderboards(cat)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "issue getting board: " + err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, res)
	})

	r.Run(":8080")
}
