package services

import "fmt"

type EmailSenderService struct {
}

type EmailSenderServiceInterface interface {
	SendEmail(email string, subject string, body string) error
}

func SendEmail(email string, subject string, body string) error {
	fmt.Println("Mocked email sent:")
	fmt.Println("To:", email)
	fmt.Println("Subject:", subject)
	fmt.Println("Body:", body)
	return nil
}
