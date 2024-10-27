package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/calebbramel/azpgp/internal/debug"
	"github.com/calebbramel/azpgp/internal/keyvault"
	"github.com/calebbramel/azpgp/internal/pgp"
	storageBlob "github.com/calebbramel/azpgp/internal/storage"
)

func KeysHandler(w http.ResponseWriter, r *http.Request) {
	debug.Logf(debugFlag, "Recieved request %s\n", r.URL.String())
	vaultName := os.Getenv("KEY_VAULT_NAME")
	azClient, err := keyvault.AuthenticateSecrets(azCredential, vaultName)
	if err != nil {
		log.Fatalf("Failed to authenticate to Key Vault: %v", err)
	}

	switch r.Method {
	case http.MethodPost:
		var body RequestBody
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&body)
		if err != nil {
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}
		privateKeyStr, publicKeyStr, fingerprint, err := pgp.GenerateKey(body.Recipient, body.ID)
		if err != nil {
			http.Error(w, "Error creating key pair", http.StatusInternalServerError)
			return
		}

		response := map[string]string{
			"privateKey":  privateKeyStr,
			"publicKey":   publicKeyStr,
			"fingerprint": fingerprint,
		}

		newData := pgp.Recipient{
			ID:          body.ID,
			Fingerprint: fingerprint,
		}

		publicKey := pgp.Key{
			Value: publicKeyStr,
			ID:    body.ID,
		}

		privateKey := pgp.Key{
			Value:       privateKeyStr,
			Fingerprint: fingerprint,
		}

		storageClient, err := storageBlob.AuthenticateAccount(azCredential, os.Getenv("STORAGE_ACCOUNT_NAME"))
		if err != nil {
			http.Error(w, "Error authenticating to blob", http.StatusInternalServerError)
			return
		}
		oldData, err := storageBlob.Get(storageClient, os.Getenv("STORAGE_CONTAINER_NAME"), "recipient.json")
		if err != nil {
			http.Error(w, "Error downloading blob data", http.StatusInternalServerError)
			return
		}
		dataBytes, err := pgp.UpdateJSON(oldData, newData)
		if err != nil {
			http.Error(w, "Error updating keyring json", http.StatusInternalServerError)
			return
		}
		secretName := sanitizeString(body.ID)
		fmt.Printf("Storing secret %s\n", secretName+"-publicKey")
		keyvault.NewPGPKeySecret(azClient, secretName+"-publicKey", &publicKey)
		fmt.Printf("Storing secret %s\n", secretName+"-privateKey")
		keyvault.NewPGPKeySecret(azClient, secretName+"-privateKey", &privateKey)
		storageBlob.Create(storageClient, os.Getenv("STORAGE_ACCOUNT_NAME"), os.Getenv("STORAGE_CONTAINER_NAME"), dataBytes, "recipient.json")

		updatedSecrets, err := keyvault.GetAllSecrets(azClient)
		if err != nil {
			log.Fatalf("could not get secrets: %s\n", err)
		}

		keyvault.Secrets = updatedSecrets // Update the global secrets variable

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case http.MethodGet:
		path := r.URL.Path
		debug.Logf(debugFlag, "Got path %s\n", path)
		if strings.HasPrefix(path, "/keys/private/") {
			fingerprint := r.URL.Path[len("/keys/private/"):]
			debug.Logf(debugFlag, "Got fingerprint %s\n", fingerprint)
			privateKey, err := keyvault.GetSecret(azClient, fingerprint)
			if err != nil {
				http.Error(w, "Error retrieving private key", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(map[string]string{"privateKey": privateKey}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else if strings.HasPrefix(path, "/keys/public") {
			recipient, err := urlDecode(r.URL.Query().Get("id"))
			if err != nil {
				http.Error(w, "Error decoding ID", http.StatusInternalServerError)
				return
			}

			if recipient == "" {
				http.Error(w, "Email query parameter is required", http.StatusBadRequest)
				return
			}

			publicKey, err := keyvault.GetSecret(azClient, recipient)
			if err != nil {
				http.Error(w, "Error retrieving public key", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(map[string]string{"publicKey": publicKey}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
