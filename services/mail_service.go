package services

import (
	"fmt"
	"os"
	"strconv"

	gomail "gopkg.in/mail.v2"
)

func SendMail(to string, subject string, body string) error {
	email := os.Getenv("EMAIL")
	fmt.Println("Email:", email)
	password := os.Getenv("EMAIL_PASSWORD")
	fmt.Println("Password:", password)
	portStr := os.Getenv("EMAIL_PORT")
	fmt.Println("Port:", portStr)
	port, err := strconv.Atoi(portStr)
	if err != nil {
		fmt.Println("Invalid port:", err)
		return err
	}
	host := os.Getenv("EMAIL_HOST")
	message := gomail.NewMessage()
	message.SetHeader("From", email)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", body)

	dailer := gomail.NewDialer(host, port, email, password)
	if err := dailer.DialAndSend(message); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Email sent successfully")
	return nil
}
