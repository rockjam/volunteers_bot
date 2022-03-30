package main

import (
	"dv/models"
	"dv/services"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	time.Sleep(5 * time.Second)
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

	messageService := services.NewMessage(db, os.Getenv("BOT_TOKEN"))

	log.Println("starting http service @ :8080")
	http.HandleFunc("/incoming_message", messageService.IncomingMessageHTTPHandler())
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
