package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Ali-Gorgani/lenslocked/models"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(err)
	}
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	es := models.NewEmailService(
		models.SMTPConfig{
			Host:     host,
			Port:     port,
			Username: username,
			Password: password,
		},
	)

	err = es.ForgotPassword("3JpjL@example.com", "http://localhost:8080/reset")
	if err != nil {
		panic(err)
	}
	fmt.Println("Email sent!")
}
