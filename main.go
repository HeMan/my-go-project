package main

import (
	"fmt"
	"log"

	"my-go-project/database"

	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
)

func main() {
	fmt.Println("Hello, World!")

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// Initialize the database
	database.Init()

	// Define a simple fasthttp route
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/ping":
			ctx.SetStatusCode(fasthttp.StatusOK)
			ctx.SetBodyString("pong")
		default:
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.SetBodyString("Not Found")
		}
	}

	// Start the fasthttp server
	log.Println("Starting server on :8080")
	if err := fasthttp.ListenAndServe(":8080", requestHandler); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
