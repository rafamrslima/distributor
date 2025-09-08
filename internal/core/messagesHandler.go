package core

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/rafamrslima/distributor/internal/db"
	"github.com/rafamrslima/distributor/internal/domain"
	"github.com/rafamrslima/distributor/internal/email"
	"github.com/rafamrslima/distributor/internal/storage"
)

func Handle(message *azservicebus.ReceivedMessage) error {
	var parsedMessage domain.Message
	err := json.Unmarshal(message.Body, &parsedMessage)

	if err != nil {
		log.Println("Error to parse received json. MessageId:", message.MessageID)
		return err
	}

	parsedMessage.MessageReceivedAt = time.Now()
	parsedMessage.Content = parsedMessage.ReportName + "dummy content" // todo
	err = db.SaveReceivedMessages(parsedMessage)

	if err != nil {
		log.Println("Error to save data in the database", err.Error())
		return err
	}

	isValid := validateMessage(parsedMessage)

	if !isValid {
		return fmt.Errorf("message is invalid")
	}

	fileName := fmt.Sprintf("%s-%s", parsedMessage.ReportName, time.Now().Format("20060102_1504"))
	file, err := os.Create(fileName)

	if err != nil {
		return err
	}

	err = storage.UploadFile(file, parsedMessage)

	if err != nil {
		return err
	}

	return nil
}

func validateMessage(message domain.Message) bool {
	if !email.IsValidEmail(message.Email) {
		log.Println("Email is invalid:", message.Email)
		return false
	}

	if len(message.Content) == 0 {
		log.Println("Content is empty.")
		return false
	}

	return true
}
