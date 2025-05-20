"""
Cloud Manager for Nessi.dev.
"""

from typing import Dict, List, Any, Optional
from .cloud_provider import CloudProvider


class CloudManager:
    """Cloud Manager for Nessi.dev."""

    def __init__(self):
        """Initialize the Cloud Manager."""
        self.providers = {}

    def get_provider(self, provider_type: str) -> CloudProvider:
        """Get a cloud provider.

        Args:
            provider_type: Type of cloud provider (aws, azure, gcp).

        Returns:
            Cloud provider instance.

        Raises:
            Exception: If provider type is not supported.
        """
        if provider_type in self.providers:
            return self.providers[provider_type]
        
        if provider_type not in ['aws', 'azure', 'gcp']:
            raise Exception(f"Unsupported provider type: {provider_type}")
        
        provider = CloudProvider(provider_type=provider_type)
        self.providers[provider_type] = provider
        
        return provider

    def register_provider(self, provider_type: str, provider: CloudProvider) -> None:
        """Register a cloud provider.

        Args:
            provider_type: Type of cloud provider (aws, azure, gcp).
            provider: Cloud provider instance.
        """
        self.providers[provider_type] = provider
