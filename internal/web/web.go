package web

import (
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/joho/godotenv"
)

type RequestBody struct {
	Recipient string `json:"username"`
	ID        string `json:"id"`
}

type Response struct {
	Message string `json:"message"`
}

var debugFlag bool
var azCredential *azidentity.ClientSecretCredential

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
	debugEnv := os.Getenv("DEBUG")
	debugFlag = debugEnv == "true"

}
