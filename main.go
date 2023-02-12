package main

import (
	"fmt"
	"log"

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
	cat := mm.GetCategories()

	m, err := mm.GetMatchup(cat[0].Name)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(m)
	}

	res, _ := mm.GetCategoryBoards(cat[0].Name)

	fmt.Println(res)
}
