#!/bin/bash

# Nessi installation script
# This script installs Nessi on Linux and macOS systems

set -e

# Color codes for pretty output
RED="\033[0;31m"
GREEN="\033[0;32m"
YELLOW="\033[0;33m"
BLUE="\033[0;34m"
MAGENTA="\033[0;35m"
CYAN="\033[0;36m"
NC="\033[0m" # No Color

# Detect OS
OS="$(uname -s)"
ARCH="$(uname -m)"

print_logo() {
  echo -e "${BLUE}"
  echo "  _   _               _    "
  echo " | \ | |             (_)   "
  echo " |  \| | ___  ___ ___ _  "
  echo " | . \  |/ _ \/ __/ __| | "
  echo " | |\  |  __/\__ \__ \ | "
  echo " |_| \_|\___||___/___/_| "
  echo -e "${NC}"
  echo -e "${CYAN}Delta Lake Quality Tool${NC}"
  echo ""
}

print_header() {
  echo -e "\n${MAGENTA}==>${NC} ${CYAN}$1${NC}\n"
}

print_success() {
  echo -e "${GREEN}✓${NC} $1"
}

print_error() {
  echo -e "${RED}✗${NC} $1"
}

print_warning() {
  echo -e "${YELLOW}!${NC} $1"
}

check_dependencies() {
  print_header "Checking dependencies"
  
  # Check for curl or wget
  if command -v curl &>/dev/null; then
    DOWNLOAD_CMD="curl -L -o"
    print_success "curl found"
  elif command -v wget &>/dev/null; then
    DOWNLOAD_CMD="wget -O"
    print_success "wget found"
  else
    print_error "Neither curl nor wget found. Please install one of them and try again."
    exit 1
  fi
  
  # Check for tar
  if command -v tar &>/dev/null; then
    print_success "tar found"
  else
    print_error "tar not found. Please install tar and try again."
    exit 1
  fi
}

get_latest_version() {
  print_header "Determining latest version"
  
  if [ "$DOWNLOAD_CMD" = "curl -L -o" ]; then
    VERSION=$(curl -s https://api.github.com/repos/nessi-dev/nessi/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
  else
    VERSION=$(wget -q -O- https://api.github.com/repos/nessi-dev/nessi/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
  fi
  
  if [ -z "$VERSION" ]; then
    print_warning "Could not determine latest version. Using default v1.0.0"
    VERSION="v1.0.0"
  else
    print_success "Latest version: $VERSION"
  fi
}

download_binary() {
  print_header "Downloading Nessi $VERSION"
  
  # Determine binary name based on OS and architecture
  if [ "$OS" = "Darwin" ]; then
    if [ "$ARCH" = "arm64" ]; then
      BINARY="nessi_${VERSION}_darwin_arm64.tar.gz"
    else
      BINARY="nessi_${VERSION}_darwin_amd64.tar.gz"
    fi
  elif [ "$OS" = "Linux" ]; then
    if [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
      BINARY="nessi_${VERSION}_linux_arm64.tar.gz"
    else
      BINARY="nessi_${VERSION}_linux_amd64.tar.gz"
    fi
  else
    print_error "Unsupported operating system: $OS"
    echo "Please download manually from https://github.com/nessi-dev/nessi/releases"
    exit 1
  fi
  
  URL="https://github.com/nessi-dev/nessi/releases/download/${VERSION}/${BINARY}"
  echo "Downloading from: $URL"
  
  TEMP_DIR=$(mktemp -d)
  $DOWNLOAD_CMD "$TEMP_DIR/$BINARY" "$URL"
  
  print_success "Downloaded to $TEMP_DIR/$BINARY"
  
  # Extract the binary
  print_header "Extracting binary"
  tar -xzf "$TEMP_DIR/$BINARY" -C "$TEMP_DIR"
  print_success "Extracted binary"
}

install_binary() {
  print_header "Installing Nessi"
  
  # Determine install location
  INSTALL_DIR="/usr/local/bin"
  if [ ! -w "$INSTALL_DIR" ]; then
    print_warning "Cannot write to $INSTALL_DIR, trying alternate location"
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"
  fi
  
  # Move binary to install location
  mv "$TEMP_DIR/nessi" "$INSTALL_DIR/nessi"
  chmod +x "$INSTALL_DIR/nessi"
  
  print_success "Installed Nessi to $INSTALL_DIR/nessi"
  
  # Clean up
  rm -rf "$TEMP_DIR"
  
  # Check if install location is in PATH
  if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    print_warning "$INSTALL_DIR is not in your PATH"
    echo "Add the following to your shell profile (.bashrc, .zshrc, etc.):"
    echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
  fi
}

setup_autocomplete() {
  print_header "Setting up shell completion"
  
  SHELL_TYPE="$(basename "$SHELL")"
  COMPLETION_DIR=""
  
  case "$SHELL_TYPE" in
    bash)
      if [ -d "/etc/bash_completion.d" ] && [ -w "/etc/bash_completion.d" ]; then
        COMPLETION_DIR="/etc/bash_completion.d"
      elif [ -d "$HOME/.local/share/bash-completion/completions" ]; then
        COMPLETION_DIR="$HOME/.local/share/bash-completion/completions"
        mkdir -p "$COMPLETION_DIR"
      else
        mkdir -p "$HOME/.bash_completion.d"
        COMPLETION_DIR="$HOME/.bash_completion.d"
        if ! grep -q "$COMPLETION_DIR" "$HOME/.bashrc" 2>/dev/null; then
          echo "[ -d $COMPLETION_DIR ] && for f in $COMPLETION_DIR/*; do source \$f; done" >> "$HOME/.bashrc"
        fi
      fi
      "$INSTALL_DIR/nessi" completion bash > "$COMPLETION_DIR/nessi"
      print_success "Bash completion installed to $COMPLETION_DIR/nessi"
      ;;
    zsh)
      if [ -d "$HOME/.zsh/completion" ]; then
        COMPLETION_DIR="$HOME/.zsh/completion"
      else
        mkdir -p "$HOME/.zsh/completion"
        COMPLETION_DIR="$HOME/.zsh/completion"
        if ! grep -q "$COMPLETION_DIR" "$HOME/.zshrc" 2>/dev/null; then
          echo "fpath=($COMPLETION_DIR \$fpath)" >> "$HOME/.zshrc"
          echo "autoload -U compinit && compinit" >> "$HOME/.zshrc"
        fi
      fi
      "$INSTALL_DIR/nessi" completion zsh > "$COMPLETION_DIR/_nessi"
      print_success "Zsh completion installed to $COMPLETION_DIR/_nessi"
      ;;
    *)
      print_warning "Shell completion not supported for $SHELL_TYPE"
      echo "To generate completion for your shell, run:"
      echo "  nessi completion [bash|zsh|fish|powershell]"
      ;;
  esac
}

print_final_instructions() {
  print_header "Installation complete!"
  
  echo -e "${GREEN}Nessi has been successfully installed!${NC}"
  echo ""
  echo "To verify the installation, run:"
  echo "  nessi version"
  echo ""
  echo "To get started, run:"
  echo "  nessi help"
  echo ""
  echo "Documentation: https://github.com/nessi-dev/nessi/docs"
  echo "Report issues: https://github.com/nessi-dev/nessi/issues"
  echo ""
  echo -e "${YELLOW}If you find Nessi useful, please consider starring the repository on GitHub!${NC}"
  echo "  https://github.com/nessi-dev/nessi"
}

# Main installation flow
print_logo
check_dependencies
get_latest_version
download_binary
install_binary
setup_autocomplete
print_final_instructions
