package core

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/rafamrslima/distributor/internal/db"
	"github.com/rafamrslima/distributor/internal/domain"
	"github.com/rafamrslima/distributor/internal/email"
)

func Handle(message *azservicebus.ReceivedMessage) error {

	var jsonMsg domain.Message
	err := json.Unmarshal(message.Body, &jsonMsg)

	if err != nil {
		fmt.Println("Error to parse received json. MessageId:", message.MessageID)
		return err
	}

	db.SaveReceivedMessages(jsonMsg)

	isValid := validateMessage(jsonMsg)

	if !isValid {
		return fmt.Errorf("message is invalid")
	}

	emailInfo := domain.Message{
		Name:              jsonMsg.Name,
		Email:             jsonMsg.Email,
		EmailCc:           jsonMsg.EmailCc,
		Content:           jsonMsg.Content,
		MessageReceivedAt: time.Now(),
	}

	mailErr := email.SendEmail(emailInfo)

	if mailErr != nil {
		return mailErr
	}

	return nil
}

func validateMessage(message domain.Message) bool {
	if !email.IsValidEmail(message.Email) {
		fmt.Println("Email is invalid:", message.Email)
		return false
	}

	if message.EmailCc != "" && !email.IsValidEmail(message.Email) {
		fmt.Println("EmailCc is invalid:", message.EmailCc)
		return false
	}

	if len(message.Content) == 0 {
		fmt.Println("Content is empty.")
		return false
	}

	return true
}
