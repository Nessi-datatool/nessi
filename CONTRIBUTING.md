# Contributing to Nessi

Thank you for your interest in contributing to Nessi! This document provides guidelines and instructions for contributing to the project.

## Code of Conduct

By participating in this project, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md).

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/your-username/nessi.git`
3. Create a new branch: `git checkout -b feature/your-feature-name`
4. Set up your development environment:
   ```bash
   python -m venv venv
   source venv/bin/activate  # On Windows: venv\Scripts\activate
   pip install -r requirements.txt
   pip install -e .
   ```

## Development Workflow

1. Make your changes
2. Run tests: `python -m pytest tests/ -v`
3. Ensure code style consistency
4. Update documentation if needed
5. Commit your changes with a descriptive message
6. Push to your fork
7. Create a Pull Request

## Testing

- All new features should include tests
- Run the full test suite before submitting a PR
- Ensure test coverage remains high
- For test coverage report: `python -m pytest tests/ --cov=src --cov-report=html`

## Code Style

- Follow PEP 8 guidelines
- Use type hints for better code documentation
- Keep functions focused and small
- Add docstrings for all public functions and classes
- Use meaningful variable and function names

## Documentation

- Update README.md for significant changes
- Add or update API.md for API changes
- Include docstrings for all new functions and classes
- Update examples if they're affected by your changes

## Pull Request Process

1. Update the README.md with details of changes if needed
2. Update the CHANGELOG.md with a summary of changes
3. The PR must pass all CI checks
4. At least one maintainer must approve the PR
5. The PR must be up to date with the main branch

## Reporting Issues

- Use the issue tracker to report bugs or request features
- Include a clear description of the issue
- Provide steps to reproduce if it's a bug
- Include relevant error messages and stack traces
- Specify your environment (OS, Python version, etc.)

## Questions?

Feel free to open an issue if you have any questions about contributing to Nessi. 