package blobhandler

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/calebbramel/azpgp/internal/logger"
	"github.com/calebbramel/azpgp/internal/pgp"
)

func UpdateRecipients(client *azblob.Client) error {
	var recipientsWrapper pgp.Recipients

	data, err := Get(client, os.Getenv("STORAGE_CONTAINER_NAME"), "recipient.json")
	if err != nil {
		return fmt.Errorf("failed to get blob: %v", err)
	}
	logger.Debugf("Raw data: %s", string(data))

	err = json.Unmarshal(data, &recipientsWrapper)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	// Log the unmarshalled recipients
	recipientsJSON, err := json.Marshal(recipientsWrapper.Recipients)
	if err != nil {
		return fmt.Errorf("failed to marshal recipients to JSON: %v", err)
	}
	logger.Debugf("Unmarshalled recipients: %s", string(recipientsJSON))

	pgp.RecipientsList = recipientsWrapper
	return nil
}
