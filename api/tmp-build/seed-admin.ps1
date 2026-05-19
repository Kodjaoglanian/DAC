# Seed Admin — cria usuário administrador no banco PostgreSQL
# Uso: .\seed-admin.ps1

Write-Host "🌱 Seed Admin — Project Tracker" -ForegroundColor Cyan

# Carrega .env se existir
$envPath = Join-Path $PSScriptRoot ".env"
if (Test-Path $envPath) {
    Get-Content $envPath | ForEach-Object {
        if ($_ -match "^\s*([^#\s=]+)\s*=\s*(.+)\s*$") {
            [System.Environment]::SetEnvironmentVariable($matches[1], $matches[2].Trim())
        }
    }
}

Set-Location $PSScriptRoot
go run cmd\seed\main.go
