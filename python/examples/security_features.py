#!/usr/bin/env python3
"""
Security features example for the Nessi monitoring client.

This example demonstrates how to use the security features of the Nessi monitoring system,
including authentication, API keys, and role-based access control.
"""

import sys
import os
import time
from datetime import datetime
import json

# Add the parent directory to the path so we can import the client
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
from nessi_client import NessiClient

def main():
    """Main function demonstrating security features of the Nessi client."""
    
    print("Nessi Monitoring Client - Security Features Example")
    print("=================================================")
    
    # Step 1: Authenticate with username and password
    print("\n1. Authenticating with username and password")
    print("------------------------------------------")
    
    try:
        # Initialize the client with username and password
        client = NessiClient(
            base_url="http://localhost:8080",
            username="admin",  # Default admin user
            password="admin"   # Default admin password
        )
        print("Successfully authenticated with username and password")
    except Exception as e:
        print(f"Authentication failed: {e}")
        return
    
    # Step 2: Create a new user with limited permissions
    print("\n2. Creating a new user with viewer role")
    print("-------------------------------------")
    
    new_user = {
        "username": "viewer_user",
        "password": "password123",
        "role": "viewer",  # Limited permissions
        "email": "viewer@example.com"
    }
    
    try:
        # Admin users can create new users
        response = client._request("post", "api/users", json=new_user)
        print(f"Created new user: {new_user['username']} with role: {new_user['role']}")
    except Exception as e:
        print(f"Failed to create user: {e}")
    
    # Step 3: Generate an API key for the admin user
    print("\n3. Generating an API key")
    print("----------------------")
    
    try:
        # Generate an API key
        response = client._request("post", "api/apikeys", json={"description": "Example API key"})
        api_key_data = response.json()
        api_key = api_key_data.get("key")
        
        if api_key:
            print(f"Generated API key: {api_key[:8]}...{api_key[-4:]}")
            print(f"Key description: {api_key_data.get('description')}")
            print(f"Created at: {api_key_data.get('created_at')}")
        else:
            print("Failed to generate API key")
            return
    except Exception as e:
        print(f"Failed to generate API key: {e}")
        return
    
    # Step 4: Authenticate with the API key
    print("\n4. Authenticating with API key")
    print("----------------------------")
    
    try:
        # Initialize a new client with the API key
        api_client = NessiClient(
            base_url="http://localhost:8080",
            api_key=api_key
        )
        print("Successfully authenticated with API key")
    except Exception as e:
        print(f"API key authentication failed: {e}")
        return
    
    # Step 5: Test role-based access control
    print("\n5. Testing role-based access control")
    print("----------------------------------")
    
    # Authenticate as the viewer user
    try:
        viewer_client = NessiClient(
            base_url="http://localhost:8080",
            username="viewer_user",
            password="password123"
        )
        print("Successfully authenticated as viewer user")
    except Exception as e:
        print(f"Viewer authentication failed: {e}")
        return
    
    # Test admin-only operations with the viewer user
    print("\n   Testing admin-only operations with viewer user:")
    try:
        # Attempt to create another user (should fail for viewers)
        another_user = {
            "username": "another_user",
            "password": "password456",
            "role": "user",
            "email": "another@example.com"
        }
        response = viewer_client._request("post", "api/users", json=another_user)
        print("  ❌ Unexpected success: Viewer was able to create a user")
    except Exception as e:
        print("  ✅ Expected failure: Viewer cannot create users")
    
    # Test allowed operations for viewers
    print("\n   Testing allowed operations for viewer user:")
    try:
        # Viewers should be able to get metrics
        response = viewer_client._request("get", "api/metrics", params={"limit": 5})
        metrics = response.json()
        print(f"  ✅ Success: Viewer can access metrics ({len(metrics)} retrieved)")
    except Exception as e:
        print(f"  ❌ Unexpected failure: Viewer cannot access metrics: {e}")
    
    # Step 6: Test SSL/HTTPS support
    print("\n6. Testing SSL/HTTPS support")
    print("-------------------------")
    
    try:
        # Initialize a client with SSL verification
        https_client = NessiClient(
            base_url="https://localhost:8443",  # HTTPS port
            username="admin",
            password="admin",
            verify_ssl=True  # Verify SSL certificates
        )
        print("Successfully connected to HTTPS endpoint with SSL verification")
    except Exception as e:
        print(f"HTTPS connection failed: {e}")
        print("Note: This is expected if you're using self-signed certificates or HTTPS is not enabled")
        
        # Try again without SSL verification
        try:
            https_client = NessiClient(
                base_url="https://localhost:8443",  # HTTPS port
                username="admin",
                password="admin",
                verify_ssl=False  # Skip SSL verification
            )
            print("Successfully connected to HTTPS endpoint without SSL verification")
        except Exception as e:
            print(f"HTTPS connection still failed: {e}")
    
    # Step 7: List and revoke API keys
    print("\n7. Managing API keys")
    print("------------------")
    
    try:
        # List API keys
        response = client._request("get", "api/apikeys")
        api_keys = response.json()
        
        print(f"Found {len(api_keys)} API keys:")
        for key in api_keys:
            print(f"  - ID: {key.get('id')}")
            print(f"    Description: {key.get('description')}")
            print(f"    Created: {key.get('created_at')}")
            print(f"    Last used: {key.get('last_used_at', 'Never')}")
            print()
        
        # Revoke the API key we created
        if api_keys and len(api_keys) > 0:
            key_id = api_keys[0].get('id')
            response = client._request("delete", f"api/apikeys/{key_id}")
            print(f"Revoked API key with ID: {key_id}")
    except Exception as e:
        print(f"Failed to manage API keys: {e}")
    
    # Step 8: Test token expiration and refresh
    print("\n8. Testing token expiration and refresh")
    print("------------------------------------")
    
    print("Note: This would normally require waiting for token expiration")
    print("For demonstration purposes, we'll simulate a token refresh")
    
    try:
        # Force a token refresh
        client._refresh_token_if_needed()
        print("Token refreshed successfully")
    except Exception as e:
        print(f"Token refresh failed: {e}")
    
    print("\nSecurity features example completed!")

if __name__ == "__main__":
    main()
