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

import os
import json
import tempfile
import pytest
from src.reports.customize_alerts import AlertCustomizer

@pytest.fixture
def sample_dashboard():
    """Create a sample dashboard with alerts for testing."""
    dashboard = {
        "panels": [
            {
                "id": 1,
                "title": "CPU Usage",
                "targets": [{"expr": "cpu_usage"}],
                "alert": {
                    "conditions": [
                        {
                            "evaluator": {
                                "params": [80.0],
                                "type": "gt"
                            }
                        }
                    ],
                    "notifications": ["email"]
                }
            },
            {
                "id": 2,
                "title": "Memory Usage",
                "targets": [{"expr": "memory_usage"}],
                "alert": {
                    "conditions": [
                        {
                            "evaluator": {
                                "params": [90.0],
                                "type": "gt"
                            }
                        }
                    ],
                    "notifications": ["slack"]
                }
            }
        ]
    }
    return dashboard

@pytest.fixture
def dashboard_file(sample_dashboard):
    """Create a temporary dashboard file for testing."""
    with tempfile.NamedTemporaryFile(mode='w', suffix='.json', delete=False) as f:
        json.dump(sample_dashboard, f)
        f.flush()
        yield f.name
    os.unlink(f.name)

def test_get_available_alerts(dashboard_file):
    """Test getting available alerts from dashboard."""
    customizer = AlertCustomizer(dashboard_file)
    alerts = customizer.get_available_alerts()
    
    assert len(alerts) == 2
    assert alerts[0]['panel_id'] == 1
    assert alerts[0]['title'] == "CPU Usage"
    assert alerts[0]['metric'] == "cpu_usage"
    assert alerts[0]['current_threshold'] == 80.0
    
    assert alerts[1]['panel_id'] == 2
    assert alerts[1]['title'] == "Memory Usage"
    assert alerts[1]['metric'] == "memory_usage"
    assert alerts[1]['current_threshold'] == 90.0

def test_update_alert(dashboard_file):
    """Test updating alert threshold and notification channels."""
    customizer = AlertCustomizer(dashboard_file)
    
    # Update CPU alert
    customizer.update_alert(1, 85.0, ["email", "slack"])
    customizer.save_dashboard()
    
    # Verify changes
    with open(dashboard_file, 'r') as f:
        updated_dashboard = json.load(f)
    
    cpu_panel = next(p for p in updated_dashboard['panels'] if p['id'] == 1)
    assert cpu_panel['alert']['conditions'][0]['evaluator']['params'][0] == 85.0
    assert cpu_panel['alert']['notifications'] == ["email", "slack"]
    
    # Verify other panel unchanged
    memory_panel = next(p for p in updated_dashboard['panels'] if p['id'] == 2)
    assert memory_panel['alert']['conditions'][0]['evaluator']['params'][0] == 90.0
    assert memory_panel['alert']['notifications'] == ["slack"]

def test_update_alert_no_channels(dashboard_file):
    """Test updating alert threshold without changing notification channels."""
    customizer = AlertCustomizer(dashboard_file)
    
    # Update CPU alert threshold only
    customizer.update_alert(1, 75.0)
    customizer.save_dashboard()
    
    # Verify changes
    with open(dashboard_file, 'r') as f:
        updated_dashboard = json.load(f)
    
    cpu_panel = next(p for p in updated_dashboard['panels'] if p['id'] == 1)
    assert cpu_panel['alert']['conditions'][0]['evaluator']['params'][0] == 75.0
    assert cpu_panel['alert']['notifications'] == ["email"]  # Should remain unchanged

def test_update_nonexistent_panel(dashboard_file):
    """Test updating alert for non-existent panel."""
    customizer = AlertCustomizer(dashboard_file)
    
    # Should not raise an exception
    customizer.update_alert(999, 50.0)
    customizer.save_dashboard()
    
    # Verify dashboard unchanged
    with open(dashboard_file, 'r') as f:
        updated_dashboard = json.load(f)
    
    assert len(updated_dashboard['panels']) == 2
    assert all(p['id'] in [1, 2] for p in updated_dashboard['panels'])