import bcrypt
import secrets
import string
import re
from typing import Dict, Any, Optional
import logging
from datetime import datetime, timedelta
from ..config.security_config import SecurityConfig

logger = logging.getLogger(__name__)

class PasswordManager:
    """Utility for secure password management."""
    
    def __init__(self, config: SecurityConfig):
        self.config = config
    
    def hash_password(self, password: str) -> str:
        """Hash a password using bcrypt."""
        try:
            salt = bcrypt.gensalt()
            hashed = bcrypt.hashpw(password.encode('utf-8'), salt)
            return hashed.decode('utf-8')
        except Exception as e:
            logger.error(f"Error hashing password: {str(e)}")
            raise
    
    def verify_password(self, password: str, hashed: str) -> bool:
        """Verify a password against its hash."""
        try:
            return bcrypt.checkpw(
                password.encode('utf-8'),
                hashed.encode('utf-8')
            )
        except Exception as e:
            logger.error(f"Error verifying password: {str(e)}")
            return False
    
    def generate_password(self) -> str:
        """Generate a secure random password."""
        try:
            policy = self.config.password_policy
            min_length = policy.get('min_length', 12)
            
            # Define character sets
            uppercase = string.ascii_uppercase
            lowercase = string.ascii_lowercase
            numbers = string.digits
            special = string.punctuation
            
            # Ensure at least one character from each required set
            password = []
            if policy.get('require_uppercase', True):
                password.append(secrets.choice(uppercase))
            if policy.get('require_lowercase', True):
                password.append(secrets.choice(lowercase))
            if policy.get('require_numbers', True):
                password.append(secrets.choice(numbers))
            if policy.get('require_special_chars', True):
                password.append(secrets.choice(special))
            
            # Fill the rest with random characters
            all_chars = []
            if policy.get('require_uppercase', True):
                all_chars.extend(uppercase)
            if policy.get('require_lowercase', True):
                all_chars.extend(lowercase)
            if policy.get('require_numbers', True):
                all_chars.extend(numbers)
            if policy.get('require_special_chars', True):
                all_chars.extend(special)
            
            while len(password) < min_length:
                password.append(secrets.choice(all_chars))
            
            # Shuffle the password
            secrets.SystemRandom().shuffle(password)
            
            return ''.join(password)
        
        except Exception as e:
            logger.error(f"Error generating password: {str(e)}")
            raise
    
    def validate_password(self, password: str) -> Dict[str, Any]:
        """Validate a password against the password policy."""
        try:
            policy = self.config.password_policy
            errors = []
            
            # Check length
            if len(password) < policy.get('min_length', 12):
                errors.append(f"Password must be at least {policy['min_length']} characters long")
            
            # Check uppercase
            if policy.get('require_uppercase', True) and not re.search(r'[A-Z]', password):
                errors.append("Password must contain at least one uppercase letter")
            
            # Check lowercase
            if policy.get('require_lowercase', True) and not re.search(r'[a-z]', password):
                errors.append("Password must contain at least one lowercase letter")
            
            # Check numbers
            if policy.get('require_numbers', True) and not re.search(r'[0-9]', password):
                errors.append("Password must contain at least one number")
            
            # Check special characters
            if policy.get('require_special_chars', True) and not re.search(r'[^A-Za-z0-9]', password):
                errors.append("Password must contain at least one special character")
            
            return {
                'valid': len(errors) == 0,
                'errors': errors
            }
        
        except Exception as e:
            logger.error(f"Error validating password: {str(e)}")
            raise
    
    def check_password_expiry(self, last_changed: datetime) -> bool:
        """Check if a password has expired."""
        try:
            expiry_days = self.config.password_policy.get('expiry_days', 90)
            expiry_date = last_changed + timedelta(days=expiry_days)
            return datetime.utcnow() > expiry_date
        except Exception as e:
            logger.error(f"Error checking password expiry: {str(e)}")
            return False 