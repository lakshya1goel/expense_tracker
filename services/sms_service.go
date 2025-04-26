package services

import (
	"fmt"
	"os"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

func SendSms(to string, otp string) error {
	client := twilio.NewRestClient()
	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(os.Getenv("TWILIO_PHONE_NUMBER"))
	params.SetBody("This is your OTP: " + otp)

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println(err.Error())
		return err
	} 
	fmt.Println("SMS sent successfully!")
	return nil
}
