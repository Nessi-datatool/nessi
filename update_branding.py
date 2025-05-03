"""
Script to update branding from Nessi to nessi.dev
"""

import os
import re

def update_file_content(file_path: str) -> None:
    """Update branding in a file."""
    with open(file_path, 'r') as f:
        content = f.read()
    
    # Replace Nessi with nessi.dev in various contexts
    content = re.sub(r'Copyright \(c\) 2025 Nessi', 'Copyright (c) 2025 nessi.dev', content)
    content = re.sub(r'Copyright © 2025 Nessi', 'Copyright © 2025 nessi.dev', content)
    content = re.sub(r'Copyright 2024 Nessi', 'Copyright 2024 nessi.dev', content)
    content = re.sub(r'Nessi is free for personal use', 'nessi.dev is free for personal use', content)
    content = re.sub(r'This version is free. Future versions of Nessi', 'This version is free. Future versions of nessi.dev', content)
    content = re.sub(r'The Software is licensed, not sold. Nessi', 'The Software is licensed, not sold. nessi.dev', content)
    content = re.sub(r'This Agreement shall be governed by and construed in accordance with the laws of the jurisdiction in which Nessi', 
                    'This Agreement shall be governed by and construed in accordance with the laws of the jurisdiction in which nessi.dev', content)
    content = re.sub(r'def get_spark_session\(app_name: str = "Nessi"\)', 'def get_spark_session(app_name: str = "nessi.dev")', content)
    content = re.sub(r'"""Nessi CLI - Free for personal use"""', '"""nessi.dev CLI - Free for personal use"""', content)
    content = re.sub(r'"""Nessi package\."""', '"""nessi.dev package."""', content)
    
    # Update GitHub references
    content = re.sub(r'https://github.com/Nessi-datatool/nessi', 'https://github.com/nessi-dev/nessi', content)
    content = re.sub(r'@Nessi-datatool/', '@nessi-dev/', content)
    
    with open(file_path, 'w') as f:
        f.write(content)
    print(f"Updated branding in {file_path}")

def main():
    """Main function to update branding in all files."""
    # Only process files in our project directory
    project_dirs = ['src', 'docs', '.github', 'tests']
    for root, _, files in os.walk('.'):
        # Skip virtual environment and other non-project directories
        if 'venv' in root or '__pycache__' in root or '.git' in root:
            continue
            
        # Only process files in project directories
        if not any(root.startswith(d) for d in project_dirs):
            continue
            
        for file in files:
            if file.endswith(('.py', '.md', '.html', '.toml', '.txt', '.yml', '.yaml')):
                file_path = os.path.join(root, file)
                try:
                    update_file_content(file_path)
                except Exception as e:
                    print(f"Error updating {file_path}: {str(e)}")

if __name__ == '__main__':
    main() 