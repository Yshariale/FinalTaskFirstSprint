package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/Yshariale/FinalTaskFirstSprint/agent"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	agent.RunAgent()
}
