package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/calebbramel/azpgp/internal/keyvault"
	storageBlob "github.com/calebbramel/azpgp/internal/storage"
	"github.com/calebbramel/azpgp/internal/web"
	"github.com/joho/godotenv"
)

var azCredential *azidentity.ClientSecretCredential
var secrets []azsecrets.SecretProperties

type Recipient struct {
	ID          string `json:"id"`
	Fingerprint string `json:"fingerprint"`
}

type Recipients struct {
	Recipients []Recipient `json:"recipients"`
}

func init() {
	// Local Testing
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	var err error
	azCredential, err = azidentity.NewClientSecretCredential(
		os.Getenv("AZURE_TENANT_ID"),
		os.Getenv("AZURE_CLIENT_ID"),
		os.Getenv("AZURE_CLIENT_SECRET"),
		nil,
	)
	if err != nil {
		log.Fatalf("failed to create credential: %v", err)
	}

	client, err := keyvault.AuthenticateSecrets(azCredential, os.Getenv("KEY_VAULT_NAME"))
	if err != nil {
		log.Fatalf("Failed to authenticate to Key Vault: %v", err)
	}

	updatedSecrets, err := keyvault.GetAllSecrets(client)
	if err != nil {
		log.Fatalf("could not get secrets: %s\n", err)
	}

	keyvault.Secrets = updatedSecrets // Update the global secrets variable
	blobClient, err := storageBlob.AuthenticateAccount(azCredential, os.Getenv("STORAGE_ACCOUNT_NAME"))
	if err != nil {
		log.Fatalf("Failed to authenticate to Storage Account: %s\n", err)
	}

	storageBlob.UpdateRecipients(blobClient)
}

func main() {
	//	go keyvault.HeartBeat(azCredential) // sync secrets every 15 seconds

	mux := http.NewServeMux()
	port := "8080" //os.Getenv("LISTENING_PORT")

	// Get Keys
	mux.HandleFunc("/secrets/", web.SecretsHandler)
	// Add Key
	mux.HandleFunc("/keys", web.KeysHandler)
	mux.HandleFunc("/keys/", web.KeysHandler)
	// Stores to blob
	mux.HandleFunc("/files/encrypt", web.EncryptHandler)
	mux.HandleFunc("/files/decrypt", web.DecryptHandler)

	log.Printf("Starting server on %v", ":"+port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("could not start server: %s\n", err)
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
