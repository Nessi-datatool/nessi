# Setting Up Branch Protection Rules for Nessi

This document provides detailed instructions for setting up branch protection rules in GitHub to ensure code quality and prevent merging PRs that fail CI checks.

## Why Branch Protection Rules?

Branch protection rules help maintain code quality and security by:

- Preventing direct pushes to important branches (like `main`)
- Requiring peer reviews before merging code
- Ensuring CI checks pass before code can be merged
- Maintaining a clean and reliable Git history

## Setting Up Branch Protection in GitHub

### 1. Access Repository Settings

1. Go to your GitHub repository
2. Click on "Settings" in the top navigation bar
3. In the left sidebar, click on "Branches"

### 2. Create Branch Protection Rule

1. Under "Branch protection rules", click "Add rule"
2. In the "Branch name pattern" field, enter `main`

### 3. Configure Protection Settings

Enable the following settings:

#### Required Checks and Reviews

- ✅ **Require a pull request before merging**
  - ✅ Require approvals (set to at least 1)
  - ✅ Dismiss stale pull request approvals when new commits are pushed

- ✅ **Require status checks to pass before merging**
  - ✅ Require branches to be up to date before merging
  - In the "Status checks that are required" section, add:
    - `Test` (from CI Pipeline workflow)
    - `Security Scan` (from CI Pipeline workflow)
    - `Databricks Integration Test` (from CI Pipeline workflow)
    - `Format and Lint` (from format-lint workflow)
    - `PR Checks` (from pr-checks workflow)
    - `Branch Protection Check` (from branch-protection workflow)

#### Additional Protections

- ✅ **Require conversation resolution before merging**
- ✅ **Require linear history**
- ✅ **Do not allow bypassing the above settings**

### 4. Save the Rule

Click "Create" or "Save changes" at the bottom of the page.

## Local Development Workflow

With branch protection rules in place, follow this workflow for development:

1. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes and commit them**:
   ```bash
   git add .
   git commit -m "feat: your descriptive commit message"
   ```

3. **Push your branch to GitHub**:
   ```bash
   git push -u origin feature/your-feature-name
   ```

4. **Create a Pull Request** in GitHub:
   - Add a descriptive title following semantic conventions (feat, fix, docs, etc.)
   - Write a detailed description
   - Link related issues using keywords like "Fixes #123"

5. **Address review feedback** by making additional commits to your branch

6. **Ensure all CI checks pass**

7. **Merge the PR** once approved and all checks pass

## Pre-Commit Hooks

We've set up local pre-commit hooks to catch issues before they reach CI:

- Formats Go files using `go fmt`
- Formats the `go.mod` file
- Runs linting if golangci-lint is installed

These hooks are automatically configured when you clone the repository.

## CI Pipeline Overview

Our CI pipeline includes:

1. **Format and Lint**: Ensures code is properly formatted and passes linting
2. **Tests**: Runs all unit and integration tests
3. **Security Scan**: Checks for security vulnerabilities
4. **Databricks Integration Tests**: Ensures Databricks functionality works correctly
5. **PR Checks**: Validates PR title format and description

## Troubleshooting

If your PR can't be merged due to failing checks:

1. Click on the failing check in GitHub to see detailed error logs
2. Fix the issues locally
3. Push the fixes to your branch
4. Wait for the checks to run again

## Best Practices

1. **Keep PRs focused and small** for easier review
2. **Write descriptive commit messages** following semantic conventions
3. **Link PRs to issues** for better traceability
4. **Run tests locally** before pushing
5. **Address all review comments** before requesting re-review

## Additional Resources

- [GitHub Documentation on Branch Protection](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/defining-the-mergeability-of-pull-requests/about-protected-branches)
- [Semantic Commit Messages](https://www.conventionalcommits.org/)
- [Effective Code Reviews](https://google.github.io/eng-practices/review/)
