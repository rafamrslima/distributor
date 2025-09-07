package core

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/rafamrslima/distributor/internal/db"
	"github.com/rafamrslima/distributor/internal/domain"
	"github.com/rafamrslima/distributor/internal/email"
)

func Handle(message *azservicebus.ReceivedMessage) error {
	var parsedMessage domain.Message
	err := json.Unmarshal(message.Body, &parsedMessage)

	if err != nil {
		log.Println("Error to parse received json. MessageId:", message.MessageID)
		return err
	}

	parsedMessage.MessageReceivedAt = time.Now()
	err = db.SaveReceivedMessages(parsedMessage)

	if err != nil {
		log.Println("Error to save data in the database", err.Error())
		return err
	}

	isValid := validateMessage(parsedMessage)

	if !isValid {
		return fmt.Errorf("message is invalid")
	}

	emailInfo := domain.Message{
		Name:              parsedMessage.Name,
		Email:             parsedMessage.Email,
		EmailCc:           parsedMessage.EmailCc,
		Content:           parsedMessage.Content,
		MessageReceivedAt: parsedMessage.MessageReceivedAt,
	}

	mailErr := email.SendEmail(emailInfo)

	if mailErr != nil {
		return mailErr
	}

	log.Println("Email sent successfully.")
	return nil
}

func validateMessage(message domain.Message) bool {
	if !email.IsValidEmail(message.Email) {
		log.Println("Email is invalid:", message.Email)
		return false
	}

	if message.EmailCc != "" && !email.IsValidEmail(message.Email) {
		log.Println("EmailCc is invalid:", message.EmailCc)
		return false
	}

	if len(message.Content) == 0 {
		log.Println("Content is empty.")
		return false
	}

	return true
}
