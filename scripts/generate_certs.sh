#!/bin/bash

# Exit on any error
set -e

# Default output directory
CERTS_DIR=${1:-"certs"}

# Create output directory if it doesn't exist
mkdir -p "$CERTS_DIR"

# Generate private key
openssl genrsa -out "$CERTS_DIR/server.key" 4096

# Generate Certificate Signing Request (CSR)
openssl req -new -key "$CERTS_DIR/server.key" -out "$CERTS_DIR/server.csr" -subj "/C=US/ST=State/L=City/O=Nessi.dev/CN=localhost"

# Generate self-signed certificate
openssl x509 -req -days 365 -in "$CERTS_DIR/server.csr" -signkey "$CERTS_DIR/server.key" -out "$CERTS_DIR/server.crt"

# Generate Diffie-Hellman parameters for perfect forward secrecy
openssl dhparam -out "$CERTS_DIR/dhparam.pem" 2048

# Set proper permissions
chmod 600 "$CERTS_DIR/server.key"
chmod 644 "$CERTS_DIR/server.crt"
chmod 644 "$CERTS_DIR/dhparam.pem"

# Clean up CSR
rm "$CERTS_DIR/server.csr"

echo "Self-signed certificates generated successfully in $CERTS_DIR:"
ls -l "$CERTS_DIR"
echo
echo "Note: These certificates are for development only. Use proper certificates for production." 