package services

import (
	"fmt"

	gomail "gopkg.in/mail.v2"
)

func SendMail(to string, subject string, body string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", "lakshya1234goel@gmail.com")
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", body)

	dailer := gomail.NewDialer("smtp.gmail.com", 587, "lakshya1234goel@gmail.com", "Lakshya@123")
	if err := dailer.DialAndSend(message); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Email sent successfully")
	return nil
}
