"""Webhook support and third-party integrations."""
import logging
from typing import Dict, List, Optional, Any, Callable
from datetime import datetime
import json
import os
from dataclasses import dataclass
from enum import Enum
import requests
from flask import Flask, request, Response

logger = logging.getLogger(__name__)

class WebhookMethod(Enum):
    """HTTP methods for webhooks."""
    GET = "GET"
    POST = "POST"
    PUT = "PUT"
    DELETE = "DELETE"

@dataclass
class Webhook:
    """Webhook definition."""
    name: str
    url: str
    method: WebhookMethod
    headers: Dict[str, str]
    payload_template: str
    description: str

class WebhookManager:
    """Webhook management system."""
    def __init__(self, webhooks_dir: str = "webhooks"):
        self.webhooks: Dict[str, Webhook] = {}
        self.webhooks_dir = webhooks_dir
        self.app = Flask(__name__)
        
        # Create webhooks directory if it doesn't exist
        os.makedirs(webhooks_dir, exist_ok=True)
        
        # Register webhook endpoint
        self.app.add_url_rule(
            "/webhooks/<name>",
            "handle_webhook",
            self._handle_webhook,
            methods=["POST"]
        )

    def register_webhook(self, name: str, url: str, method: WebhookMethod,
                        headers: Dict[str, str], payload_template: str,
                        description: str) -> Webhook:
        """Register a new webhook.
        
        Args:
            name: Name of the webhook
            url: URL to send the webhook to
            method: HTTP method to use
            headers: Headers to include
            payload_template: Template for the payload
            description: Description of the webhook
            
        Returns:
            Created webhook
        """
        try:
            webhook = Webhook(
                name=name,
                url=url,
                method=method,
                headers=headers,
                payload_template=payload_template,
                description=description
            )
            
            self.webhooks[name] = webhook
            
            # Save webhook configuration
            self._save_webhook(webhook)
            
            logger.info(f"Registered webhook: {name}")
            
            return webhook
            
        except Exception as e:
            logger.error(f"Error registering webhook: {str(e)}")
            raise

    def _save_webhook(self, webhook: Webhook) -> None:
        """Save webhook configuration to disk."""
        try:
            webhook_file = os.path.join(self.webhooks_dir, f"{webhook.name}.json")
            
            webhook_data = {
                "name": webhook.name,
                "url": webhook.url,
                "method": webhook.method.value,
                "headers": webhook.headers,
                "payload_template": webhook.payload_template,
                "description": webhook.description
            }
            
            with open(webhook_file, "w") as f:
                json.dump(webhook_data, f, indent=2)
                
        except Exception as e:
            logger.error(f"Error saving webhook: {str(e)}")
            raise

    def trigger_webhook(self, name: str, data: Dict[str, Any]) -> Dict[str, Any]:
        """Trigger a webhook.
        
        Args:
            name: Name of the webhook to trigger
            data: Data to include in the payload
            
        Returns:
            Dictionary with webhook results
        """
        try:
            if name not in self.webhooks:
                raise ValueError(f"Webhook {name} not found")
                
            webhook = self.webhooks[name]
            
            # Format payload
            payload = self._format_payload(webhook.payload_template, data)
            
            # Send request
            response = requests.request(
                method=webhook.method.value,
                url=webhook.url,
                headers=webhook.headers,
                json=payload
            )
            
            return {
                "triggered": True,
                "webhook": name,
                "status_code": response.status_code,
                "response": response.json() if response.text else None
            }
            
        except Exception as e:
            logger.error(f"Error triggering webhook: {str(e)}")
            raise

    def _format_payload(self, template: str, data: Dict[str, Any]) -> Dict[str, Any]:
        """Format payload using template and data."""
        try:
            # In a real implementation, this would use a proper template engine
            # For now, just do simple string replacement
            payload_str = template
            
            for key, value in data.items():
                payload_str = payload_str.replace(f"${key}", str(value))
                
            return json.loads(payload_str)
            
        except Exception as e:
            logger.error(f"Error formatting payload: {str(e)}")
            raise

    def _handle_webhook(self, name: str) -> Response:
        """Handle incoming webhook requests."""
        try:
            if name not in self.webhooks:
                return Response(
                    json.dumps({"error": f"Webhook {name} not found"}),
                    status=404,
                    mimetype="application/json"
                )
                
            webhook = self.webhooks[name]
            
            # Get request data
            data = request.get_json()
            
            # Log webhook
            logger.info(f"Received webhook: {name}")
            logger.debug(f"Webhook data: {data}")
            
            # Return success
            return Response(
                json.dumps({"status": "ok"}),
                status=200,
                mimetype="application/json"
            )
            
        except Exception as e:
            logger.error(f"Error handling webhook: {str(e)}")
            return Response(
                json.dumps({"error": str(e)}),
                status=500,
                mimetype="application/json"
            )

    def start_server(self, host: str = "0.0.0.0", port: int = 5000) -> None:
        """Start the webhook server.
        
        Args:
            host: Host to bind to
            port: Port to listen on
        """
        try:
            self.app.run(host=host, port=port)
            
        except Exception as e:
            logger.error(f"Error starting webhook server: {str(e)}")
            raise

class ThirdPartyIntegration:
    """Base class for third-party integrations."""
    def __init__(self, name: str, config: Dict[str, Any]):
        self.name = name
        self.config = config
        self.webhook_manager = WebhookManager()

    def register_webhooks(self) -> None:
        """Register webhooks for the integration."""
        raise NotImplementedError

    def handle_event(self, event_type: str, data: Dict[str, Any]) -> None:
        """Handle an event from the third-party service.
        
        Args:
            event_type: Type of event
            data: Event data
        """
        raise NotImplementedError

class SlackIntegration(ThirdPartyIntegration):
    """Slack integration."""
    def __init__(self, config: Dict[str, Any]):
        super().__init__("slack", config)
        
    def register_webhooks(self) -> None:
        """Register Slack webhooks."""
        self.webhook_manager.register_webhook(
            name="slack_message",
            url=self.config["webhook_url"],
            method=WebhookMethod.POST,
            headers={"Content-Type": "application/json"},
            payload_template=json.dumps({
                "text": "${message}",
                "channel": "${channel}",
                "username": "Nessi Bot"
            }),
            description="Send messages to Slack"
        )
        
    def handle_event(self, event_type: str, data: Dict[str, Any]) -> None:
        """Handle Slack events."""
        if event_type == "message":
            self.webhook_manager.trigger_webhook("slack_message", {
                "message": data["text"],
                "channel": data["channel"]
            })

class GitHubIntegration(ThirdPartyIntegration):
    """GitHub integration."""
    def __init__(self, config: Dict[str, Any]):
        super().__init__("github", config)
        
    def register_webhooks(self) -> None:
        """Register GitHub webhooks."""
        self.webhook_manager.register_webhook(
            name="github_event",
            url=self.config["webhook_url"],
            method=WebhookMethod.POST,
            headers={
                "Content-Type": "application/json",
                "X-GitHub-Event": "${event_type}"
            },
            payload_template=json.dumps({
                "repository": "${repository}",
                "action": "${action}",
                "sender": "${sender}"
            }),
            description="Handle GitHub events"
        )
        
    def handle_event(self, event_type: str, data: Dict[str, Any]) -> None:
        """Handle GitHub events."""
        self.webhook_manager.trigger_webhook("github_event", {
            "event_type": event_type,
            "repository": data["repository"]["full_name"],
            "action": data.get("action", ""),
            "sender": data["sender"]["login"]
        }) 