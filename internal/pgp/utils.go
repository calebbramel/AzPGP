package pgp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/ProtonMail/gopenpgp/v3/constants"
	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/ProtonMail/gopenpgp/v3/profile"
	"github.com/calebbramel/azpgp/internal/debug"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

var debugFlag bool
var PGPHandler *crypto.PGPHandle

type Recipient struct {
	ID          string `json:"id"`
	Fingerprint string `json:"fingerprint"`
}

type Recipients struct {
	Recipients []Recipient `json:"recipients"`
}

var RecipientsList Recipients

func init() {
	// Local Testing
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	debugEnv := os.Getenv("DEBUG")
	debugFlag = debugEnv == "true"

	RFC9850 := os.Getenv("RFC9850_ENABLED")
	if RFC9850 == "true" {
		PGPHandler = crypto.PGPWithProfile(profile.RFC9580())
	} else {
		PGPHandler = crypto.PGP() // for default
	}
}

func GenerateKey(username string, email string) (string, string, string, error) {
	keyHandler := PGPHandler.KeyGeneration().
		AddUserId(username, email).
		New()
	key, err := keyHandler.GenerateKeyWithSecurity(constants.HighSecurity)
	if err != nil {
		return "", "", "", fmt.Errorf("error generating key: %w", err)
	}

	armoredKey, err := key.Armor()
	if err != nil {
		return "", "", "", fmt.Errorf("error armoring key: %w", err)
	}

	publicKey, err := key.GetArmoredPublicKey()
	if err != nil {
		return "", "", "", fmt.Errorf("error getting public key: %w", err)
	}

	fingerprint := key.GetFingerprint()

	return armoredKey, publicKey, fingerprint, nil
}

func GetFingerprintFromEncryptedFile(encryptedFile []byte) (string, error) {
	// Create a reader for the encrypted file
	ciphertextReader := bytes.NewReader(encryptedFile)

	// Decode the armored data
	block, err := armor.Decode(ciphertextReader)
	if err != nil {
		return "", fmt.Errorf("error decoding armored data: %w", err)
	}

	// Create a packet reader
	packetReader := packet.NewReader(block.Body)

	// Read the encrypted message
	ciphertext, err := packetReader.Next()
	if err != nil {
		return "", fmt.Errorf("error reading PGP packet: %w", err)
	}

	// Check if the packet is an encrypted key packet
	encryptedKeyPacket, ok := ciphertext.(*packet.EncryptedKey)
	if !ok {
		return "", fmt.Errorf("no encrypted key packet found")
	}

	// Get the key ID from the encrypted key packet
	keyID := encryptedKeyPacket.KeyId

	// Convert the key ID to a fingerprint string
	fingerprint := fmt.Sprintf("%X", keyID)

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
	debug.Logf(debugFlag, "Looking for fingerprint associated with %s\n", id)
	for _, recipient := range recipients.Recipients {
		if recipient.ID == id {
			debug.Logf(debugFlag, "Found fingerprint: %s\n", recipient.Fingerprint)
			jsonData, err := json.Marshal(recipients)
			if err != nil {
				return "", fmt.Errorf("error marshalling recipients: %v", err)
			}
			debug.Logln(debugFlag, string(jsonData))
			return recipient.Fingerprint, nil
		}
	}
	jsonData, err := json.Marshal(recipients)
	if err != nil {
		return "", fmt.Errorf("error marshalling recipients: %v", err)
	}
	debug.Logln(debugFlag, string(jsonData))
	return "", fmt.Errorf("recipient not found")
}
