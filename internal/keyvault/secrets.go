package keyvault

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
)

var Secrets []azsecrets.SecretProperties //[]Secret

var secretValueCache = make(map[string]string)

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
