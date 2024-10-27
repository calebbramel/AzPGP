package pgp

import (
	"bytes"
	"fmt"
	"io"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/calebbramel/azpgp/internal/logger"
)

func Encrypt(PGPHandler *crypto.PGPHandle, publicKeyStr string, privateKeyStr string, sourceFile []byte) ([]byte, error) {
	logger.Debugf("PrivateKey %s\n", privateKeyStr)

	// Encrypt data with a public key and sign with private key streaming
	publicKey, err := crypto.NewKeyFromArmored(publicKeyStr)
	if err != nil {
		return nil, fmt.Errorf("error loading public key: %v: %w", publicKeyStr, err)
	}
	privateKey, err := crypto.NewKeyFromArmored(privateKeyStr)
	if err != nil {
		return nil, fmt.Errorf("error loading private key: %w", err)
	}
	defer privateKey.ClearPrivateParams()

	encHandle, err := PGPHandler.Encryption().
		Recipient(publicKey).
		SigningKey(privateKey).
		New()
	if err != nil {
		return nil, fmt.Errorf("error creating encryption handle: %w", err)
	}

	fileReader := bytes.NewReader(sourceFile)
	var ciphertextWriter bytes.Buffer
	ptWriter, err := encHandle.EncryptingWriter(&ciphertextWriter, crypto.Armor)
	if err != nil {
		return nil, fmt.Errorf("error creating encrypting writer: %w", err)
	}

	if _, err = io.Copy(ptWriter, fileReader); err != nil {
		return nil, fmt.Errorf("error writing to encrypting writer: %w", err)
	}

	if err = ptWriter.Close(); err != nil {
		return nil, fmt.Errorf("error closing encrypting writer: %w", err)
	}

	return ciphertextWriter.Bytes(), nil
}

func Decrypt(PGPHandler *crypto.PGPHandle, publicKeyStr string, privateKeyStr string, encryptedFile []byte) ([]byte, error) {
	// Decrypt armored encrypted message using the private key and obtain the plaintext with streaming
	publicKey, err := crypto.NewKeyFromArmored(publicKeyStr)
	if err != nil {
		return nil, fmt.Errorf("error loading public key: %w", err)
	}
	privateKey, err := crypto.NewKeyFromArmored(privateKeyStr)
	if err != nil {
		return nil, fmt.Errorf("error loading private key: %w", err)
	}
	defer privateKey.ClearPrivateParams()

	decHandle, err := PGPHandler.
		Decryption().
		DecryptionKey(privateKey).
		VerificationKey(publicKey).
		New()
	if err != nil {
		return nil, fmt.Errorf("error creating decryption handle: %w", err)
	}

	ciphertextReader := bytes.NewReader(encryptedFile)
	decryptedReader, err := decHandle.DecryptingReader(ciphertextReader, crypto.Armor)
	if err != nil {
		return nil, fmt.Errorf("error creating decrypting reader: %w", err)
	}

	var decrypted bytes.Buffer
	if _, err = io.Copy(&decrypted, decryptedReader); err != nil {
		return nil, fmt.Errorf("error reading decrypted data: %w", err)
	}

	return decrypted.Bytes(), nil
}

func GetRecipient(publicKeyStr string) (Key, error) {
	key, err := crypto.NewKeyFromArmored(publicKeyStr)
	if err != nil {
		return Key{}, fmt.Errorf("error reading decrypted data: %w", err)
	}

	// Extract recipient information
	entity := key.GetEntity()
	var name, id string
	for _, identity := range entity.Identities {
		name = identity.Name
		id = identity.UserId.Email
		break
	}

	return Key{
		Value: publicKeyStr,
		Name:  name,
		ID:    id,
	}, nil
}
