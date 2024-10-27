# FreshService eBonding Middleware

This project is a PGP tool leveraging the Azure backend for processing and storage.

## Table of Contents

- Installation
- Usage
- API Endpoints

## Requirements
1. Azure Functions Core Tools

### Windows
```PowerShell
curl -L -o azure-functions-core-tools.msi https://go.microsoft.com/fwlink/?linkid=2174087
Start-Process msiexec.exe -ArgumentList '/i azure-functions-core-tools.msi /quiet /norestart' -NoNewWindow -Wait
```

### MacOS
```bash
brew tap azure/functions
brew install azure-functions-core-tools@4
# if upgrading on a machine that has 2.x or 3.x installed:
brew link --overwrite azure-functions-core-tools@4
```

### Linux
#### Add Microsoft Key
```bash
curl https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > microsoft.gpg
sudo mv microsoft.gpg /etc/apt/trusted.gpg.d/microsoft.gpg
```

#### Ubuntu
```bash
sudo sh -c 'echo "deb [arch=amd64] https://packages.microsoft.com/repos/microsoft-ubuntu-$(lsb_release -cs)-prod $(lsb_release -cs) main" > /etc/apt/sources.list.d/dotnetdev.list'
```
#### Debian
```bash
sudo sh -c 'echo "deb [arch=amd64] https://packages.microsoft.com/debian/$(lsb_release -rs | cut -d'.' -f 1)/prod $(lsb_release -cs) main" > /etc/apt/sources.list.d/dotnetdev.list'
```
2. Go

Navigate to https://go.dev/doc/install for download instructions

## Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/calebbramel/azpgp.git
    ```
2. Install the dependencies:
    ```bash
	cd azpgp
    func init azpgp
    ```
## Usage

To start the application, run:
```bash
func start
```

### Endpoints and Expected Fields
#### Key Ring
##### List Secrets 
`GET /secrets/<keyVaultName>`

Response: `200 OK`
##### Refresh Secrets
`POST /secrets/<keyVaultName>`
Response: `200 OK`
#### Key Management
##### Create New Private Key
POST /keys/<keyVaultName>

Body:
```JSON
{
    "Recipient": "Friendly Name <friendly.name@example.com>",
    "ID": "friendly.name@example.com"
}
```
Response: `201 Created`
```JSON
{
    "fingerprint": "<fingerprint>",
    "privateKey": "<privateKey>",
    "publicKey": "<publicKey>"
}
```
##### Get Private Key
GET /keys/private/<fingerprint>

Response: `200 OK`
```JSON
{
"privateKey": "<privateKey>"
}
```

##### Get Public Key
GET /keys/public?id=<urlEncodedRecipient>
GET /keys/public?id=friendly.name%40example.com

Response: `200 OK`
```JSON
{
"publicKey": "publicKey"
}
```

#### File Operations
##### Encrypt File
POST /files/encrypt?recipient=<urlEncodedRecipient>&filename=<filename>
POST /files/encrypt?recipient=friendly.name%40example.com&filename=file.txt

Body: binary

Response: 201 Created
```JSON
{
    "URL": "https://<STORAGE_ACCOUNT_NAME>.blob.core.windows.net/<STORAGE_CONTAINER_NAME>/<filename>.gpg"
}
```

##### Decrypt File
POST /files/decrypt?recipient=<urlEncodedRecipient>&filename=<filename>
POST /files/decrypt?recipient=friendly.name%40example.com&filename=file.txt

Body: binary

Response: 201 Created
```JSON
{
    "URL": "https://<STORAGE_ACCOUNT_NAME>.blob.core.windows.net/<STORAGE_CONTAINER_NAME>/<filename-.gpg>"
}
```

##### Create Ticket
```
#### Add Attachment
This requires a query to be added to the endpoint:
```
/api/tickets?record=incident&number=<u_inc_number>&filename=<filename.ext>
```
The file should be attached as a binary stream in the request body. The file will be uploaded as the name in the query.
