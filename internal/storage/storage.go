package storage

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/joho/godotenv"
	"github.com/rafamrslima/distributor/internal/domain"
)

const CONTAINER = "reports"

func UploadFile(file *os.File, message domain.Message) error {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system env.")
	}

	connString := os.Getenv("BLOB_STORAGE_CONNECTION_STRING")

	svc, err := azblob.NewClientFromConnectionString(connString, nil)

	if err != nil {
		log.Println("Error to connect to storage.")
		return err
	}

	defer file.Close()

	blobName := fmt.Sprintf("%s/%s", message.ClientName, file.Name())

	_, err = svc.UploadFile(context.Background(), CONTAINER, blobName, file, nil)
	if err != nil {
		log.Println("Error to upload file to storage.")
		return err
	}

	log.Println("Upload completed:", blobName)
	return nil
}
