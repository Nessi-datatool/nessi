#!/usr/bin/env python3
import argparse
import sys
from pathlib import Path

from src.security.ip_manager import IPManager
from src.config.security_config import SecurityConfig

def main():
    parser = argparse.ArgumentParser(description="Manage IP allowlist")
    subparsers = parser.add_subparsers(dest="command", help="Command to execute")
    
    # Add IP command
    add_parser = subparsers.add_parser("add", help="Add IP to allowlist")
    add_parser.add_argument("ip", help="IP address to add")
    
    # Remove IP command
    remove_parser = subparsers.add_parser("remove", help="Remove IP from allowlist")
    remove_parser.add_argument("ip", help="IP address to remove")
    
    # List command
    subparsers.add_parser("list", help="List allowed IPs")
    
    # Clear command
    subparsers.add_parser("clear", help="Clear allowlist")
    
    args = parser.parse_args()
    
    # Initialize IP manager
    config = SecurityConfig()
    ip_manager = IPManager(config.ip_allowlist_file)
    
    if args.command == "add":
        try:
            ip_manager.add_ip(args.ip)
            print(f"Added {args.ip} to allowlist")
        except ValueError as e:
            print(f"Error: {e}", file=sys.stderr)
            sys.exit(1)
    
    elif args.command == "remove":
        try:
            ip_manager.remove_ip(args.ip)
            print(f"Removed {args.ip} from allowlist")
        except ValueError as e:
            print(f"Error: {e}", file=sys.stderr)
            sys.exit(1)
    
    elif args.command == "list":
        allowed_ips = ip_manager.get_allowed_ips()
        if allowed_ips:
            print("Allowed IPs:")
            for ip in allowed_ips:
                print(f"- {ip}")
        else:
            print("No IPs in allowlist")
    
    elif args.command == "clear":
        ip_manager.clear_allowlist()
        print("Cleared allowlist")
    
    else:
        parser.print_help()
        sys.exit(1)

if __name__ == "__main__":
    main() 