package web

import (
	"encoding/json"
	"net/http"

	"github.com/calebbramel/azpgp/internal/azenv"
	"github.com/calebbramel/azpgp/internal/keyvault"
	"github.com/calebbramel/azpgp/internal/logger"
)

/*	case http.MethodGet:
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(keyvault.Secrets) // Respond with the JSON
*/

func SecretsHandler(w http.ResponseWriter, r *http.Request) {
	vaultName := r.URL.Path[len("/secrets/"):]

	client, err := keyvault.AuthenticateSecrets(azenv.AzCredential, vaultName)
	logger.HandleErrf("Failed to authenticate to Key Vault: %v", err)

	allSecrets, err := keyvault.GetAllSecrets(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.HandleErrf("Failed to collect Key Vault secrets: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		if err := json.NewEncoder(w).Encode(keyvault.Secrets); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case http.MethodPost:
		keyvault.Secrets = allSecrets
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
