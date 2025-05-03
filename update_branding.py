"""
Script to update branding from nessi.dev to nessi.dev
"""

import os
import re

def update_gitignore(file_path: str) -> None:
    """Update .gitignore file."""
    with open(file_path, 'r') as f:
        content = f.read()
    content = content.replace('# .gitignore for nessi.dev project', '# .gitignore for nessi.dev project')
    with open(file_path, 'w') as f:
        f.write(content)
    print(f"Updated branding in {file_path}")

def update_codeowners(file_path: str) -> None:
    """Update CODEOWNERS file."""
    with open(file_path, 'r') as f:
        content = f.read()
    content = content.replace('@nessi-dev/', '@nessi-dev/')
    with open(file_path, 'w') as f:
        f.write(content)
    print(f"Updated branding in {file_path}")

def update_bug_report(file_path: str) -> None:
    """Update bug report template."""
    with open(file_path, 'r') as f:
        content = f.read()
    content = content.replace('nessi.dev Version:', 'nessi.dev Version:')
    with open(file_path, 'w') as f:
        f.write(content)
    print(f"Updated branding in {file_path}")

def update_gitleaks(file_path: str) -> None:
    """Update gitleaks configuration."""
    with open(file_path, 'r') as f:
        content = f.read()
    content = content.replace('title = "nessi.dev Gitleaks Configuration"', 'title = "nessi.dev Gitleaks Configuration"')
    with open(file_path, 'w') as f:
        f.write(content)
    print(f"Updated branding in {file_path}")

def update_notice(file_path: str) -> None:
    """Update NOTICE file."""
    with open(file_path, 'r') as f:
        content = f.read()
    content = content.replace('Nessi\n', 'nessi.dev\n')
    content = content.replace('Copyright 2024 nessi.dev', 'Copyright 2024 nessi.dev')
    with open(file_path, 'w') as f:
        f.write(content)
    print(f"Updated branding in {file_path}")

def update_license(file_path: str) -> None:
    """Update LICENSE file."""
    with open(file_path, 'r') as f:
        content = f.read()
    replacements = [
        ('Copyright (c) 2025 nessi.dev', 'Copyright (c) 2025 nessi.dev'),
        ('nessi.dev is free for personal use', 'nessi.dev is free for personal use'),
        ('This version is free. Future versions of nessi.dev', 'This version is free. Future versions of nessi.dev'),
        ('The Software is licensed, not sold. nessi.dev', 'The Software is licensed, not sold. nessi.dev'),
        ('This Agreement shall be governed by and construed in accordance with the laws of the jurisdiction in which nessi.dev',
         'This Agreement shall be governed by and construed in accordance with the laws of the jurisdiction in which nessi.dev')
    ]
    for old, new in replacements:
        content = content.replace(old, new)
    with open(file_path, 'w') as f:
        f.write(content)
    print(f"Updated branding in {file_path}")

def update_file_content(file_path: str) -> None:
    """Update branding in a file."""
    if file_path.endswith('.gitignore'):
        update_gitignore(file_path)
    elif file_path.endswith('CODEOWNERS'):
        update_codeowners(file_path)
    elif file_path.endswith('bug_report.md'):
        update_bug_report(file_path)
    elif file_path.endswith('gitleaks.toml'):
        update_gitleaks(file_path)
    elif file_path.endswith('NOTICE'):
        update_notice(file_path)
    elif file_path.endswith('LICENSE') or file_path.endswith('LICENSE.md'):
        update_license(file_path)
    else:
        with open(file_path, 'r') as f:
            content = f.read()
        
        # Replace Nessi with nessi.dev in various contexts
        replacements = [
            (r'Copyright \(c\) 2025 Nessi', 'Copyright (c) 2025 nessi.dev'),
            (r'Copyright © 2025 nessi.dev', 'Copyright © 2025 nessi.dev'),
            (r'Copyright 2024 nessi.dev', 'Copyright 2024 nessi.dev'),
            (r'nessi.dev is free for personal use', 'nessi.dev is free for personal use'),
            (r'This version is free. Future versions of nessi.dev', 'This version is free. Future versions of nessi.dev'),
            (r'The Software is licensed, not sold. nessi.dev', 'The Software is licensed, not sold. nessi.dev'),
            (r'This Agreement shall be governed by and construed in accordance with the laws of the jurisdiction in which nessi.dev', 
             'This Agreement shall be governed by and construed in accordance with the laws of the jurisdiction in which nessi.dev'),
            (r'def get_spark_session\(app_name: str = "Nessi"\)', 'def get_spark_session(app_name: str = "nessi.dev")'),
            (r'def get_spark_session\(app_name="NessiApp"\)', 'def get_spark_session(app_name="nessi.dev")'),
            (r'def get_session\(self, app_name="NessiApp"\)', 'def get_session(self, app_name="nessi.dev")'),
            (r'"""nessi.dev CLI - Free for personal use"""', '"""nessi.dev CLI - Free for personal use"""'),
            (r'"""nessi.dev package\."""', '"""nessi.dev package."""'),
            (r'https://github.com/nessi-dev/nessi', 'https://github.com/nessi-dev/nessi'),
            (r'@nessi-dev/', '@nessi-dev/'),
            (r'team@nessi.dev', 'team@nessi.dev'),
            (r'Generated by nessi.dev -', 'Generated by nessi.dev -'),
            (r'# Nessi\n', '# nessi.dev\n'),
            (r'# Contributing to Nessi\n', '# Contributing to nessi.dev\n'),
            (r'nessi.dev Version:', 'nessi.dev Version:'),
            (r'title = "nessi.dev ', 'title = "nessi.dev '),
            (r'appName\("NessiTest"\)', 'appName("nessi.dev-test")'),
            (r'appName\("NessiDemo"\)', 'appName("nessi.dev-demo")'),
            (r'# .gitignore for nessi.dev', '# .gitignore for nessi.dev'),
            (r'{ name = "nessi.dev Team"', '{ name = "nessi.dev Team"'),
            (r'# nessi.dev Report Generation', '# nessi.dev Report Generation'),
            (r'# nessi.dev Monitoring Guide', '# nessi.dev Monitoring Guide'),
            (r'# nessi.dev Demo Tutorial', '# nessi.dev Demo Tutorial'),
            (r'# nessi.dev Installation Guide', '# nessi.dev Installation Guide'),
            (r'# nessi.dev API Documentation', '# nessi.dev API Documentation'),
            (r'# nessi.dev Features', '# nessi.dev Features'),
            (r'# nessi.dev Quickstart Guide', '# nessi.dev Quickstart Guide'),
            (r'# nessi.dev Monitoring System', '# nessi.dev Monitoring System'),
            (r'nessi.dev provides', 'nessi.dev provides'),
            (r'nessi.dev can generate', 'nessi.dev can generate'),
            (r'nessi.dev supports', 'nessi.dev supports'),
            (r'nessi.dev is a powerful', 'nessi.dev is a powerful'),
            (r'nessi.dev is a Python-based', 'nessi.dev is a Python-based'),
            (r'Install nessi.dev', 'Install nessi.dev'),
            (r'Start the nessi.dev CLI', 'Start the nessi.dev CLI'),
            (r'"""[^"]*Nessi[^"]*"""', lambda m: m.group(0).replace('Nessi', 'nessi.dev')),
            (r'"title": "nessi.dev', '"title": "nessi.dev'),
            (r'The nessi.dev', 'The nessi.dev'),
            (r'the nessi.dev', 'the nessi.dev'),
            (r'Nessi\'s', 'nessi.dev\'s'),
            (r'nessi.dev Reports Dashboard', 'nessi.dev Reports Dashboard'),
            (r'nessi.dev Demo Dashboard', 'nessi.dev Demo Dashboard'),
            (r'nessi.dev Table Metrics', 'nessi.dev Table Metrics'),
            (r'nessi.dev Metrics Job', 'nessi.dev Metrics Job'),
            (r'nessi.dev Gitleaks Configuration', 'nessi.dev Gitleaks Configuration'),
            (r'Nessi\. All Rights Reserved', 'nessi.dev. All Rights Reserved'),
            (r'nessi.dev software', 'nessi.dev software'),
            (r'nessi.dev system', 'nessi.dev system'),
            (r'nessi.dev monitoring system', 'nessi.dev monitoring system'),
            (r'nessi.dev Monitoring System', 'nessi.dev Monitoring System'),
            (r'nessi.dev generates', 'nessi.dev generates'),
            (r'Nessi$', 'nessi.dev'),
            (r'running nessi.dev', 'running nessi.dev'),
            (r'get started with nessi.dev', 'get started with nessi.dev'),
            (r'system for nessi.dev', 'system for nessi.dev'),
            (r'│ nessi.dev │', '│ nessi.dev │'),
            (r'get_spark_session\("NessiDemo"\)', 'get_spark_session("nessi.dev-demo")'),
            (r'get_spark_session\("NessiTest"\)', 'get_spark_session("nessi.dev-test")'),
            (r'Before installing nessi.dev', 'Before installing nessi.dev'),
            (r'- nessi.dev backend', '- nessi.dev backend'),
            (r'nessi.dev project', 'nessi.dev project'),
            (r'@nessi-dev/owners', '@nessi-dev/owners'),
            (r'@nessi-dev/python-reviewers', '@nessi-dev/python-reviewers'),
            (r'@nessi-dev/documentation-reviewers', '@nessi-dev/documentation-reviewers'),
            (r'@nessi-dev/test-reviewers', '@nessi-dev/test-reviewers'),
        ]
        
        for pattern, replacement in replacements:
            if callable(replacement):
                content = re.sub(pattern, replacement, content)
            else:
                content = re.sub(pattern, replacement, content)
        
        with open(file_path, 'w') as f:
            f.write(content)
        print(f"Updated branding in {file_path}")

def main():
    """Main function to update branding in all files."""
    # Only process files in our project directory
    project_dirs = ['src', 'docs', '.github', 'tests', 'backend']
    for root, _, files in os.walk('.'):
        # Skip virtual environment and other non-project directories
        if 'venv' in root or '__pycache__' in root or '.git' in root:
            continue
            
        # Only process files in project directories
        if not any(d in root for d in project_dirs) and not os.path.basename(root) in project_dirs and root != '.':
            continue
            
        for file in files:
            if file.endswith(('.py', '.md', '.html', '.toml', '.txt', '.yml', '.yaml', '.json', 'CODEOWNERS', '.gitignore', 'NOTICE', 'LICENSE')):
                file_path = os.path.join(root, file)
                try:
                    update_file_content(file_path)
                except Exception as e:
                    print(f"Error updating {file_path}: {str(e)}")

if __name__ == '__main__':
    main() 