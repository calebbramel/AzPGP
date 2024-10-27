$keyVault = Read-Host("Input Keyvault Name")

Write-Output("Testing Secret Create")
$headers = New-Object "System.Collections.Generic.Dictionary[[String],[String]]"
$headers.Add("Content-Type", "text/plain")

$body = @"
{
    `"Recipient`": `"Test User <test.user@example.com>`",
    `"ID`": `"test.user@example.com`"
}
"@


$response = Invoke-RestMethod 'http://localhost:8080/keys' -Method 'POST' -Headers $headers -Body $body
$response | ConvertTo-Json

Start-Sleep 1

Write-Output("Testing Secrets Pull")
$response = Invoke-RestMethod 'http://localhost:8080/secrets/KV-INFRA-SB-CUS-001' -Method 'GET' -Headers $headers
$response.attributes
$response | ConvertTo-Json

Start-Sleep 1

Write-Output("Testing Secrets Refresh")
$response = Invoke-RestMethod 'http://localhost:8080/secrets/KV-INFRA-SB-CUS-001' -Method 'POST' -Headers $headers
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

$response = Invoke-RestMethod 'http://localhost:8080/keys/<keyVaultName>' -Method 'POST' -Headers $headers -Body $body
$response | ConvertTo-Json

Start-Sleep 1

Write-Output("Testing Private Key Get")
$response = Invoke-RestMethod "http://localhost:8080/keys/private/$($response.fingerprint)" -Method 'GET' -Headers $headers
$response | ConvertTo-Json

Start-Sleep 1

Write-Output("Testing Public Key Get")
$response = Invoke-RestMethod 'http://localhost:8080/keys/public?id=test.user%40example.com' -Method 'GET' -Headers $headers
$response | ConvertTo-Json

Start-Sleep 1

Write-Output("Testing File Encrypt")
$headers = New-Object "System.Collections.Generic.Dictionary[[String],[String]]"
$headers.Add("Content-Type", "text/plain")

$body = "test content"

$response = Invoke-RestMethod 'http://localhost:8080/files/encrypt?filename=file.txt&recipient=test.user%40example.com' -Method 'POST' -Headers $headers -Body $body
$response | ConvertTo-Json
