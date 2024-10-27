package test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/calebbramel/azpgp/internal/keyvault"
	"github.com/calebbramel/azpgp/internal/pgp"
	storageBlob "github.com/calebbramel/azpgp/internal/storage"
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

func test_json() {
	data := Recipients{
		Recipients: []Recipient{
			{ID: "recipient1", Fingerprint: "ABCD1234EFGH5678"},
			{ID: "recipient2", Fingerprint: "IJKL9101MNOP2345"},
			{ID: "recipient3", Fingerprint: "QRST5678UVWX9101"},
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(jsonData))
}

func test_upload() {
	accountName := "stinfrasbcus001"
	client, err := storageBlob.AuthenticateAccount(azCredential, accountName)
	if err != nil {
		log.Fatalf("Failed to authenticate to storage account: %v", err)
	}

	containerName := "caleb"
	filePath := "test2.txt"
	testContent := []byte("test content 2: electric boogaloo")
	upload, err := storageBlob.Create(client, accountName, containerName, testContent, filePath)
	if err != nil {
		log.Fatalf("Failed to upload data: %v", err)
	}
	fmt.Println("Uploaded data:", upload)
}

func test_recipient() {
	publicKey := pgp.Key{
		Value: `-----BEGIN PGP PUBLIC KEY BLOCK-----

xiYEZIbSkxsHknQrXGfb+kM2iOsOvin8yE05ff5hF8KE6k+saspAZc0VdXNlciA8
dXNlckB0ZXN0LnRlc3Q+wocEExsIAD0FAmSG0pMJkEHsytogdrSJFiEEamc2vcEG
XMMaYxmDQezK2iB2tIkCGwMCHgECGQECCwcCFQgCFgADJwcCAABTnme46ymbAs0X
7tX3xWu+9O+LLdM0aAUyV6FwUNWcy47IfmTunwdqHZ2CbUGLLb+OR/9yci1aIHDJ
xXmJh3kj9wDOJgRkhtKTGX6Xe04jkL+7ikivpOB0/ZSq+fnZr2+76Mf/InbOrpxJ
wnQEGBsIACoFAmSG0pMJkEHsytogdrSJFiEEamc2vcEGXMMaYxmDQezK2iB2tIkC
GwwAAMJizYj3AFqQi70eHGzhHcmr0XwnsAfLGw0vQaiZn6HGITQw5nBGvXQPF9Vp
FpsXV9x/08dIdfZLAQVdQowgeBsxCw==
=JIkN

-----END PGP PUBLIC KEY BLOCK-----`,
	}
	recipient, err := pgp.GetRecipient(publicKey.Value)
	if err != nil {
		log.Fatalf("Failed to get secret recipient: %v", err)
	}
	fmt.Println("Secret Recipient:", recipient)

}

func test_retrieve() {

	vaultName := os.Getenv("KEY_VAULT_NAME")

	// Example usage
	client, err := keyvault.AuthenticateSecrets(azCredential, vaultName)
	if err != nil {
		log.Fatalf("Failed to authenticate to Key Vault: %v", err)
	}

	secrets, _ := keyvault.GetAllSecrets(client)
	fmt.Println("Secret Value: %v", secrets)
	fingerprint := "61CE9C95894276DA"

	secretValue, err := keyvault.GetSecret(client, fingerprint)
	if err != nil {
		log.Fatalf("Failed to get secret by fingerprint: %v", err)
	}
	fmt.Println("Secret Value:", secretValue)

}

func test_pgp() {
	privateKey := pgp.Key{
		Value: `-----BEGIN PGP PRIVATE KEY BLOCK-----

xUkEZIbSkxsHknQrXGfb+kM2iOsOvin8yE05ff5hF8KE6k+saspAZQCy/kfFUYc2
GkpOHc42BI+MsysKzk4ofjBAfqM+bb7goQ3hzRV1c2VyIDx1c2VyQHRlc3QudGVz
dD7ChwQTGwgAPQUCZIbSkwmQQezK2iB2tIkWIQRqZza9wQZcwxpjGYNB7MraIHa0
iQIbAwIeAQIZAQILBwIVCAIWAAMnBwIAAFOeZ7jrKZsCzRfu1ffFa77074st0zRo
BTJXoXBQ1ZzLjsh+ZO6fB2odnYJtQYstv45H/3JyLVogcMnFeYmHeSP3AMdJBGSG
0pMZfpd7TiOQv7uKSK+k4HT9lKr5+dmvb7vox/8ids6unEkAF1v8fCKogIrtBWVT
nVbwnovjM3LLexpXFZSgTKRcNMgPRMJ0BBgbCAAqBQJkhtKTCZBB7MraIHa0iRYh
BGpnNr3BBlzDGmMZg0HsytogdrSJAhsMAADCYs2I9wBakIu9Hhxs4R3Jq9F8J7AH
yxsNL0GomZ+hxiE0MOZwRr10DxfVaRabF1fcf9PHSHX2SwEFXUKMIHgbMQs=
=bJqd

-----END PGP PRIVATE KEY BLOCK-----`,
		Fingerprint: "61CE9C95894276DA",
	}

	publicKey := pgp.Key{
		Value: `-----BEGIN PGP PUBLIC KEY BLOCK-----

xiYEZIbSkxsHknQrXGfb+kM2iOsOvin8yE05ff5hF8KE6k+saspAZc0VdXNlciA8
dXNlckB0ZXN0LnRlc3Q+wocEExsIAD0FAmSG0pMJkEHsytogdrSJFiEEamc2vcEG
XMMaYxmDQezK2iB2tIkCGwMCHgECGQECCwcCFQgCFgADJwcCAABTnme46ymbAs0X
7tX3xWu+9O+LLdM0aAUyV6FwUNWcy47IfmTunwdqHZ2CbUGLLb+OR/9yci1aIHDJ
xXmJh3kj9wDOJgRkhtKTGX6Xe04jkL+7ikivpOB0/ZSq+fnZr2+76Mf/InbOrpxJ
wnQEGBsIACoFAmSG0pMJkEHsytogdrSJFiEEamc2vcEGXMMaYxmDQezK2iB2tIkC
GwwAAMJizYj3AFqQi70eHGzhHcmr0XwnsAfLGw0vQaiZn6HGITQw5nBGvXQPF9Vp
FpsXV9x/08dIdfZLAQVdQowgeBsxCw==
=JIkN

-----END PGP PUBLIC KEY BLOCK-----`,
	}
	// Example usage
	filePath := "test2.txt"
	testContent := []byte("test content")

	err := os.WriteFile(filePath, testContent, 0644)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	encryptedData, err := pgp.Encrypt(pgp.PGPHandler, publicKey.Value, privateKey.Value, fileContent)
	if err != nil {
		fmt.Println("Error encrypting file:", err)
		return
	}

	// Write the encrypted data to a new file
	outputFilePath := fmt.Sprintf("%v.pgp", filePath)
	err = os.WriteFile(outputFilePath, encryptedData, 0644)
	if err != nil {
		fmt.Println("Error writing encrypted file:", err)
		return
	}

	fmt.Println("File encrypted and saved as", outputFilePath)

	fingerprint, err := pgp.GetFingerprintFromEncryptedFile(encryptedData)

	if err != nil {
		fmt.Println("Error reading fingerprint:", err)
		return
	}
	fmt.Println("Private Fingerprint identified: ", fingerprint)

	vaultName := os.Getenv("KEY_VAULT_NAME")

	// Authenticate to Key Vault
	client, err := keyvault.AuthenticateSecrets(azCredential, vaultName)
	if err != nil {
		log.Fatalf("Failed to authenticate to Key Vault: %v", err)
	}

	newSecret, err := keyvault.NewPGPKeySecret(client, "key1", &privateKey)
	if err != nil {
		fmt.Println("Error creating secret:", err)
		return
	}
	fmt.Println("Secret created: ", newSecret)

	test_keyvault(fingerprint)
	// decryptedData, err := pgp.Decrypt(PGPHandler, publicKeyStr, privateKeyStr, outputFilePath)
}

func test_keyvault(fingerprint string) {
	vaultName := os.Getenv("KEY_VAULT_NAME")

	// Authenticate to Key Vault
	client, err := keyvault.AuthenticateSecrets(azCredential, vaultName)
	if err != nil {
		log.Fatalf("Failed to authenticate to Key Vault: %v", err)
	}

	if err != nil {
		log.Fatalf("Failed to retrieve secret pager: %v", err)
	}

	keyvault.Secrets, _ = keyvault.GetAllSecrets(client)

	// Retrieve the secret
	secretValue, err := keyvault.GetSecret(client, fingerprint)
	if err != nil {
		log.Fatalf("Failed to retrieve secret: %v", err)
	}
	// Print confirmation message
	fmt.Println("Secret retrieved successfully:", secretValue)
}
