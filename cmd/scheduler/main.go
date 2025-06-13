package main

import (
	"kronos-scheduler/config"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		println("No .env file found.")
	}

	app := config.CreateServer()

	port := os.Getenv("PORT")

	if port == "" {
		port = "5000"
	}

	if err := app.Listen(":" + port); err != nil {
		panic(err)
	}
}
