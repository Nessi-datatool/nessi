# Nessi PowerShell Installation Script for Windows

# Function to display colorful output
function Write-ColorOutput {
    param (
        [Parameter(Mandatory=$true)]
        [string]$Message,
        
        [Parameter(Mandatory=$false)]
        [string]$ForegroundColor = "White"
    )
    
    $originalColor = $host.UI.RawUI.ForegroundColor
    $host.UI.RawUI.ForegroundColor = $ForegroundColor
    Write-Output $Message
    $host.UI.RawUI.ForegroundColor = $originalColor
}

# Function to display the logo
function Show-Logo {
    Write-ColorOutput "  _   _               _    " "Blue"
    Write-ColorOutput " | \ | |             (_)   " "Blue"
    Write-ColorOutput " |  \| | ___  ___ ___ _  " "Blue"
    Write-ColorOutput " | . \  |/ _ \/ __/ __| | " "Blue"
    Write-ColorOutput " | |\  |  __/\__ \__ \ | " "Blue"
    Write-ColorOutput " |_| \_|\___||___/___/_| " "Blue"
    Write-ColorOutput "Delta Lake Quality Tool" "Cyan"
    Write-Output ""
}

# Function to display section headers
function Show-Header {
    param (
        [Parameter(Mandatory=$true)]
        [string]$Title
    )
    
    Write-Output ""
    Write-ColorOutput "==> $Title" "Magenta"
    Write-Output ""
}

# Function to check if running as administrator
function Test-Administrator {
    $user = [Security.Principal.WindowsIdentity]::GetCurrent();
    $principal = New-Object Security.Principal.WindowsPrincipal $user
    return $principal.IsInRole([Security.Principal.WindowsBuiltinRole]::Administrator)
}

# Main installation script
Show-Logo

# Check if running as administrator
if (-not (Test-Administrator)) {
    Write-ColorOutput "Warning: Not running as Administrator. Installation may fail or require additional permissions." "Yellow"
    Write-Output "Consider restarting this script as Administrator if you encounter permission issues."
    Write-Output ""
}

# Check system architecture
Show-Header "Checking system architecture"
$architecture = $env:PROCESSOR_ARCHITECTURE
if ($architecture -eq "AMD64") {
    $arch = "amd64"
    Write-ColorOutput "✓ Detected 64-bit system" "Green"
} elseif ($architecture -eq "ARM64") {
    $arch = "arm64"
    Write-ColorOutput "✓ Detected ARM64 system" "Green"
} else {
    Write-ColorOutput "✗ Unsupported architecture: $architecture" "Red"
    Write-Output "Please download manually from https://github.com/nessi-dev/nessi/releases"
    exit 1
}

# Get latest version
Show-Header "Determining latest version"
try {
    $latestRelease = Invoke-RestMethod -Uri "https://api.github.com/repos/nessi-dev/nessi/releases/latest" -ErrorAction Stop
    $version = $latestRelease.tag_name
    Write-ColorOutput "✓ Latest version: $version" "Green"
} catch {
    $version = "v1.0.0"
    Write-ColorOutput "! Could not determine latest version. Using default $version" "Yellow"
}

# Create temporary directory
$tempDir = [System.IO.Path]::GetTempPath() + [System.Guid]::NewGuid().ToString()
New-Item -ItemType Directory -Path $tempDir | Out-Null

# Download binary
Show-Header "Downloading Nessi $version"
$binaryName = "nessi_${version}_windows_${arch}.zip"
$downloadUrl = "https://github.com/nessi-dev/nessi/releases/download/${version}/${binaryName}"
Write-Output "Downloading from: $downloadUrl"

$downloadPath = "$tempDir\$binaryName"
try {
    Invoke-WebRequest -Uri $downloadUrl -OutFile $downloadPath -ErrorAction Stop
    Write-ColorOutput "✓ Downloaded to $downloadPath" "Green"
} catch {
    Write-ColorOutput "✗ Failed to download Nessi: $_" "Red"
    exit 1
}

# Extract the binary
Show-Header "Extracting binary"
try {
    Expand-Archive -Path $downloadPath -DestinationPath $tempDir -ErrorAction Stop
    Write-ColorOutput "✓ Extracted binary" "Green"
} catch {
    Write-ColorOutput "✗ Failed to extract archive: $_" "Red"
    exit 1
}

# Determine install location
Show-Header "Installing Nessi"
$installDir = "$env:LOCALAPPDATA\Nessi"
if (Test-Administrator) {
    $installDir = "$env:ProgramFiles\Nessi"
}

# Create install directory if it doesn't exist
if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir | Out-Null
}

# Copy binary to install location
Copy-Item -Path "$tempDir\nessi.exe" -Destination "$installDir\nessi.exe" -Force
Write-ColorOutput "✓ Installed Nessi to $installDir\nessi.exe" "Green"

# Add to PATH if not already there
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if (-not $userPath.Contains($installDir)) {
    Show-Header "Updating PATH"
    [Environment]::SetEnvironmentVariable("Path", $userPath + ";$installDir", "User")
    Write-ColorOutput "✓ Added Nessi to your PATH" "Green"
    Write-ColorOutput "! Note: You'll need to restart your terminal for the PATH change to take effect" "Yellow"
}

# Clean up
Remove-Item -Path $tempDir -Recurse -Force

# Setup PowerShell completion
Show-Header "Setting up PowerShell completion"
$profileDir = Split-Path -Parent $PROFILE
if (-not (Test-Path $profileDir)) {
    New-Item -ItemType Directory -Path $profileDir -Force | Out-Null
}

$completionPath = "$profileDir\nessi_completion.ps1"
& "$installDir\nessi.exe" completion powershell > $completionPath

# Add completion to profile if not already there
$profileContent = ""
if (Test-Path $PROFILE) {
    $profileContent = Get-Content $PROFILE -Raw
}

if (-not $profileContent.Contains("nessi_completion.ps1")) {
    $completionCommand = ". $completionPath"
    Add-Content -Path $PROFILE -Value "`n# Nessi completion`n$completionCommand"
    Write-ColorOutput "✓ Added completion to PowerShell profile" "Green"
} else {
    Write-ColorOutput "✓ Completion already in PowerShell profile" "Green"
}

# Final instructions
Show-Header "Installation complete!"

Write-ColorOutput "Nessi has been successfully installed!" "Green"
Write-Output ""
Write-Output "To verify the installation, run:"
Write-Output "  nessi version"
Write-Output ""
Write-Output "To get started, run:"
Write-Output "  nessi help"
Write-Output ""
Write-Output "Documentation: https://github.com/nessi-dev/nessi/docs"
Write-Output "Report issues: https://github.com/nessi-dev/nessi/issues"
Write-Output ""
Write-ColorOutput "If you find Nessi useful, please consider starring the repository on GitHub!" "Yellow"
Write-Output "  https://github.com/nessi-dev/nessi"
