"""
NESSI - FREE FOR PERSONAL USE LICENSE

Copyright (c) 2025 nessi.dev. All rights reserved.

nessi.dev is free for personal use.

A paid license is required for enterprise or consulting use.

This version is free. Future versions of nessi.dev will require a license for all users.
"""

from typing import Optional
from cryptography.fernet import Fernet
import base64
import os

class Encryption:
    """Handles data encryption and decryption."""
    
    def __init__(self, key: Optional[str] = None):
        """Initialize the encryption system.
        
        Args:
            key: Encryption key (if None, a new one will be generated)
        """
        if key is None:
            self.key = Fernet.generate_key()
        else:
            self.key = base64.urlsafe_b64encode(key.encode())
            
        self.cipher_suite = Fernet(self.key)
        
    def encrypt(self, data: str) -> str:
        """Encrypt data.
        
        Args:
            data: Data to encrypt
            
        Returns:
            Encrypted data as a string
        """
        return self.cipher_suite.encrypt(data.encode()).decode()
        
    def decrypt(self, encrypted_data: str) -> str:
        """Decrypt data.
        
        Args:
            encrypted_data: Encrypted data to decrypt
            
        Returns:
            Decrypted data as a string
        """
        return self.cipher_suite.decrypt(encrypted_data.encode()).decode()
        
    def rotate_key(self) -> None:
        """Generate a new encryption key."""
        self.key = Fernet.generate_key()
        self.cipher_suite = Fernet(self.key) 