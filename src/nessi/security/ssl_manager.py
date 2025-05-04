from typing import Optional, Tuple
import logging
import ssl
from pathlib import Path
import os
from datetime import datetime, timedelta
import subprocess
from cryptography import x509
from cryptography.x509.oid import NameOID
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.primitives.asymmetric import rsa

class SSLManager:
    """Manages SSL/TLS certificates and secure communication."""
    
    def __init__(self, cert_dir: str = "certs"):
        self.logger = logging.getLogger(__name__)
        self.cert_dir = Path(cert_dir)
        self.cert_dir.mkdir(exist_ok=True)
    
    def generate_self_signed_cert(
        self,
        common_name: str,
        organization: str = "Nessi.dev",
        country: str = "US",
        validity_days: int = 365
    ) -> Tuple[str, str]:
        """Generate a self-signed SSL certificate."""
        # Generate private key
        private_key = rsa.generate_private_key(
            public_exponent=65537,
            key_size=2048
        )
        
        # Generate public key
        public_key = private_key.public_key()
        
        # Generate certificate
        subject = issuer = x509.Name([
            x509.NameAttribute(NameOID.COUNTRY_NAME, country),
            x509.NameAttribute(NameOID.ORGANIZATION_NAME, organization),
            x509.NameAttribute(NameOID.COMMON_NAME, common_name),
        ])
        
        cert = x509.CertificateBuilder().subject_name(
            subject
        ).issuer_name(
            issuer
        ).public_key(
            public_key
        ).serial_number(
            x509.random_serial_number()
        ).not_valid_before(
            datetime.utcnow()
        ).not_valid_after(
            datetime.utcnow() + timedelta(days=validity_days)
        ).add_extension(
            x509.SubjectAlternativeName([x509.DNSName(common_name)]),
            critical=False,
        ).sign(private_key, hashes.SHA256())
        
        # Write private key
        private_key_path = self.cert_dir / "private.key"
        with open(private_key_path, "wb") as f:
            f.write(private_key.private_bytes(
                encoding=serialization.Encoding.PEM,
                format=serialization.PrivateFormat.PKCS8,
                encryption_algorithm=serialization.NoEncryption()
            ))
        
        # Write certificate
        cert_path = self.cert_dir / "certificate.crt"
        with open(cert_path, "wb") as f:
            f.write(cert.public_bytes(serialization.Encoding.PEM))
        
        return str(cert_path), str(private_key_path)
    
    def create_ssl_context(
        self,
        cert_path: Optional[str] = None,
        key_path: Optional[str] = None,
        verify: bool = True
    ) -> ssl.SSLContext:
        """Create an SSL context for secure communication."""
        context = ssl.create_default_context()
        
        if cert_path and key_path:
            context.load_cert_chain(cert_path, key_path)
        
        if verify:
            context.verify_mode = ssl.CERT_REQUIRED
            context.check_hostname = True
        else:
            context.verify_mode = ssl.CERT_NONE
            context.check_hostname = False
        
        return context
    
    def verify_certificate(self, cert_path: str) -> bool:
        """Verify the validity of a certificate."""
        try:
            with open(cert_path, "rb") as f:
                cert = x509.load_pem_x509_certificate(f.read())
            
            if cert.not_valid_after < datetime.utcnow():
                self.logger.warning("Certificate has expired")
                return False
            
            return True
        except Exception as e:
            self.logger.error(f"Error verifying certificate: {str(e)}")
            return False
    
    def renew_certificate(
        self,
        cert_path: str,
        key_path: str,
        validity_days: int = 365
    ) -> Tuple[str, str]:
        """Renew an existing certificate."""
        if not self.verify_certificate(cert_path):
            # Extract common name from existing certificate
            with open(cert_path, "rb") as f:
                cert = x509.load_pem_x509_certificate(f.read())
                common_name = cert.subject.get_attributes_for_oid(NameOID.COMMON_NAME)[0].value
            
            # Generate new certificate
            return self.generate_self_signed_cert(
                common_name=common_name,
                validity_days=validity_days
            )
        
        return cert_path, key_path 