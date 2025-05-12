#!/usr/bin/env python3
"""
Basic usage example for the Nessi monitoring client.
"""

import sys
import os
import time
from datetime import datetime, timedelta
import random

# Add the parent directory to the path so we can import the client
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
from nessi_client import NessiClient, Metric, Rule

def main():
    """Main function demonstrating basic usage of the Nessi client."""
    
    print("Nessi Monitoring Client - Basic Usage Example")
    print("=============================================")
    
    # Initialize the client
    # In a real application, you would use your actual server URL and credentials
    client = NessiClient(
        base_url="http://localhost:8080",
        username="admin",  # Default admin user
        password="admin"   # Default admin password
    )
    
    print("\n1. Sending metrics")
    print("-----------------")
    
    # Send some CPU usage metrics
    for i in range(5):
        # Create a simulated CPU usage value between 10% and 90%
        cpu_value = random.uniform(10, 90)
        
        metric = Metric(
            name="cpu_usage",
            value=cpu_value,
            timestamp=datetime.now() - timedelta(minutes=i*5),  # Backdate metrics
            tags={
                "host": "server-01",
                "environment": "production",
                "region": "us-west"
            }
        )
        
        try:
            result = client.send_metric(metric)
            print(f"Sent metric: cpu_usage = {cpu_value:.2f}%")
        except Exception as e:
            print(f"Error sending metric: {e}")
    
    # Send some memory usage metrics
    for i in range(5):
        # Create a simulated memory usage value between 20% and 80%
        memory_value = random.uniform(20, 80)
        
        metric = Metric(
            name="memory_usage",
            value=memory_value,
            timestamp=datetime.now() - timedelta(minutes=i*5),  # Backdate metrics
            tags={
                "host": "server-01",
                "environment": "production",
                "region": "us-west"
            }
        )
        
        try:
            result = client.send_metric(metric)
            print(f"Sent metric: memory_usage = {memory_value:.2f}%")
        except Exception as e:
            print(f"Error sending metric: {e}")
    
    # Let's wait a moment to ensure metrics are processed
    time.sleep(1)
    
    print("\n2. Retrieving metrics")
    print("-------------------")
    
    try:
        # Get CPU usage metrics
        cpu_metrics = client.get_metrics(
            name="cpu_usage",
            limit=10
        )
        
        print(f"Retrieved {len(cpu_metrics)} CPU usage metrics:")
        for metric in cpu_metrics:
            print(f"  {metric.timestamp}: {metric.name} = {metric.value:.2f}%")
        
        # Get memory usage metrics
        memory_metrics = client.get_metrics(
            name="memory_usage",
            limit=10
        )
        
        print(f"\nRetrieved {len(memory_metrics)} memory usage metrics:")
        for metric in memory_metrics:
            print(f"  {metric.timestamp}: {metric.name} = {metric.value:.2f}%")
    
    except Exception as e:
        print(f"Error retrieving metrics: {e}")
    
    print("\n3. Creating data quality rules")
    print("----------------------------")
    
    # Create a rule for CPU usage
    cpu_rule = Rule(
        id="cpu_high",
        name="High CPU Usage",
        description="Alerts when CPU usage exceeds 80%",
        severity="warning",
        field="value",
        rule_type="range",
        config={
            "min": 0,
            "max": 80
        },
        tags=["cpu", "performance"]
    )
    
    try:
        result = client.create_rule(cpu_rule)
        print(f"Created rule: {cpu_rule.name}")
    except Exception as e:
        print(f"Error creating rule: {e}")
    
    # Create a rule for memory usage
    memory_rule = Rule(
        id="memory_high",
        name="High Memory Usage",
        description="Alerts when memory usage exceeds 75%",
        severity="warning",
        field="value",
        rule_type="range",
        config={
            "min": 0,
            "max": 75
        },
        tags=["memory", "performance"]
    )
    
    try:
        result = client.create_rule(memory_rule)
        print(f"Created rule: {memory_rule.name}")
    except Exception as e:
        print(f"Error creating rule: {e}")
    
    print("\n4. Validating data")
    print("----------------")
    
    # Create some test data
    test_data = [
        {"metric": "cpu_usage", "value": 85, "host": "server-01"},  # Should trigger cpu_high rule
        {"metric": "memory_usage", "value": 60, "host": "server-01"},  # Should not trigger memory_high rule
        {"metric": "memory_usage", "value": 90, "host": "server-02"}   # Should trigger memory_high rule
    ]
    
    try:
        # Validate the data against our rules
        results = client.validate_data(test_data)
        
        print(f"Validation results ({len(results)} total):")
        for result in results:
            status = "PASSED" if result.passed else "FAILED"
            print(f"  {status}: {result.rule_id} - {result.message}")
    
    except Exception as e:
        print(f"Error validating data: {e}")
    
    print("\n5. Retrieving rules")
    print("----------------")
    
    try:
        # Get all rules
        rules = client.get_rules()
        
        print(f"Retrieved {len(rules)} rules:")
        for rule in rules:
            print(f"  {rule.id}: {rule.name} ({rule.severity})")
            print(f"    Description: {rule.description}")
            print(f"    Type: {rule.rule_type}")
            print(f"    Config: {rule.config}")
            print()
    
    except Exception as e:
        print(f"Error retrieving rules: {e}")
    
    print("\nExample completed successfully!")

if __name__ == "__main__":
    main()
