"""
NESSI - FREE FOR PERSONAL USE LICENSE

Copyright (c) 2025 Nessi. All rights reserved.

Nessi is free for personal use.

A paid license is required for enterprise or consulting use.

This version is free. Future versions of Nessi will require a license for all users.

1. PERSONAL USE LICENSE
   a) The Software is provided free of charge for personal use
   b) Personal use includes:
      - Individual data analysis and processing
      - Educational purposes
      - Non-commercial research
   c) Enterprise or consulting use requires a paid license

2. RESTRICTIONS
   You shall not:
   a) Copy, modify, adapt, translate, reverse engineer, decompile, or disassemble the Software
   b) Create derivative works based on the Software
   c) Rent, lease, loan, sell, sublicense, distribute, transmit, or otherwise transfer the Software
   d) Remove or alter any proprietary notices or labels on the Software
   e) Use the Software for enterprise or consulting purposes without a valid paid license

3. OWNERSHIP
   The Software is licensed, not sold. Nessi retains all right, title, and interest in and to the Software, including all intellectual property rights.

4. DISCLAIMER OF WARRANTY
   THE SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING, BUT NOT LIMITED TO, THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND NONINFRINGEMENT.

5. LIMITATION OF LIABILITY
   IN NO EVENT SHALL NESSI BE LIABLE FOR ANY CLAIM, DAMAGES, OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT, OR OTHERWISE, ARISING FROM, OUT OF, OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
"""

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