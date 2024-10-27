package keyvault

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/calebbramel/azpgp/internal/pgp"
)

/*
	type Secret struct {
		Properties azsecrets.SecretProperties
		Value      string
	}
*/

var Secrets []azsecrets.SecretProperties //[]Secret

var secretValueCache = make(map[string]string)

func HeartBeat(credential *azidentity.ClientSecretCredential) {
	vaultName := os.Getenv("KEY_VAULT_NAME")

	client, err := AuthenticateSecrets(credential, vaultName)
	if err != nil {
		log.Fatalf("Failed to authenticate to Key Vault: %v", err)
	}

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		updatedSecrets, err := GetAllSecrets(client)
		if err != nil {
			log.Fatalf("could not get secrets: %s\n", err)
		}

		Secrets = updatedSecrets // Update the global secrets variable

		secretsJSON, _ := json.Marshal(Secrets)
		fmt.Println(string(secretsJSON)) // Print the JSON string to the console
	}
}

func AuthenticateSecrets(credential *azidentity.ClientSecretCredential, vaultName string) (*azsecrets.Client, error) {
	vaultURI := fmt.Sprintf("https://%s.vault.azure.net/", vaultName)

	client, err := azsecrets.NewClient(vaultURI, credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create key vault client: %v", err)
	}

	return client, nil
}

func GetSecret(client *azsecrets.Client, ContentType string) (string, error) {
	fmt.Printf("Getting secret %s\n", ContentType)
	if value, found := secretValueCache[ContentType]; found {
		return value, nil
	}
	fmt.Printf("Secret %s not in cache\n", ContentType)
	// If not in the cache, retrieve the secret value
	for _, secret := range Secrets {
		if secret.ContentType != nil && *secret.ContentType == ContentType {
			fmt.Printf("Found Secret ID: %s\n", *secret.ID)
			secretResp, err := client.GetSecret(context.TODO(), secret.ID.Name(), "", nil)
			if err != nil {
				return "", err
			}
			secretValue := *secretResp.Value
			// Store the secret value in the cache
			secretValueCache[ContentType] = secretValue
			return secretValue, nil
		}
	}
	return "", fmt.Errorf("secret with ContentType %s not found", ContentType)
}

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

func GetAllSecrets(client *azsecrets.Client) ([]azsecrets.SecretProperties, error) {
	pager := client.NewListSecretPropertiesPager(nil)

	var allSecrets []azsecrets.SecretProperties

	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}
		for _, secretPtr := range page.Value {
			if secretPtr != nil {
				allSecrets = append(allSecrets, *secretPtr)
			}
		}
	}
	return allSecrets, nil
}
