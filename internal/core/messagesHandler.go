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
		fmt.Println("Error to parse received json.")
		return err
	}

	db.SaveReceivedMessages(jsonMsg)

	if !email.IsValidEmail(jsonMsg.Email) {
		fmt.Println("Email is invalid.")
		return fmt.Errorf("email is invalid")
	}

	if jsonMsg.EmailCc != "" && !email.IsValidEmail(jsonMsg.Email) {
		fmt.Println("EmailCc is invalid.")
		return fmt.Errorf("emailCc is invalid")
	}

	if len(jsonMsg.Content) == 0 {
		fmt.Println("Content is empty.")
		return fmt.Errorf("content is empty")
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
