# PowerShell script to run tests with the required environment variable
$env:ENCRYPTION_KEY="0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

Write-Host "Running all backend tests..." -ForegroundColor Green
go test ./... -v

Write-Host "`nRunning tests with coverage..." -ForegroundColor Green
go test ./... -cover

Write-Host "`nRunning benchmarks..." -ForegroundColor Green
go test ./... -bench=.
