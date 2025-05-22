# Nessi Installation Guide

This guide provides detailed instructions for installing and setting up Nessi on various platforms.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation Methods](#installation-methods)
  - [Using the Installation Script](#using-the-installation-script)
  - [Using Pre-built Binaries](#using-pre-built-binaries)
  - [Building from Source](#building-from-source)
  - [Using Docker](#using-docker)
- [Platform-Specific Instructions](#platform-specific-instructions)
  - [Linux](#linux)
  - [macOS](#macos)
  - [Windows](#windows)
- [Verifying Installation](#verifying-installation)
- [Configuration](#configuration)
- [Troubleshooting](#troubleshooting)
- [Upgrading](#upgrading)
- [Uninstallation](#uninstallation)

## Prerequisites

Before installing Nessi, ensure your system meets the following requirements:

- **Operating System**: Linux, macOS, or Windows
- **Disk Space**: At least 200 MB of free disk space
- **Memory**: Minimum 1 GB RAM (4 GB recommended for larger datasets)
- **Go**: Version 1.21 or later (only required for building from source)
- **Git**: Latest version (only required for building from source)

For premium features, the following additional requirements apply:

- **Docker**: Latest version (for containerized deployments)
- **AWS CLI**: Configured with appropriate credentials (for AWS S3 integration)
- **Databricks CLI**: Configured with appropriate credentials (for Databricks integration)

## Installation Methods

### Using the Installation Script

The easiest way to install Nessi is using our installation script, which automatically detects your operating system and installs the appropriate version:

```bash
curl -sSL https://nessi.dev/install.sh | bash
```

For a specific version:

```bash
curl -sSL https://nessi.dev/install.sh | bash -s -- --version v1.2.3
```

For a non-root installation:

```bash
curl -sSL https://nessi.dev/install.sh | bash -s -- --dir $HOME/.local/bin
```

### Using Pre-built Binaries

1. Download the appropriate binary for your platform from the [releases page](https://github.com/nessi-dev/nessi/releases).

2. Extract the archive:

   ```bash
   # Linux/macOS
   tar -xzf nessi_1.2.3_linux_amd64.tar.gz
   
   # Windows
   Expand-Archive -Path nessi_1.2.3_windows_amd64.zip -DestinationPath .
   ```

3. Move the binary to a directory in your PATH:

   ```bash
   # Linux/macOS
   sudo mv nessi /usr/local/bin/
   
   # Windows
   # Move nessi.exe to a directory in your PATH
   ```

4. Make the binary executable (Linux/macOS only):

   ```bash
   sudo chmod +x /usr/local/bin/nessi
   ```

### Building from Source

To build Nessi from source, you need Go 1.21 or later and Git installed on your system.

1. Clone the repository:

   ```bash
   git clone https://github.com/nessi-dev/nessi.git
   cd nessi
   ```

2. Build the binary:

   ```bash
   go build -o nessi ./cmd/nessi
   ```

3. Move the binary to a directory in your PATH:

   ```bash
   # Linux/macOS
   sudo mv nessi /usr/local/bin/
   
   # Windows
   # Move nessi.exe to a directory in your PATH
   ```

### Using Docker

Nessi is also available as a Docker image:

```bash
# Pull the latest image
docker pull nessi/nessi:latest

# Run Nessi
docker run --rm -v $(pwd):/data nessi/nessi:latest --help
```

For a specific version:

```bash
docker pull nessi/nessi:1.2.3
```

## Platform-Specific Instructions

### Linux

#### Debian/Ubuntu

```bash
# Install dependencies
sudo apt-get update
sudo apt-get install -y curl

# Install Nessi
curl -sSL https://nessi.dev/install.sh | bash
```

#### Red Hat/CentOS/Fedora

```bash
# Install dependencies
sudo yum install -y curl

# Install Nessi
curl -sSL https://nessi.dev/install.sh | bash
```

#### Arch Linux

```bash
# Install from AUR
yay -S nessi

# Or manually
curl -sSL https://nessi.dev/install.sh | bash
```

### macOS

#### Using Homebrew

```bash
# Install Homebrew if not already installed
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install Nessi
brew tap nessi-dev/nessi
brew install nessi
```

#### Manual Installation

```bash
# Install Nessi
curl -sSL https://nessi.dev/install.sh | bash
```

### Windows

#### Using Chocolatey

```powershell
# Install Chocolatey if not already installed
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

# Install Nessi
choco install nessi
```

#### Using Scoop

```powershell
# Install Scoop if not already installed
iwr -useb get.scoop.sh | iex

# Add the bucket and install Nessi
scoop bucket add nessi https://github.com/nessi-dev/scoop-bucket.git
scoop install nessi
```

#### Manual Installation

1. Download the Windows binary from the [releases page](https://github.com/nessi-dev/nessi/releases).
2. Extract the ZIP file.
3. Add the directory containing `nessi.exe` to your PATH.

## Verifying Installation

To verify that Nessi has been installed correctly, run:

```bash
nessi version
```

You should see output similar to:

```
Nessi CLI v1.2.3
Build Date: 2025-05-01
Git Commit: abcdef123456
Go Version: go1.21.0
OS/Arch: linux/amd64
```

To check if all dependencies are correctly installed and configured:

```bash
nessi info
```

## Configuration

After installation, you can configure Nessi using the configuration commands:

```bash
# Initialize configuration
nessi config init --interactive

# Set specific configuration values
nessi config set log.level info
nessi config set delta.default_path /path/to/delta
```

For detailed configuration options, see the [Configuration documentation](cli/config.md).

### Environment Variables

You can also configure Nessi using environment variables:

```bash
# Set log level
export NESSI_LOG_LEVEL=debug

# Set license path
export NESSI_LICENSE_PATH=/path/to/license.json
```

For a complete list of environment variables, see the [CLI Reference](cli/REFERENCE.md#environment-variables).

## Troubleshooting

### Common Issues

#### Permission Denied

If you encounter permission issues when running Nessi:

```bash
# Make sure the binary is executable
chmod +x /path/to/nessi

# Or run with sudo (not recommended for regular use)
sudo nessi --help
```

#### Command Not Found

If you get a "command not found" error:

```bash
# Check if the binary is in your PATH
which nessi

# If not, add it to your PATH
export PATH=$PATH:/path/to/nessi/directory
```

#### Missing Dependencies

If Nessi reports missing dependencies:

```bash
# Check system information
nessi info

# Install required dependencies based on the output
```

### Logging

To enable debug logging:

```bash
nessi --log-level debug <command>
```

Or using environment variables:

```bash
export NESSI_LOG_LEVEL=debug
nessi <command>
```

### Getting Help

If you encounter issues not covered here:

- Check the [Frequently Asked Questions](FREQUENTLY_ASKED_QUESTIONS.md)
- Visit the [Nessi website](https://nessi.dev) for documentation
- Join the [Nessi community forum](https://community.nessi.dev) for support
- Report issues on [GitHub](https://github.com/nessi-dev/nessi/issues)

## Upgrading

To upgrade Nessi to the latest version:

```bash
# Using the installation script
curl -sSL https://nessi.dev/install.sh | bash

# Using Homebrew on macOS
brew upgrade nessi

# Using Chocolatey on Windows
choco upgrade nessi

# Using Scoop on Windows
scoop update nessi

# Using Docker
docker pull nessi/nessi:latest
```

To upgrade to a specific version:

```bash
# Using the installation script
curl -sSL https://nessi.dev/install.sh | bash -s -- --version v1.2.3

# Using Docker
docker pull nessi/nessi:1.2.3
```

## Uninstallation

To uninstall Nessi:

### Linux/macOS

```bash
# If installed using the installation script
sudo rm /usr/local/bin/nessi

# If installed using Homebrew
brew uninstall nessi
```

### Windows

```powershell
# If installed using Chocolatey
choco uninstall nessi

# If installed using Scoop
scoop uninstall nessi

# If installed manually
# Remove nessi.exe from your PATH
```

### Docker

```bash
# Remove the Docker image
docker rmi nessi/nessi:latest
```

### Configuration Files

To remove configuration files:

```bash
# Remove global configuration
rm -rf ~/.nessi

# Remove project-specific configuration
rm -f ./.nessi.yaml
```
