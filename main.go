package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"lottery/config"
	"lottery/database"
	routers "lottery/router"
	"lottery/service"

	"github.com/joho/godotenv"
	// "github.com/aws/aws-lambda-go/events"
	// "github.com/aws/aws-lambda-go/lambda"
	// ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
)

func main() {
	cfg := config.Load()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	redisClient := database.ConnectRedis()

	go func() {
		time.Sleep(1 * time.Second)
		for {
			var date string
			fmt.Print("Lottery date (YYYY-MM-DD): ")
			fmt.Scanln(&date)
			if date == "" {
				continue
			}
			lotteryService := service.NewLotteryService(db, redisClient, cfg)
			lotteryService.Add(date)
		}

	}()

	// r := routes.Setup(db, cfg)
	r := routers.Setup(db, redisClient)

	port := os.Getenv("PORT")
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
