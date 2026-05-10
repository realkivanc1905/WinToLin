$ErrorActionPreference = 'Stop'

$repo = "realkivanc1905/WinToLin"
Write-Host "Downloading WinToLin ($repo)..."

try {
    $release = Invoke-RestMethod -Uri "https://api.github.com/repos/$repo/releases/latest"
} catch {
    Write-Host "ERROR: Could not fetch latest release. Please ensure you have created a Release on GitHub." -ForegroundColor Red
    return
}

$asset = $release.assets | Where-Object { $_.name -eq "wintolin-windows-amd64.exe" }

if (-not $asset) {
    Write-Host "ERROR: Windows binary not found in the latest release." -ForegroundColor Red
    return
}

$url = $asset.browser_download_url
$installDir = "$env:USERPROFILE\wintolin"
if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir | Out-Null
}

$dest = "$installDir\wintolin.exe"
Write-Host "Downloading: $url -> $dest"
Invoke-WebRequest -Uri $url -OutFile $dest

$userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($userPath -notlike "*$installDir*") {
    $newPath = "$userPath;$installDir"
    [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
    Write-Host "$installDir directory added to User PATH." -ForegroundColor Cyan
}

Write-Host "Installation complete! Please restart your terminal and run the 'wintolin' command." -ForegroundColor Green
