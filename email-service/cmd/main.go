package main

import (
	"fmt"
	"log"

	"golang/email-service/internal/email"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env variables")
	}

	to := "dilara.galimkizi@gmail.com"
	subject := "from Dilyara"
	body := "Hello from Go Microservice! Aruzhan Ali we love you!"

	err = email.SendEmail(to, subject, body)
	if err != nil {
		log.Fatalf("Error sending email: %v", err)
	}

	fmt.Println("Email sent successfully!")
}
