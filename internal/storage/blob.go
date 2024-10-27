package storageBlob

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/calebbramel/azpgp/internal/debug"
	"github.com/calebbramel/azpgp/internal/pgp"
	"github.com/joho/godotenv"
)

// Azure Storage Quickstart Sample - Demonstrate how to upload, list, download, and delete blobs.
//
// Documentation References:
// - What is a Storage Account - https://docs.microsoft.com/azure/storage/common/storage-create-storage-account
// - Blob Service Concepts - https://docs.microsoft.com/rest/api/storageservices/Blob-Service-Concepts
// - Blob Service Go SDK API - https://godoc.org/github.com/Azure/azure-storage-blob-go
// - Blob Service REST API - https://docs.microsoft.com/rest/api/storageservices/Blob-Service-REST-API
// - Scalability and performance targets - https://docs.microsoft.com/azure/storage/common/storage-scalability-targets
// - Azure Storage Performance and Scalability checklist https://docs.microsoft.com/azure/storage/common/storage-performance-checklist
// - Storage Emulator - https://docs.microsoft.com/azure/storage/common/storage-use-emulator
var debugFlag bool
var azCredential *azidentity.ClientSecretCredential

func init() {
	// Local Testing
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	var err error
	azCredential, err = azidentity.NewClientSecretCredential(
		os.Getenv("AZURE_TENANT_ID"),
		os.Getenv("AZURE_CLIENT_ID"),
		os.Getenv("AZURE_CLIENT_SECRET"),
		nil,
	)
	if err != nil {
		log.Fatalf("failed to create credential: %v", err)
	}
	debugEnv := os.Getenv("DEBUG")
	debugFlag = debugEnv == "true"

}

func AuthenticateAccount(credential *azidentity.ClientSecretCredential, accountName string) (*azblob.Client, error) {
	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	client, err := azblob.NewClient(blobURL, credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate to storage account: %v", err)
	}

	return client, nil
}

func Create(client *azblob.Client, accountName string, containerName string, data []byte, blobName string) (string, error) {
	ctx := context.Background()

	_, err := client.UploadBuffer(ctx, containerName, blobName, data, &azblob.UploadBufferOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to create blob: %v", err)
	}

	// Construct the blob URL
	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", accountName, containerName, blobName)
	fmt.Printf("Uploaded blob %s\n", blobName)
	return blobURL, nil
}

func Get(client *azblob.Client, containerName string, blobName string) ([]byte, error) {
	blobDownloadResponse, err := client.DownloadStream(context.TODO(), containerName, blobName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to download data: %v", err)
	}

	reader := blobDownloadResponse.Body

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read downloaded data: %v", err)
	}

	return data, nil
}

func UpdateRecipients(client *azblob.Client) error {
	var recipientsWrapper pgp.Recipients

	data, err := Get(client, os.Getenv("STORAGE_CONTAINER_NAME"), "recipient.json")
	if err != nil {
		return fmt.Errorf("failed to get blob: %v", err)
	}
	debug.Logf(debugFlag, "Raw data: %s", string(data))

	err = json.Unmarshal(data, &recipientsWrapper)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	// Log the unmarshalled recipients
	recipientsJSON, err := json.Marshal(recipientsWrapper.Recipients)
	if err != nil {
		return fmt.Errorf("failed to marshal recipients to JSON: %v", err)
	}
	debug.Logf(debugFlag, "Unmarshalled recipients: %s", string(recipientsJSON))

	pgp.RecipientsList = recipientsWrapper
	return nil
}
