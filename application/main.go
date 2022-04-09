package main

import (
	"dv/application/handler"
	"dv/models"
	"dv/services"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
		),
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&models.Message{}, &models.MessageHashtag{})
	if err != nil {
		panic(err)
	}

	groupChatID, err := strconv.ParseInt(os.Getenv("GROUP_CHAT_ID"), 10, 64)
	if err != nil {
		panic(err)
	}

	messageService := services.NewMessage(db, os.Getenv("BOT_TOKEN"), os.Getenv("BOT_NAME"), groupChatID)
	eventHandler := handler.NewEventHandler(messageService)
	lambda.Start(eventHandler.HandleRequest)
}
