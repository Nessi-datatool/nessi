"""
Audit Logger for Nessi.dev.
"""

from typing import Dict, List, Any, Optional
from datetime import datetime


class AuditLogger:
    """Audit Logger for Nessi.dev."""

    def __init__(self):
        """Initialize the Audit Logger."""
        self.logs = []

    def log_api_request(self, user_id: str, method: str, path: str, status: str, **kwargs) -> Dict[str, Any]:
        """Log an API request.

        Args:
            user_id: User ID.
            method: HTTP method.
            path: Request path.
            status: Response status.
            **kwargs: Additional log data.

        Returns:
            Dict with log information.
        """
        timestamp = datetime.now().isoformat()
        
        log = {
            'log_id': f"log-{len(self.logs) + 1}",
            'log_type': 'api_request',
            'user_id': user_id,
            'method': method,
            'path': path,
            'status': status,
            'timestamp': timestamp,
            **kwargs,
        }
        
        self.logs.append(log)
        return log

    def log_user_action(self, user_id: str, action: str, resource: str, **kwargs) -> Dict[str, Any]:
        """Log a user action.

        Args:
            user_id: User ID.
            action: Action name.
            resource: Resource name.
            **kwargs: Additional log data.

        Returns:
            Dict with log information.
        """
        timestamp = datetime.now().isoformat()
        
        log = {
            'log_id': f"log-{len(self.logs) + 1}",
            'log_type': 'user_action',
            'user_id': user_id,
            'action': action,
            'resource': resource,
            'timestamp': timestamp,
            **kwargs,
        }
        
        self.logs.append(log)
        return log

    def get_logs(self, log_type: Optional[str] = None, user_id: Optional[str] = None) -> List[Dict[str, Any]]:
        """Get logs.

        Args:
            log_type: Optional log type filter.
            user_id: Optional user ID filter.

        Returns:
            List of logs.
        """
        filtered_logs = self.logs
        
        if log_type:
            filtered_logs = [log for log in filtered_logs if log['log_type'] == log_type]
        
        if user_id:
            filtered_logs = [log for log in filtered_logs if log['user_id'] == user_id]
        
        return filtered_logs
