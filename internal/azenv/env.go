package azenv

import (
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/joho/godotenv"
)

var DebugFlag bool
var AzCredential *azidentity.ClientSecretCredential

func Load() {
	// Local Testing
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	var err error
	// Loading Service Principal Connection
	AzCredential, err = azidentity.NewClientSecretCredential(
		os.Getenv("AZURE_TENANT_ID"),
		os.Getenv("AZURE_CLIENT_ID"),
		os.Getenv("AZURE_CLIENT_SECRET"),
		nil,
	)
	if err != nil {
		log.Fatalf("failed to create credential: %v", err)
	}
	// Debug
	debugEnv := os.Getenv("DEBUG")
	DebugFlag = debugEnv == "true"
}
