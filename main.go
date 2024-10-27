package main

import (
	"log"
	"net/http"
	"os"

	"github.com/calebbramel/azpgp/internal/azenv"
	"github.com/calebbramel/azpgp/internal/blobhandler"
	"github.com/calebbramel/azpgp/internal/keyvault"
	"github.com/calebbramel/azpgp/internal/web"
)

type Recipient struct {
	ID          string `json:"id"`
	Fingerprint string `json:"fingerprint"`
}

type Recipients struct {
	Recipients []Recipient `json:"recipients"`
}

//var secrets []azsecrets.SecretProperties

func init() {
	// Local Testing
	// loads values from .env into the system
	azenv.Load()

	client, err := keyvault.AuthenticateSecrets(azenv.AzCredential, os.Getenv("KEY_VAULT_NAME"))
	if err != nil {
		log.Fatalf("Failed to authenticate to Key Vault: %v", err)
	}

	updatedSecrets, err := keyvault.GetAllSecrets(client)
	if err != nil {
		log.Fatalf("could not get secrets: %s\n", err)
	}

	keyvault.Secrets = updatedSecrets // Update the global secrets variable
	blobClient, err := blobhandler.AuthenticateAccount(azenv.AzCredential, os.Getenv("STORAGE_ACCOUNT_NAME"))
	if err != nil {
		log.Fatalf("Failed to authenticate to Storage Account: %s\n", err)
	}

	blobhandler.UpdateRecipients(blobClient)
}

func main() {
	mux := http.NewServeMux()
	port := os.Getenv("LISTENING_PORT")
	if port == "" {
		port = "8080"
	}

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
