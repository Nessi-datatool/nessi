from typing import Any, Optional, Dict
import logging
import json
import redis
from datetime import timedelta
import pickle

class CacheManager:
    """Manages caching operations with Redis."""
    
    def __init__(self, host: str = "localhost", port: int = 6379, db: int = 0):
        self.logger = logging.getLogger(__name__)
        self.redis_client = redis.Redis(
            host=host,
            port=port,
            db=db,
            decode_responses=True
        )
    
    def get(self, key: str) -> Optional[Any]:
        """Get a value from cache."""
        try:
            value = self.redis_client.get(key)
            if value:
                return pickle.loads(value)
            return None
        except Exception as e:
            self.logger.error(f"Error getting value from cache: {str(e)}")
            return None
    
    def set(self, key: str, value: Any, ttl: Optional[int] = None) -> bool:
        """Set a value in cache with optional TTL."""
        try:
            serialized_value = pickle.dumps(value)
            if ttl:
                return self.redis_client.setex(key, ttl, serialized_value)
            return self.redis_client.set(key, serialized_value)
        except Exception as e:
            self.logger.error(f"Error setting value in cache: {str(e)}")
            return False
    
    def delete(self, key: str) -> bool:
        """Delete a value from cache."""
        try:
            return bool(self.redis_client.delete(key))
        except Exception as e:
            self.logger.error(f"Error deleting value from cache: {str(e)}")
            return False
    
    def exists(self, key: str) -> bool:
        """Check if a key exists in cache."""
        try:
            return bool(self.redis_client.exists(key))
        except Exception as e:
            self.logger.error(f"Error checking key existence: {str(e)}")
            return False
    
    def clear(self) -> bool:
        """Clear all values from cache."""
        try:
            return bool(self.redis_client.flushdb())
        except Exception as e:
            self.logger.error(f"Error clearing cache: {str(e)}")
            return False
    
    def get_stats(self) -> Dict[str, Any]:
        """Get cache statistics."""
        try:
            info = self.redis_client.info()
            return {
                "used_memory": info.get("used_memory_human", "0"),
                "total_keys": info.get("db0", {}).get("keys", 0),
                "uptime": info.get("uptime_in_seconds", 0),
                "connected_clients": info.get("connected_clients", 0)
            }
        except Exception as e:
            self.logger.error(f"Error getting cache stats: {str(e)}")
            return {}

class CacheDecorator:
    """Decorator for caching function results."""
    
    def __init__(self, cache_manager: CacheManager, ttl: Optional[int] = None):
        self.cache_manager = cache_manager
        self.ttl = ttl
    
    def __call__(self, func):
        def wrapper(*args, **kwargs):
            # Generate cache key from function name and arguments
            key = f"{func.__name__}:{str(args)}:{str(kwargs)}"
            
            # Try to get from cache
            cached_result = self.cache_manager.get(key)
            if cached_result is not None:
                return cached_result
            
            # Execute function and cache result
            result = func(*args, **kwargs)
            self.cache_manager.set(key, result, self.ttl)
            return result
        return wrapper 