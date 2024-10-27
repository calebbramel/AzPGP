package pgp

import (
	"fmt"
	"os"

	"github.com/ProtonMail/gopenpgp/v3/constants"
	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/ProtonMail/gopenpgp/v3/profile"
	"github.com/calebbramel/azpgp/internal/azenv"
)

type Recipient struct {
	ID          string `json:"id"`
	Fingerprint string `json:"fingerprint"`
}

type Recipients struct {
	Recipients []Recipient `json:"recipients"`
}

type Key struct {
	Value       string
	Fingerprint string
	Name        string
	ID          string
}

var PGPHandler *crypto.PGPHandle
var RecipientsList Recipients

func init() {
	azenv.Load()

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
