#!/bin/bash

# Script to lint Databricks integration code
# Created on 2025-05-19

set -e

echo "Linting Databricks integration code..."

# Run gofmt on Databricks integration code
GOFMT_OUTPUT=$(gofmt -l pkg/catalog/databricks pkg/datalake)
if [ -n "$GOFMT_OUTPUT" ]; then
    echo "The following files need formatting:"
    echo "$GOFMT_OUTPUT"
    echo "Running gofmt to fix formatting issues..."
    gofmt -w pkg/catalog/databricks pkg/datalake
    echo "Formatting issues fixed."
else
    echo "All files are properly formatted."
fi

# Run go vet on Databricks integration code
echo "Running go vet on Databricks integration code..."
go vet ./pkg/catalog/databricks/...
go vet ./pkg/datalake/...
echo "go vet completed successfully."

# Check for unused imports
echo "Checking for unused imports..."
UNUSED_IMPORTS_OUTPUT=$(go build ./pkg/catalog/databricks/... ./pkg/datalake/... 2>&1)
UNUSED_IMPORTS=$(echo "$UNUSED_IMPORTS_OUTPUT" | grep -E "imported and not used" || true)
if [ -n "$UNUSED_IMPORTS" ]; then
    echo "Found unused imports:"
    echo "$UNUSED_IMPORTS"
    echo "Please fix these unused imports."
else
    echo "No unused imports found."
fi

# Check for proper error handling
echo "Checking for proper error handling..."
MISSING_ERROR_CHECKS=$(grep -r "err :=" pkg/catalog/databricks pkg/datalake | grep -v "if err" | grep -v "return err" | grep -v "errors.Wrap" | grep -v "errors.New")
if [ -n "$MISSING_ERROR_CHECKS" ]; then
    echo "Potential missing error checks:"
    echo "$MISSING_ERROR_CHECKS"
    echo "Please review these lines to ensure proper error handling."
else
    echo "Error handling looks good."
fi

echo "Linting completed successfully!"
