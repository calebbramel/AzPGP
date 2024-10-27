package blobhandler

import (
	"context"
	"fmt"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/calebbramel/azpgp/internal/azenv"
)

// - Storage Emulator - https://docs.microsoft.com/azure/storage/common/storage-use-emulator

func init() {
	azenv.Load()
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
