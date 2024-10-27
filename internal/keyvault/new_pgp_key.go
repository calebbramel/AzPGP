package keyvault

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/calebbramel/azpgp/internal/pgp"
)

func NewPGPKeySecret(client *azsecrets.Client, secretName string, PGPKey *pgp.Key) (string, error) {
	var params azsecrets.SetSecretParameters

	if PGPKey.Fingerprint != "" {
		params = azsecrets.SetSecretParameters{Value: &PGPKey.Value, ContentType: &PGPKey.Fingerprint}
	} else {
		params = azsecrets.SetSecretParameters{Value: &PGPKey.Value, ContentType: &PGPKey.ID}
	}

	resp, err := client.SetSecret(context.TODO(), secretName, params, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create a secret: %w", err)
	}
	fmt.Printf("Updated secreted %s\n", secretName)
	return fmt.Sprint(*resp.ID), nil
}
