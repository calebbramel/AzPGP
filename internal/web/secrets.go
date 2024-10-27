package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/calebbramel/azpgp/internal/keyvault"
)

func SecretsHandler(w http.ResponseWriter, r *http.Request) {
	vaultName := r.URL.Path[len("/secrets/"):]

	client, err := keyvault.AuthenticateSecrets(azCredential, vaultName)
	if err != nil {
		log.Fatalf("Failed to authenticate to Key Vault: %v", err)
	}

	allSecrets, err := keyvault.GetAllSecrets(client) // Update the package-level secrets variable
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(keyvault.Secrets) // Respond with the JSON
	case http.MethodPost:
		keyvault.Secrets = allSecrets
		w.WriteHeader(http.StatusOK) // Respond with 200 OK
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
