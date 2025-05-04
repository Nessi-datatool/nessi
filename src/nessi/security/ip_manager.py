from typing import List, Set, Optional
import ipaddress
from pathlib import Path
import json
import logging
from datetime import datetime

logger = logging.getLogger(__name__)

class IPManager:
    """Manages IP allowlisting and access control."""
    
    def __init__(self, config_path: str = "config/ip_allowlist.json"):
        self.config_path = Path(config_path)
        self.allowed_ips: Set[ipaddress.IPv4Address] = set()
        self._load_config()
    
    def _load_config(self) -> None:
        """Load IP allowlist configuration."""
        try:
            if self.config_path.exists():
                with open(self.config_path) as f:
                    config = json.load(f)
                    self.allowed_ips = {
                        ipaddress.IPv4Address(ip) for ip in config.get('allowed_ips', [])
                    }
            else:
                self._save_config()
        except Exception as e:
            logger.error(f"Error loading IP allowlist config: {str(e)}")
            raise
    
    def _save_config(self) -> None:
        """Save IP allowlist configuration."""
        try:
            config = {
                'allowed_ips': [str(ip) for ip in self.allowed_ips],
                'last_updated': datetime.now().isoformat()
            }
            self.config_path.parent.mkdir(parents=True, exist_ok=True)
            with open(self.config_path, 'w') as f:
                json.dump(config, f, indent=2)
        except Exception as e:
            logger.error(f"Error saving IP allowlist config: {str(e)}")
            raise
    
    def is_allowed(self, ip: str) -> bool:
        """Check if an IP address is allowed."""
        try:
            ip_addr = ipaddress.IPv4Address(ip)
            return ip_addr in self.allowed_ips
        except ValueError:
            logger.warning(f"Invalid IP address: {ip}")
            return False
    
    def add_ip(self, ip: str) -> None:
        """Add an IP address to the allowlist."""
        try:
            ip_addr = ipaddress.IPv4Address(ip)
            self.allowed_ips.add(ip_addr)
            self._save_config()
            logger.info(f"Added IP to allowlist: {ip}")
        except ValueError:
            logger.error(f"Invalid IP address: {ip}")
            raise
    
    def remove_ip(self, ip: str) -> None:
        """Remove an IP address from the allowlist."""
        try:
            ip_addr = ipaddress.IPv4Address(ip)
            self.allowed_ips.remove(ip_addr)
            self._save_config()
            logger.info(f"Removed IP from allowlist: {ip}")
        except ValueError:
            logger.error(f"Invalid IP address: {ip}")
            raise
        except KeyError:
            logger.warning(f"IP not in allowlist: {ip}")
            raise
    
    def get_allowed_ips(self) -> List[str]:
        """Get list of allowed IP addresses."""
        return [str(ip) for ip in sorted(self.allowed_ips)]
    
    def clear_allowlist(self) -> None:
        """Clear the IP allowlist."""
        self.allowed_ips.clear()
        self._save_config()
        logger.info("Cleared IP allowlist")
    
    def is_private_ip(self, ip: str) -> bool:
        """Check if an IP address is private."""
        try:
            ip_addr = ipaddress.IPv4Address(ip)
            return ip_addr.is_private
        except ValueError:
            logger.warning(f"Invalid IP address: {ip}")
            return False
    
    def is_loopback_ip(self, ip: str) -> bool:
        """Check if an IP address is loopback."""
        try:
            ip_addr = ipaddress.IPv4Address(ip)
            return ip_addr.is_loopback
        except ValueError:
            logger.warning(f"Invalid IP address: {ip}")
            return False 