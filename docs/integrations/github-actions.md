# GitHub Actions Integration

This guide shows how to integrate Nessi.dev with GitHub Actions for automated data quality checks and monitoring.

## Example Workflow

Create a file `.github/workflows/data-quality.yml`:

```yaml
name: Data Quality Check

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]
  schedule:
    - cron: '0 0 * * *'  # Run daily at midnight

jobs:
  quality-check:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v2
    
    - name: Set up Python
      uses: actions/setup-python@v2
      with:
        python-version: '3.9'
    
    - name: Install Nessi.dev
      run: |
        python -m pip install --upgrade pip
        pip install nessi
    
    - name: Run Data Quality Check
      env:
        NESSI_API_KEY: ${{ secrets.NESSI_API_KEY }}
      run: |
        nessi quality check \
          --table-path s3://bucket/table \
          --rules .nessi/rules.json \
          --output quality-report.html
    
    - name: Upload Quality Report
      uses: actions/upload-artifact@v2
      with:
        name: quality-report
        path: quality-report.html
    
    - name: Check Quality Score
      env:
        NESSI_API_KEY: ${{ secrets.NESSI_API_KEY }}
      run: |
        SCORE=$(nessi quality score --table-path s3://bucket/table)
        if (( $(echo "$SCORE < 0.9" | bc -l) )); then
          echo "Quality score below threshold: $SCORE"
          exit 1
        fi
```

## Configuration

1. Add your Nessi.dev API key to GitHub Secrets:
   - Go to your repository settings
   - Navigate to Secrets
   - Add a new secret named `NESSI_API_KEY`

2. Create a rules file `.nessi/rules.json`:
```json
{
  "completeness": 0.95,
  "accuracy": 0.98,
  "consistency": 0.97,
  "rules": [
    {
      "name": "age_range",
      "condition": "age >= 0 AND age <= 120"
    },
    {
      "name": "email_format",
      "pattern": "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
    }
  ]
}
```

## Advanced Usage

### Custom Alerts

Add custom alert thresholds to the workflow:

```yaml
    - name: Check Alerts
      env:
        NESSI_API_KEY: ${{ secrets.NESSI_API_KEY }}
      run: |
        ALERTS=$(nessi alerts list)
        if [ -n "$ALERTS" ]; then
          echo "Active alerts found:"
          echo "$ALERTS"
          exit 1
        fi
```

### Slack Notifications

Add Slack notifications for quality check results:

```yaml
    - name: Send Slack Notification
      if: always()
      uses: 8398a7/action-slack@v3
      with:
        status: ${{ job.status }}
        fields: repo,message,commit,author,action,eventName,ref,workflow
      env:
        SLACK_WEBHOOK_URL: ${{ secrets.SLACK_WEBHOOK_URL }}
```

### Matrix Testing

Run quality checks against multiple tables:

```yaml
    strategy:
      matrix:
        table:
          - s3://bucket/table1
          - s3://bucket/table2
          - s3://bucket/table3
    
    steps:
    - name: Run Data Quality Check
      env:
        NESSI_API_KEY: ${{ secrets.NESSI_API_KEY }}
      run: |
        nessi quality check \
          --table-path ${{ matrix.table }} \
          --rules .nessi/rules.json \
          --output quality-report-${{ matrix.table }}.html
``` 