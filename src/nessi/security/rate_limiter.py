from typing import Dict, List, Optional, Set
import time
import logging
from dataclasses import dataclass
from datetime import datetime, timedelta
import ipaddress

@dataclass
class RateLimit:
    """Rate limit configuration."""
    requests: int
    window: int  # in seconds
    block_duration: int = 300  # in seconds

class RateLimiter:
    """Rate limiter with IP allowlisting support."""
    
    def __init__(self):
        self.logger = logging.getLogger(__name__)
        self.requests: Dict[str, List[float]] = {}  # IP -> list of timestamps
        self.blocked_ips: Dict[str, float] = {}  # IP -> block expiry timestamp
        self.allowed_ips: Set[str] = set()  # Allowlisted IPs
        self.allowed_networks: List[ipaddress.IPv4Network] = []  # Allowlisted networks
    
    def add_to_allowlist(self, ip_or_network: str) -> None:
        """Add an IP or network to the allowlist."""
        try:
            # Try parsing as network first
            network = ipaddress.IPv4Network(ip_or_network)
            self.allowed_networks.append(network)
        except ValueError:
            # If not a network, treat as single IP
            try:
                ip = ipaddress.IPv4Address(ip_or_network)
                self.allowed_ips.add(str(ip))
            except ValueError as e:
                self.logger.error(f"Invalid IP or network: {str(e)}")
                raise ValueError(f"Invalid IP or network: {ip_or_network}")
    
    def remove_from_allowlist(self, ip_or_network: str) -> None:
        """Remove an IP or network from the allowlist."""
        try:
            # Try parsing as network first
            network = ipaddress.IPv4Network(ip_or_network)
            if network in self.allowed_networks:
                self.allowed_networks.remove(network)
        except ValueError:
            # If not a network, treat as single IP
            try:
                ip = ipaddress.IPv4Address(ip_or_network)
                self.allowed_ips.discard(str(ip))
            except ValueError as e:
                self.logger.error(f"Invalid IP or network: {str(e)}")
                raise ValueError(f"Invalid IP or network: {ip_or_network}")
    
    def is_allowed(self, ip: str) -> bool:
        """Check if an IP is allowlisted."""
        try:
            ip_obj = ipaddress.IPv4Address(ip)
            # Check if IP is directly allowlisted
            if str(ip_obj) in self.allowed_ips:
                return True
            # Check if IP belongs to any allowlisted network
            return any(ip_obj in network for network in self.allowed_networks)
        except ValueError as e:
            self.logger.error(f"Invalid IP address: {str(e)}")
            return False
    
    def is_blocked(self, ip: str) -> bool:
        """Check if an IP is currently blocked."""
        if ip in self.blocked_ips:
            if time.time() < self.blocked_ips[ip]:
                return True
            # Block duration expired
            del self.blocked_ips[ip]
        return False
    
    def block_ip(self, ip: str, duration: int) -> None:
        """Block an IP for the specified duration."""
        self.blocked_ips[ip] = time.time() + duration
    
    def check_rate_limit(self, ip: str, rate_limit: RateLimit) -> bool:
        """Check if request is within rate limit."""
        # Allowlisted IPs bypass rate limiting
        if self.is_allowed(ip):
            return True
        
        # Check if IP is blocked
        if self.is_blocked(ip):
            return False
        
        current_time = time.time()
        window_start = current_time - rate_limit.window
        
        # Initialize or update request timestamps
        if ip not in self.requests:
            self.requests[ip] = []
        
        # Remove old timestamps
        self.requests[ip] = [ts for ts in self.requests[ip] if ts > window_start]
        
        # Check rate limit
        if len(self.requests[ip]) >= rate_limit.requests:
            self.logger.warning(f"Rate limit exceeded for IP: {ip}")
            self.block_ip(ip, rate_limit.block_duration)
            return False
        
        # Add current request timestamp
        self.requests[ip].append(current_time)
        return True
    
    def clear_expired(self) -> None:
        """Clear expired blocks and old request records."""
        current_time = time.time()
        
        # Clear expired blocks
        self.blocked_ips = {
            ip: expiry
            for ip, expiry in self.blocked_ips.items()
            if expiry > current_time
        }
        
        # Clear old request records
        for ip in list(self.requests.keys()):
            if not self.requests[ip]:
                del self.requests[ip] 