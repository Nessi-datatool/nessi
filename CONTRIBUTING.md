# Contributing to Nessi

Thank you for your interest in contributing to Nessi! This document provides guidelines and instructions for contributing to the project.

## Code of Conduct

By participating in this project, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md).

## Developer Certificate of Origin (DCO)

All contributors must sign off on their commits. This certifies that you wrote or have the right to submit the code you are contributing to the project. To sign off on a commit, add the `-s` flag:

```bash
git commit -s -m "Your commit message"
```

Or manually add a line like this to your commit message:
```
Signed-off-by: Your Name <your.email@example.com>
```

## Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/your-username/nessi.git
   cd nessi
   ```
3. Create a new branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```
4. Set up your development environment:
   ```bash
   python -m venv venv
   source venv/bin/activate  # On Windows: venv\Scripts\activate
   pip install -r requirements.txt
   pip install -e .
   ```

## Development Workflow

1. Make your changes
2. Run tests:
   ```bash
   python -m pytest tests/ -v
   ```
3. Ensure code style:
   ```bash
   flake8 .
   ```
4. Sign off on your commits:
   ```bash
   git commit -s -m "Your commit message"
   ```
5. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```
6. Create a pull request

## Pull Request Process

1. Update the README.md with details of changes if needed
2. Update the CHANGELOG.md
3. Ensure all CI checks pass
4. Get maintainer approval
5. Merge the pull request

## Testing

- Include tests for new features
- Ensure all tests pass
- Maintain or improve test coverage
- Test with multiple Python versions

## Code Style

- Follow PEP 8 guidelines
- Use type hints
- Write docstrings for all public functions
- Keep functions small and focused
- Use meaningful variable names

## Documentation

- Update documentation for significant changes
- Document new features
- Update API documentation if needed
- Keep examples up to date

## Security

- Follow security best practices
- Report vulnerabilities responsibly
- Keep dependencies up to date
- Review code for security issues

## Questions?

If you have questions about contributing, please open an issue. 