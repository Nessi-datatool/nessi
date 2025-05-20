# Branch Protection Rules for Nessi

This document outlines how to set up branch protection rules in GitHub to ensure code quality and prevent merging PRs that don't pass the CI pipeline.

## Setting Up Branch Protection Rules

1. Go to your GitHub repository settings
2. Click on "Branches" in the left sidebar
3. Under "Branch protection rules", click "Add rule"
4. In the "Branch name pattern" field, enter `main`
5. Enable the following options:

   - ✅ Require a pull request before merging
   - ✅ Require approvals (set to at least 1)
   - ✅ Dismiss stale pull request approvals when new commits are pushed
   - ✅ Require status checks to pass before merging
   - ✅ Require branches to be up to date before merging

6. In the "Status checks that are required" section, search for and select:
   - `Format and Lint`
   - `Test`
   - `Branch Protection Check`
   - Any other CI checks you've configured

7. Optionally, enable:
   - ✅ Require conversation resolution before merging
   - ✅ Require signed commits

8. Click "Create" or "Save changes"

## Local Pre-commit Hooks

We've set up local pre-commit hooks to catch formatting and linting issues before you push your code. These hooks are automatically installed when you clone the repository.

### What the Pre-commit Hook Does

- Formats Go files using `go fmt`
- Formats the `go.mod` file
- Runs linting using `golangci-lint` if installed

### Manual Setup

If the hooks aren't working, you can set them up manually:

```bash
# Make sure the hook is executable
chmod +x .githooks/pre-commit

# Configure Git to use our hooks directory
git config core.hooksPath .githooks
```

## CI Pipeline Enforcement

Our CI pipeline includes:

1. **Format and Lint Check**: Ensures code is properly formatted and passes linting
2. **Tests**: Runs all tests to ensure functionality
3. **Branch Protection Check**: Ensures all required checks pass before merging

## Best Practices

1. Always run `go fmt ./...` before committing
2. Run `go mod edit -fmt` to format the go.mod file
3. Install and run `golangci-lint run` locally to catch linting issues early
4. Write tests for new functionality
5. Keep PRs focused and small for easier review

## Troubleshooting

If your PR can't be merged due to failing checks:

1. Check the CI logs to identify the specific issue
2. Fix formatting issues with `go fmt ./...`
3. Fix go.mod format issues with `go mod edit -fmt`
4. Address linting issues identified by golangci-lint
5. Make sure all tests are passing locally with `go test ./...`
