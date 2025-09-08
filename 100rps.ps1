# скрипт тормозит в начале, надо ждать

$walletIds = @(
    "a1b2c3e4-5678-9012-3456-789012345678",
    "b2c3d4e5-6789-0123-4567-890123456789",
    "c3d4e5f6-7890-1234-5678-901234567890",
    "d4e5f6a7-8901-2345-6789-012345678901",
    "e5f6a7b8-9012-3456-7890-123456789012"
)

$operations = @("DEPOSIT", "WITHDRAW")

$maxParallel = 100
$totalRequests = 100
$completed = 0

Write-Host "Starting load test: $totalRequests requests, $maxParallel parallel"

for ($i = 1; $i -le $totalRequests; $i += $maxParallel) {
    $jobs = @()
    $batchSize = [Math]::Min($maxParallel, $totalRequests - $i + 1)
    
    Write-Host "Starting batch: requests $i - $($i + $batchSize - 1)"
    
    # Запускаем пачку
    for ($j = 0; $j -lt $batchSize; $j++) {
        $requestNum = $i + $j
        $scriptBlock = {
            param($reqNum, $walletIds, $operations)
            
            $walletId = $walletIds | Get-Random
            $operation = $operations | Get-Random
            $amount = Get-Random -Minimum 10 -Maximum 1000
            
            $body = @{
                valletId = $walletId
                operationType = $operation
                amount = $amount
            } | ConvertTo-Json

            try {
                $response = Invoke-RestMethod -Uri "http://localhost:3000/api/v1/wallet" -Method Post -Headers @{ "Content-Type" = "application/json" } -Body $body
                return ("Request {0}: {1} {2} Success - {3}" -f $reqNum, $operation, $amount, $response.reference_id)
            } catch {
                return ("Request {0}: {1} {2} Error" -f $reqNum, $operation, $amount)
            }
        }
        
        $job = Start-Job -ScriptBlock $scriptBlock -ArgumentList $requestNum, $walletIds, $operations
        $jobs += $job
    }
    
    # Ждем завершения пачки и выводим результаты
    $results = $jobs | Wait-Job | Receive-Job
    $results | ForEach-Object { Write-Host $_ }
    $jobs | Remove-Job
    
    $completed += $batchSize
    Write-Host "Completed $completed / $totalRequests requests"
}

Write-Host "Load test completed!"