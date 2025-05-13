#!/bin/bash
#
# Nessi Shell Script Integration Example
#
# This example demonstrates how to integrate Nessi with shell scripts
# to run data quality checks and process the results.
#

set -e

# Define colors for terminal output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Default values
TABLES_DIR=""
CONFIG_FILE=""
OUTPUT_FILE="report.html"
FORMAT="html"
WEBHOOK=false
SLACK_WEBHOOK=""

# Function to print usage
usage() {
  echo -e "${BOLD}Usage:${NC} $0 [options]"
  echo -e "  ${BOLD}-t, --tables-dir${NC}     Directory containing Delta tables (required)"
  echo -e "  ${BOLD}-c, --config${NC}         Path to Nessi configuration file"
  echo -e "  ${BOLD}-o, --output${NC}         Path to save the report (default: report.html)"
  echo -e "  ${BOLD}-f, --format${NC}         Report format: html, json, csv, pdf (default: html)"
  echo -e "  ${BOLD}-w, --webhook${NC}        Trigger Nessi webhook with results"
  echo -e "  ${BOLD}-s, --slack-webhook${NC}  Slack webhook URL for notifications"
  echo -e "  ${BOLD}-h, --help${NC}           Show this help message"
  exit 1
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  key="$1"
  case $key in
    -t|--tables-dir)
      TABLES_DIR="$2"
      shift 2
      ;;
    -c|--config)
      CONFIG_FILE="$2"
      shift 2
      ;;
    -o|--output)
      OUTPUT_FILE="$2"
      shift 2
      ;;
    -f|--format)
      FORMAT="$2"
      shift 2
      ;;
    -w|--webhook)
      WEBHOOK=true
      shift
      ;;
    -s|--slack-webhook)
      SLACK_WEBHOOK="$2"
      shift 2
      ;;
    -h|--help)
      usage
      ;;
    *)
      echo -e "${RED}Error: Unknown option $1${NC}"
      usage
      ;;
  esac
done

# Check required arguments
if [ -z "$TABLES_DIR" ]; then
  echo -e "${RED}Error: --tables-dir is required${NC}"
  usage
fi

# Check if Nessi is installed
if ! command -v nessi &> /dev/null; then
  echo -e "${RED}Error: Nessi is not installed or not in PATH${NC}"
  exit 1
fi

# Print Nessi version
NESSI_VERSION=$(nessi version)
echo -e "${BLUE}Using Nessi: $NESSI_VERSION${NC}"

# Create temporary directory for results
TEMP_DIR=$(mktemp -d)
trap 'rm -rf "$TEMP_DIR"' EXIT

# Function to validate a table
validate_table() {
  local table_path="$1"
  local table_name=$(basename "$table_path")
  local output_file="$TEMP_DIR/${table_name}_result.json"
  
  echo -e "${BOLD}Validating table: $table_name${NC}"
  
  # Build command
  local cmd="nessi validate \"$table_path\" --format json"
  if [ -n "$CONFIG_FILE" ]; then
    cmd="$cmd --config \"$CONFIG_FILE\""
  fi
  
  # Run validation
  echo -e "${BLUE}Running: $cmd${NC}"
  eval "$cmd" > "$output_file"
  
  # Check if validation passed
  local passed=$(jq -r '.passed' "$output_file")
  if [ "$passed" = "true" ]; then
    echo -e "${GREEN}✓ Table $table_name passed validation${NC}"
    return 0
  else
    echo -e "${RED}✗ Table $table_name failed validation${NC}"
    # Print failed rules
    jq -r '.rules[] | select(.passed == false) | "  - \(.name): \(.message)"' "$output_file" | 
      while read line; do
        echo -e "${YELLOW}$line${NC}"
      done
    return 1
  fi
}

# Function to send Slack notification
send_slack_notification() {
  local status="$1"
  local message="$2"
  local color="good"
  
  if [ "$status" != "success" ]; then
    color="danger"
  fi
  
  if [ -n "$SLACK_WEBHOOK" ]; then
    echo -e "${BLUE}Sending Slack notification...${NC}"
    
    curl -s -X POST -H 'Content-type: application/json' --data "{
      \"attachments\": [
        {
          \"color\": \"$color\",
          \"title\": \"Nessi Data Quality Check: $status\",
          \"text\": \"$message\",
          \"footer\": \"Nessi Integration Script\",
          \"ts\": $(date +%s)
        }
      ]
    }" "$SLACK_WEBHOOK"
  fi
}

# Function to trigger Nessi webhook
trigger_webhook() {
  local status="$1"
  local results_file="$2"
  
  if [ "$WEBHOOK" = true ]; then
    echo -e "${BLUE}Triggering Nessi webhook...${NC}"
    
    # Create payload with summary
    local total_tables=$(jq -r 'length' "$results_file")
    local passed_tables=$(jq -r '[.[] | select(.passed == true)] | length' "$results_file")
    local failed_tables=$(jq -r '[.[] | select(.passed == false)] | length' "$results_file")
    
    local payload="{
      \"status\": \"$status\",
      \"timestamp\": \"$(date -u +"%Y-%m-%dT%H:%M:%SZ")\",
      \"summary\": {
        \"total_tables\": $total_tables,
        \"passed_tables\": $passed_tables,
        \"failed_tables\": $failed_tables
      },
      \"results\": $(cat "$results_file")
    }"
    
    nessi webhook trigger --event "validation.complete" --payload "$payload"
  fi
}

# Validate all tables
echo -e "${BOLD}Starting validation of tables in: $TABLES_DIR${NC}"
all_passed=true
results=()

for table_path in "$TABLES_DIR"/*; do
  if [ -d "$table_path" ]; then
    if ! validate_table "$table_path"; then
      all_passed=false
    fi
    
    # Add to results array
    table_name=$(basename "$table_path")
    results+=("\"$table_name\": $(cat "$TEMP_DIR/${table_name}_result.json")")
  fi
done

# Combine all results into a single JSON file
echo "{$(IFS=,; echo "${results[*]}")}" > "$TEMP_DIR/all_results.json"

# Generate report
echo -e "${BLUE}Generating $FORMAT report: $OUTPUT_FILE${NC}"
nessi report generate --input "$TEMP_DIR/*_result.json" --output "$OUTPUT_FILE" --format "$FORMAT"

# Print summary
echo -e "\n${BOLD}Validation Results Summary:${NC}"
echo "----------------------------------------"
jq -r 'to_entries[] | "\(.key): \(if .value.passed then "✓ PASSED" else "✗ FAILED" end)"' "$TEMP_DIR/all_results.json" |
  while read line; do
    if [[ $line == *"PASSED"* ]]; then
      echo -e "${GREEN}$line${NC}"
    else
      echo -e "${RED}$line${NC}"
    fi
  done
echo "----------------------------------------"

# Determine overall status
if [ "$all_passed" = true ]; then
  echo -e "${GREEN}${BOLD}Overall Status: PASSED${NC}"
  status="success"
  message="All tables passed data quality checks."
else
  echo -e "${RED}${BOLD}Overall Status: FAILED${NC}"
  status="failure"
  message="Some tables failed data quality checks. See the report for details."
fi

# Trigger webhook
trigger_webhook "$status" "$TEMP_DIR/all_results.json"

# Send Slack notification
send_slack_notification "$status" "$message"

# Exit with appropriate status code
if [ "$all_passed" = true ]; then
  exit 0
else
  exit 1
fi
