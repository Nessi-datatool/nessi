"""
Script to add license headers to Python files.
"""

import os
import re

LICENSE_HEADER = '''"""
NESSI DATA TOOL - FREE FOR PERSONAL USE LICENSE

Copyright (c) 2025 Nessi Data Tool. All rights reserved.

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
   The Software is licensed, not sold. Nessi Data Tool retains all right, title, and interest in and to the Software, including all intellectual property rights.

4. DISCLAIMER OF WARRANTY
   THE SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING, BUT NOT LIMITED TO, THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND NONINFRINGEMENT.

5. LIMITATION OF LIABILITY
   IN NO EVENT SHALL NESSI DATA TOOL BE LIABLE FOR ANY CLAIM, DAMAGES, OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT, OR OTHERWISE, ARISING FROM, OUT OF, OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
"""

'''

def clean_content(content):
    """Remove all existing license headers from content."""
    # Remove any triple-quoted docstrings that contain NESSI DATA TOOL
    content = re.sub(r'"""[^"]*?NESSI DATA TOOL.*?"""', '', content, flags=re.DOTALL | re.MULTILINE)
    # Clean up any resulting multiple blank lines
    content = re.sub(r'\n\s*\n\s*\n', '\n\n', content)
    return content.strip()

def add_license_header(file_path):
    """Add license header to a file."""
    try:
        with open(file_path, 'r') as f:
            content = f.read()

        # Clean the content
        content = clean_content(content)

        # Handle shebang lines
        if content.startswith('#!'):
            shebang, rest = content.split('\n', 1)
            new_content = f"{shebang}\n{LICENSE_HEADER}{rest.strip()}"
        else:
            new_content = f"{LICENSE_HEADER}{content}"

        with open(file_path, 'w') as f:
            f.write(new_content)
        print(f"Updated license header in {file_path}")
    except Exception as e:
        print(f"Error processing {file_path}: {e}")

def process_directory(directory):
    """Process all Python files in a directory and its subdirectories."""
    for root, _, files in os.walk(directory):
        for file in files:
            if file.endswith('.py'):
                file_path = os.path.join(root, file)
                add_license_header(file_path)

if __name__ == "__main__":
    # Process backend directory
    process_directory("backend")
    # Process tests directory
    process_directory("tests")
    # Process src directory
    process_directory("src") 