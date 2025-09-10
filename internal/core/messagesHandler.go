package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/jung-kurt/gofpdf"
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
	err = db.SaveReceivedMessages(parsedMessage)

	if err != nil {
		log.Println("Error to save data in the database", err.Error())
		return err
	}

	isValid := validateMessage(parsedMessage)

	if !isValid {
		return fmt.Errorf("message is invalid")
	}

	reportInfo, err := db.GetReportInfo(parsedMessage.Email, parsedMessage.ReportName)

	if err != nil {
		log.Println("Error to get report info from database.", err.Error())
		return err
	}

	var fileContent string

	if len(reportInfo) > 0 {
		fileContent = fmt.Sprintf("Balance for the day is $%v", (reportInfo[0].Gains - reportInfo[0].Losses))
	} else {
		fileContent = "Balance for the day does not exist in the database"
	}
	pdfFile, err := createPdfReport(fileContent)

	if err != nil {
		log.Println("Error to create pdf report.", err.Error())
		return err
	}

	fileName := fmt.Sprintf("%s-%s.pdf", parsedMessage.ReportName, time.Now().Format("20060102_1504"))

	err = storage.UploadFile(pdfFile, fileName, parsedMessage.ClientName)

	if err != nil {
		return err
	}
	return nil
}

func validateMessage(message domain.Message) bool {
	if !email.IsValid(message.Email) {
		log.Println("Email is invalid:", message.Email)
		return false
	}

	if message.ClientName == "" {
		log.Println("Client name is invalid:", message.ClientName)
		return false
	}

	if message.ReportName == "" {
		log.Println("Report name is invalid:", message.ReportName)
		return false
	}
	return true
}

func createPdfReport(content string) (bytes.Buffer, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, content)

	var buf bytes.Buffer
	err := pdf.Output(&buf) // write PDF into buffer
	if err != nil {
		return bytes.Buffer{}, err
	}

	return buf, nil
}
