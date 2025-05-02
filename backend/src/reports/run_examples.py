"""Script to run all example components together."""

import subprocess
import time
import sys
from pathlib import Path
from generate_examples import ExampleGenerator
from setup_grafana import GrafanaSetup
from test_setup import MonitoringTest

def run_docker_compose() -> bool:
    """Start Docker services."""
    try:
        print("Starting Docker services...")
        docker_compose_path = Path(__file__).parent.parent.parent / "docker-compose.test.yml"
        subprocess.run([
            "docker-compose", "-f", str(docker_compose_path), "up", "-d"
        ], check=True)
        print("Docker services started successfully")
        return True
    except subprocess.CalledProcessError as e:
        print(f"Error starting Docker services: {str(e)}")
        return False

def main():
    """Run all example components."""
    try:
        # Start Docker services
        if not run_docker_compose():
            sys.exit(1)
            
        # Wait for services to be ready
        print("\nWaiting for services to be ready...")
        time.sleep(30)  # Give services time to start
        
        # Run Grafana setup
        print("\nSetting up Grafana...")
        grafana_setup = GrafanaSetup()
        if not grafana_setup.run():
            print("Grafana setup failed")
            sys.exit(1)
            
        # Start example generator
        print("\nStarting example generator...")
        generator = ExampleGenerator()
        
        # Run in background
        import threading
        generator_thread = threading.Thread(target=generator.run)
        generator_thread.daemon = True
        generator_thread.start()
        
        # Run verification tests
        print("\nRunning verification tests...")
        test = MonitoringTest()
        results = test.run_tests()
        
        # Print results
        print("\nTest Results:")
        for test_name, passed in results.items():
            print(f"{test_name}: {'✓' if passed else '✗'}")
            
        if all(results.values()):
            print("\nAll components are working correctly!")
            print("\nYou can now access:")
            print("1. Grafana Dashboard: http://localhost:3000/d/nessi-reports")
            print("2. Prometheus Metrics: http://localhost:9090")
            print("3. Example Reports: Check the 'reports' directory")
            print("\nPress Ctrl+C to stop the example generator...")
            
            # Keep the script running
            try:
                while True:
                    time.sleep(1)
            except KeyboardInterrupt:
                print("\nStopping example generator...")
                
        else:
            print("\nSome components failed. Please check the logs above.")
            sys.exit(1)
            
    except Exception as e:
        print(f"Error running examples: {str(e)}")
        sys.exit(1)
        
    finally:
        # Clean up Docker services
        print("\nCleaning up Docker services...")
        docker_compose_path = Path(__file__).parent.parent.parent / "docker-compose.test.yml"
        subprocess.run([
            "docker-compose", "-f", str(docker_compose_path), "down"
        ])
        
if __name__ == "__main__":
    main() 