package pgp

import (
	"encoding/json"
	"fmt"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/calebbramel/azpgp/internal/logger"
)

func GetFingerprintFromEncryptedFile(encryptedFile []byte) (string, error) {
	logger.Debugf("file bytes %x", encryptedFile[:50])

	// Leverage gopenpgp to get key IDs
	pgpMessage, err := crypto.NewPGPMessageFromArmored(string(encryptedFile))
	if err != nil {
		return "", fmt.Errorf("error creating PGP message: %w", err)
	}

	keyIDs, _ := pgpMessage.EncryptionKeyIDs()
	if len(keyIDs) == 0 {
		return "", fmt.Errorf("no key IDs found in the PGP message")
	}

	// Convert the key ID to a fingerprint string
	fingerprint := fmt.Sprintf("%X", keyIDs[0])

	return fingerprint, nil
}

func UpdateJSON(data []byte, newRecipient Recipient) ([]byte, error) {
	var recipients Recipients
	err := json.Unmarshal(data, &recipients)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	// Append or update the recipient
	found := false
	for i, recipient := range recipients.Recipients {
		if recipient.ID == newRecipient.ID {
			recipients.Recipients[i] = newRecipient
			found = true
			break
		}
	}
	if !found {
		recipients.Recipients = append(recipients.Recipients, newRecipient)
	}

	return json.Marshal(recipients)
}

func FindFingerprintByID(recipients Recipients, id string) (string, error) {
	logger.Debugf("Looking for fingerprint associated with %s\n", id)
	for _, recipient := range recipients.Recipients {
		if recipient.ID == id {
			logger.Debugf("Found fingerprint: %s\n", recipient.Fingerprint)
			jsonData, err := json.Marshal(recipients)
			if err != nil {
				return "", fmt.Errorf("error marshalling recipients: %v", err)
			}
			logger.Debugln(string(jsonData))
			return recipient.Fingerprint, nil
		}
	}
	jsonData, err := json.Marshal(recipients)
	if err != nil {
		return "", fmt.Errorf("error marshalling recipients: %v", err)
	}
	logger.Debugln(string(jsonData))
	return "", fmt.Errorf("recipient not found")
}
