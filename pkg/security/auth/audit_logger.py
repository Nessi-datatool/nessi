"""
Audit Logger for Nessi.dev.
"""

from typing import Dict, List, Any, Optional
from datetime import datetime


class AuditLogger:
    """Audit Logger for Nessi.dev."""

    def __init__(self, log_file=None):
        """Initialize the Audit Logger.

        Args:
            log_file: Path to the log file.
        """
        self.log_file = log_file
        self.logs = []

    def log_api_request(self, user_id: str, method: str, path: str, status: str) -> Dict[str, Any]:
        """Log an API request.

        Args:
            user_id: User ID.
            method: HTTP method.
            path: Request path.
            status: Request status.

        Returns:
            Dict with log entry details.
        """
        timestamp = datetime.now().isoformat()
        log_entry = {
            'timestamp': timestamp,
            'user_id': user_id,
            'method': method,
            'path': path,
            'status': status,
            'log_id': f'log-{timestamp}-{user_id}',
        }
        
        self.logs.append(log_entry)
        
        # Write to log file if specified
        if self.log_file:
            with open(self.log_file, 'a') as f:
                f.write(f"{timestamp} | {user_id} | {method} | {path} | {status}\n")
        
        return log_entry

    def log_user_action(self, user_id: str, action: str, details: Dict[str, Any]) -> Dict[str, Any]:
        """Log a user action.

        Args:
            user_id: User ID.
            action: Action name.
            details: Action details.

        Returns:
            Dict with log entry details.
        """
        timestamp = datetime.now().isoformat()
        log_entry = {
            'timestamp': timestamp,
            'user_id': user_id,
            'action': action,
            'details': details,
            'log_id': f'log-{timestamp}-{user_id}',
        }
        
        self.logs.append(log_entry)
        
        # Write to log file if specified
        if self.log_file:
            with open(self.log_file, 'a') as f:
                f.write(f"{timestamp} | {user_id} | {action} | {str(details)}\n")
        
        return log_entry

    def get_logs(self, user_id: Optional[str] = None, action: Optional[str] = None) -> List[Dict[str, Any]]:
        """Get logs, optionally filtered by user ID or action.

        Args:
            user_id: User ID to filter by.
            action: Action to filter by.

        Returns:
            List of log entries.
        """
        filtered_logs = self.logs
        
        if user_id:
            filtered_logs = [log for log in filtered_logs if log.get('user_id') == user_id]
        
        if action:
            filtered_logs = [log for log in filtered_logs if log.get('action') == action]
        
        return filtered_logs
