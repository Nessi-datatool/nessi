import logging
from datetime import datetime
from pathlib import Path
import json
from typing import Dict, Any, Optional
import uuid

class AuditLogger:
    """Logs security events and user actions for audit purposes."""
    
    def __init__(self, log_dir: str = "logs/audit"):
        self.log_dir = Path(log_dir)
        self.log_dir.mkdir(parents=True, exist_ok=True)
        
        # Setup file handler
        self.log_file = self.log_dir / f"audit_{datetime.now().strftime('%Y%m%d')}.log"
        self.logger = logging.getLogger('audit')
        self.logger.setLevel(logging.INFO)
        
        if not self.logger.handlers:
            handler = logging.FileHandler(self.log_file)
            handler.setFormatter(logging.Formatter(
                '%(asctime)s - %(levelname)s - %(message)s'
            ))
            self.logger.addHandler(handler)
    
    def log_event(self, event_type: str, user: str, details: Dict[str, Any],
                 success: bool = True, ip_address: Optional[str] = None) -> None:
        """Log a security event or user action."""
        event = {
            'id': str(uuid.uuid4()),
            'timestamp': datetime.now().isoformat(),
            'event_type': event_type,
            'user': user,
            'success': success,
            'ip_address': ip_address,
            'details': details
        }
        
        self.logger.info(json.dumps(event))
    
    def log_login(self, user: str, success: bool, ip_address: Optional[str] = None,
                 details: Optional[Dict[str, Any]] = None) -> None:
        """Log a login attempt."""
        self.log_event(
            event_type='login',
            user=user,
            success=success,
            ip_address=ip_address,
            details=details or {}
        )
    
    def log_logout(self, user: str, ip_address: Optional[str] = None) -> None:
        """Log a logout event."""
        self.log_event(
            event_type='logout',
            user=user,
            ip_address=ip_address,
            details={}
        )
    
    def log_access(self, user: str, resource: str, action: str,
                  success: bool, ip_address: Optional[str] = None) -> None:
        """Log resource access."""
        self.log_event(
            event_type='access',
            user=user,
            success=success,
            ip_address=ip_address,
            details={
                'resource': resource,
                'action': action
            }
        )
    
    def log_config_change(self, user: str, config_type: str,
                         changes: Dict[str, Any], ip_address: Optional[str] = None) -> None:
        """Log configuration changes."""
        self.log_event(
            event_type='config_change',
            user=user,
            ip_address=ip_address,
            details={
                'config_type': config_type,
                'changes': changes
            }
        )
    
    def log_security_event(self, user: str, event_type: str,
                          details: Dict[str, Any], ip_address: Optional[str] = None) -> None:
        """Log a security-related event."""
        self.log_event(
            event_type=f'security_{event_type}',
            user=user,
            ip_address=ip_address,
            details=details
        )
    
    def get_events(self, start_time: Optional[datetime] = None,
                  end_time: Optional[datetime] = None,
                  event_type: Optional[str] = None,
                  user: Optional[str] = None) -> list:
        """Retrieve logged events based on filters."""
        events = []
        
        # Get all log files in the date range
        log_files = self._get_log_files(start_time, end_time)
        
        for log_file in log_files:
            with open(log_file) as f:
                for line in f:
                    try:
                        event = json.loads(line)
                        if self._matches_filters(event, event_type, user):
                            events.append(event)
                    except json.JSONDecodeError:
                        continue
        
        return events
    
    def _get_log_files(self, start_time: Optional[datetime],
                      end_time: Optional[datetime]) -> list:
        """Get log files within the specified time range."""
        log_files = []
        
        for log_file in self.log_dir.glob('audit_*.log'):
            file_date = datetime.strptime(log_file.stem.split('_')[1], '%Y%m%d')
            
            if start_time and file_date < start_time.date():
                continue
            if end_time and file_date > end_time.date():
                continue
            
            log_files.append(log_file)
        
        return log_files
    
    def _matches_filters(self, event: Dict[str, Any],
                        event_type: Optional[str],
                        user: Optional[str]) -> bool:
        """Check if an event matches the specified filters."""
        if event_type and event['event_type'] != event_type:
            return False
        if user and event['user'] != user:
            return False
        return True
    
    def _store_event(self, event: Dict[str, Any]) -> None:
        """Store event in audit database."""
        # TODO: Implement database storage
        pass
    
    def generate_report(self, start_time: datetime,
                       end_time: datetime) -> Dict[str, Any]:
        """Generate an audit report for the specified time period."""
        events = self.get_events(start_time, end_time)
        
        report = {
            'period': {
                'start': start_time.isoformat(),
                'end': end_time.isoformat()
            },
            'total_events': len(events),
            'event_types': {},
            'users': {},
            'success_rate': 0,
            'security_events': []
        }
        
        success_count = 0
        for event in events:
            # Count event types
            report['event_types'][event['event_type']] = \
                report['event_types'].get(event['event_type'], 0) + 1
            
            # Count user activity
            report['users'][event['user']] = \
                report['users'].get(event['user'], 0) + 1
            
            # Track success rate
            if event['success']:
                success_count += 1
            
            # Track security events
            if event['event_type'] in ['login_failure', 'access_denied']:
                report['security_events'].append(event)
        
        if events:
            report['success_rate'] = success_count / len(events)
        
        return report 