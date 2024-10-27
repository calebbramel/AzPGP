$keyVault = Read-Host("Input Keyvault Name")
$URI = Read-Host("Input server URI, I.E 'http:/localhost:8080'") 

Write-Output("Testing Secret Create")
$headers = New-Object "System.Collections.Generic.Dictionary[[String],[String]]"
$headers.Add("Content-Type", "text/plain")

$body = @"
{
    `"Recipient`": `"Test User <test.user@example.com>`",
    `"ID`": `"test.user@example.com`"
}
"@


$response = Invoke-RestMethod "$URI/keys" -Method 'POST' -Headers $headers -Body $body
$response | ConvertTo-Json

Start-Sleep 1

Write-Output("Testing Secrets Pull")
$response = Invoke-RestMethod "$URI/secrets/$keyVault" -Method 'GET' -Headers $headers
$response.attributes
$response | ConvertTo-Json

Start-Sleep 1

Write-Output("Testing Secrets Refresh")
$response = Invoke-RestMethod "$URI/secrets/$keyVault" -Method 'POST' -Headers $headers
$response | ConvertTo-Json

Start-Sleep 1

Write-Output("Testing Key Generate")
$headers = New-Object "System.Collections.Generic.Dictionary[[String],[String]]"
$headers.Add("Content-Type", "text/plain")

$body = @"
{
    `"Recipient`": `"Test User <test.user@example.com>`",
    `"ID`": `"test.user@example.com`"
}
"@

$response = Invoke-RestMethod "$URI/keys/$keyVault" -Method 'POST' -Headers $headers -Body $body
$response | ConvertTo-Json


Start-Sleep 1

Write-Output("Testing Private Key Get")
$response = Invoke-RestMethod "$URI/keys/private/$($response.fingerprint)" -Method 'GET' -Headers $headers
$response | ConvertTo-Json

Start-Sleep 1

Write-Output("Testing Public Key Get")
$response = Invoke-RestMethod "$URI/keys/public?id=test.user%40example.com" -Method 'GET' -Headers $headers
$response | ConvertTo-Json

Start-Sleep 1

Write-Output("Testing File Encrypt")
$headers = New-Object "System.Collections.Generic.Dictionary[[String],[String]]"
$headers.Add("Content-Type", "text/plain")

$body = "test content"

$response = Invoke-RestMethod "$URI/files/encrypt?filename=file.txt&recipient=test.user%40example.com" -Method 'POST' -Headers $headers -Body $body
$response | ConvertTo-Json


$encryptedFile = Split-Path -Path $response.URL -Leaf
Invoke-WebRequest -Uri $response.URL -OutFile $filename



Write-Output("Testing File Encrypt")
$headers = New-Object "System.Collections.Generic.Dictionary[[String],[String]]"
$headers.Add("Content-Type", "text/plain")

$body = "test content"

$response = Invoke-RestMethod "http://localhost:8080/files/decrypt?filename=$($encryptedFile)&recipient=test.user%40example.com" -Method 'POST' -Headers $headers -Body $body
$response | ConvertTo-Json


$decryptedFile = Split-Path -Path $response.URL -Leaf
Invoke-WebRequest -Uri $response.URL -OutFile $decryptedFile