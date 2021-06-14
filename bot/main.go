package main

import (
	"github.com/joho/godotenv"
	"goatbot/bot"
	"log"
)

func main() {
	log.Println("Loading environment config...")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bot.Init()
	bot.Start()
}
