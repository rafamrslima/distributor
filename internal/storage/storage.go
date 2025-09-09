package storage

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/joho/godotenv"
)

const CONTAINER = "reports"

func UploadFile(file bytes.Buffer, fileName string, directory string) error {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system env.")
	}

	connString := os.Getenv("BLOB_STORAGE_CONNECTION_STRING")
	svc, err := azblob.NewClientFromConnectionString(connString, nil)

	if err != nil {
		log.Println("Error to connect to storage.")
		return err
	}

	blobName := fmt.Sprintf("%s/%s", directory, fileName)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	headers := blob.HTTPHeaders{
		BlobContentType: to.Ptr("application/pdf"),
	}

	_, err = svc.UploadBuffer(ctx, CONTAINER, blobName, file.Bytes(),
		&azblob.UploadBufferOptions{
			HTTPHeaders: &headers,
		})

	if err != nil {
		log.Println("Error to upload file to storage.")
		return err
	}

	log.Println("Upload completed:", blobName)
	return nil
}
