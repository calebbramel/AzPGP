package web

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/calebbramel/azpgp/internal/azenv"
	"github.com/calebbramel/azpgp/internal/blobhandler"
	"github.com/calebbramel/azpgp/internal/keyvault"
	"github.com/calebbramel/azpgp/internal/logger"
	"github.com/calebbramel/azpgp/internal/pgp"
)

func DecryptHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		vaultName := os.Getenv("KEY_VAULT_NAME")
		azClient, err := keyvault.AuthenticateSecrets(azenv.AzCredential, vaultName)
		if err != nil {
			http.Error(w, "Unable to authenticate to keyvault", http.StatusBadRequest)
			return
		}

		// Retrieve the recipient name from the query parameters
		recipientName := r.URL.Query().Get("recipient")
		if recipientName == "" {
			http.Error(w, "Recipient name is required", http.StatusBadRequest)
			return
		}
		filename, err := urlDecode(r.URL.Query().Get("filename"))
		if err != nil {
			http.Error(w, "Unable to parse filename", http.StatusBadRequest)
			return
		}

		var encryptedFile []byte

		// Check the Content-Type of the request
		contentType := r.Header.Get("Content-Type")
		if contentType == "multipart/form-data" {
			// Parse the form to retrieve the file
			err = r.ParseMultipartForm(10 << 20) // 10 MB
			if err != nil {
				log.Printf("Error parsing form: %v", err)
				http.Error(w, "Unable to parse form", http.StatusBadRequest)
				return
			}

			// Retrieve the file from the form
			file, _, err := r.FormFile("file")
			if err != nil {
				http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
				return
			}
			defer file.Close()

			// Read the file into a byte slice
			encryptedFile, err = io.ReadAll(file)
			if err != nil {
				http.Error(w, "Unable to read file", http.StatusInternalServerError)
				return
			}
		} else if contentType == "text/plain" || contentType == "application/pgp-encrypted" {
			// Read the plain text or PGP encrypted body
			encryptedFile, err = io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Unable to read body", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Unsupported content type", http.StatusBadRequest)
			return
		}

		recipient, err := urlDecode(recipientName)
		if err != nil {
			http.Error(w, "Unable to decode recipient", http.StatusBadRequest)
			return
		}
		publicKeyStr, err := keyvault.GetSecret(azClient, recipient)
		if err != nil {
			http.Error(w, "Error retrieving public key", http.StatusInternalServerError)
			return
		}

		fingerprint, err := pgp.FindFingerprintByID(pgp.RecipientsList, recipient)
		if err != nil {
			http.Error(w, "Unable to locate private key fingerprint", http.StatusInternalServerError)
			return
		}
		fileFingerprint, err := pgp.GetFingerprintFromEncryptedFile(encryptedFile)
		if err != nil {
			logger.HandleErrf("error retrieving fingerprint from file: %v", err)
			http.Error(w, "Unable to locate private key fingerprint", http.StatusInternalServerError)
			return
		}
		logger.Debugf("fingerprint from file: %v", fileFingerprint)
		privateKeyStr, err := keyvault.GetSecret(azClient, fingerprint)
		if err != nil {
			http.Error(w, "Unable to locate private key", http.StatusInternalServerError)
			return
		}

		//		logger.Debugf("public key: %s", publicKeyStr)
		//		logger.Debugf("private key: %s", privateKeyStr)
		//		logger.Debugf("fingerprint: %s", fingerprint)

		// Decrypt the file
		decryptedFile, err := pgp.Decrypt(pgp.PGPHandler, publicKeyStr, privateKeyStr, encryptedFile)
		if err != nil {
			logger.Debugf("Unable to decrypt file: %v", err)
			http.Error(w, "Unable to decrypt file", http.StatusInternalServerError)
			return
		}

		// Upload the decrypted file to Azure Blob Storage
		blobClient, err := blobhandler.AuthenticateAccount(azenv.AzCredential, os.Getenv("STORAGE_ACCOUNT_NAME"))
		if err != nil {
			log.Fatalf("Failed to authenticate to Storage Account: %s\n", err)
		}
		blobName := filename[:len(filename)-len(".pgp")]
		blobURL, err := blobhandler.Create(blobClient, os.Getenv("STORAGE_ACCOUNT_NAME"), os.Getenv("STORAGE_CONTAINER_NAME"), decryptedFile, blobName)
		if err != nil {
			http.Error(w, "Unable to upload file to Azure Blob Storage", http.StatusInternalServerError)
			return
		}

		response := map[string]string{
			"URL": blobURL,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
