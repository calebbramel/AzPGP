package keyvault

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

// Unused, can update secret cache on a cadence
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
