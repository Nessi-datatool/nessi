#!/bin/bash

# Documentation Verification Script for Nessi
# This script verifies that code examples in documentation are accurate and functional

set -e

# Colors for output
GREEN="\033[0;32m"
YELLOW="\033[1;33m"
RED="\033[0;31m"
NC="\033[0m" # No Color

echo -e "${YELLOW}Nessi Documentation Verification${NC}"
echo "==============================="

# Check if nessi binary exists
if [ ! -f "./nessi" ]; then
    echo -e "${YELLOW}Building nessi binary...${NC}"
    go build -o nessi ./cmd/nessi
fi

# Function to extract and verify code examples from markdown files
verify_markdown_examples() {
    local file=$1
    local temp_dir=$(mktemp -d)
    local example_count=0
    local success_count=0
    
    echo -e "\n${YELLOW}Verifying examples in $file${NC}"
    
    # Extract code blocks with ```go ... ``` format
    # Use a more robust approach with sed and a counter
    count=0
    sed -n '/```go/,/```/ { /```go/d; /```/d; /./p; }' $file | while read -r line; do
        if [[ ! -z $line ]]; then
            if [[ ! -f "$temp_dir/example_$count.go" ]]; then
                count=$((count+1))
                touch "$temp_dir/example_$count.go"
            fi
            echo "$line" >> "$temp_dir/example_$count.go"
        else
            if [[ -f "$temp_dir/example_$count.go" && -s "$temp_dir/example_$count.go" ]]; then
                count=$((count+1))
            fi
        fi
    done
    
    # Count the number of extracted examples
    example_count=$(ls $temp_dir/example_*.go 2>/dev/null | wc -l)
    
    if [ $example_count -eq 0 ]; then
        echo "  No Go code examples found"
        rm -rf $temp_dir
        return
    fi
    
    echo "  Found $example_count code examples"
    
    # Verify each example
    for example in $temp_dir/example_*.go; do
        example_name=$(basename $example)
        echo -n "  Verifying $example_name: "
        
        # Check if the example is a complete program or just a snippet
        if grep -q "package main" $example && grep -q "func main" $example; then
            # It's a complete program, try to compile it
            if go build -o /dev/null $example &>/dev/null; then
                echo -e "${GREEN}✓ Compiles successfully${NC}"
                success_count=$((success_count + 1))
            else
                echo -e "${RED}✗ Compilation failed${NC}"
            fi
        else
            # It's a code snippet, check for syntax errors
            if go vet $example &>/dev/null; then
                echo -e "${GREEN}✓ Syntax looks good${NC}"
                success_count=$((success_count + 1))
            else
                echo -e "${RED}✗ Syntax issues detected${NC}"
            fi
        fi
    done
    
    # Report success rate
    if [ $example_count -gt 0 ]; then
        success_rate=$((success_count * 100 / example_count))
        echo -e "  ${YELLOW}Success rate: $success_rate% ($success_count/$example_count)${NC}"
    fi
    
    # Clean up
    rm -rf $temp_dir
}

# Function to verify CLI examples
verify_cli_examples() {
    local file=$1
    local temp_dir=$(mktemp -d)
    local example_count=0
    local success_count=0
    
    echo -e "\n${YELLOW}Verifying CLI examples in $file${NC}"
    
    # Extract code blocks with ```bash ... ``` or ```shell ... ``` format that contain nessi commands
    count=0
    for block_start in "```bash" "```shell"; do
        sed -n "/$block_start/,/```/ { /$block_start/d; /```/d; /./p; }" $file | while read -r line; do
            if [[ ! -z $line ]]; then
                if [[ ! -f "$temp_dir/example_$count.sh" ]]; then
                    count=$((count+1))
                    touch "$temp_dir/example_$count.sh"
                fi
                echo "$line" >> "$temp_dir/example_$count.sh"
            else
                if [[ -f "$temp_dir/example_$count.sh" && -s "$temp_dir/example_$count.sh" ]]; then
                    count=$((count+1))
                fi
            fi
        done
    done
    
    # Filter to keep only examples with nessi commands
    shopt -s nullglob
    for example in $temp_dir/example_*.sh; do
        if [[ -f "$example" ]]; then
            if ! grep -q "nessi" "$example"; then
                rm "$example"
            fi
        fi
    done
    shopt -u nullglob
    
    # Count the number of extracted examples
    example_count=$(ls $temp_dir/example_*.sh 2>/dev/null | wc -l)
    
    if [ $example_count -eq 0 ]; then
        echo "  No CLI examples found"
        rm -rf $temp_dir
        return
    fi
    
    echo "  Found $example_count CLI examples"
    
    # Verify each example
    for example in $temp_dir/example_*.sh; do
        example_name=$(basename $example)
        echo -n "  Verifying $example_name: "
        
        # Extract nessi commands (only the command, not the output)
        grep "^nessi " $example > $temp_dir/commands.txt || true
        
        # Check if each command is valid (just check help output, don't actually run the command)
        while read -r cmd; do
            # Extract just the command part (without arguments)
            cmd_parts=($cmd)
            base_cmd="${cmd_parts[0]} ${cmd_parts[1]}"
            
            # Check if the command exists in help
            if ./nessi help | grep -q "${cmd_parts[1]}"; then
                echo -e "${GREEN}✓ Valid command: $cmd${NC}"
                success_count=$((success_count + 1))
            else
                echo -e "${RED}✗ Unknown command: $cmd${NC}"
            fi
        done < $temp_dir/commands.txt
    done
    
    # Report success rate
    if [ $example_count -gt 0 ]; then
        success_rate=$((success_count * 100 / example_count))
        echo -e "  ${YELLOW}Success rate: $success_rate% ($success_count/$example_count)${NC}"
    fi
    
    # Clean up
    rm -rf $temp_dir
}

# Verify viral growth documentation
echo -e "\n${YELLOW}Verifying Viral Growth Documentation${NC}"
echo "----------------------------------"

# List of documentation files to verify
docs=(
    "docs/VIRAL_GROWTH.md"
    "docs/PLUGINS.md"
    "docs/PLUGIN_API.md"
    "docs/integration_guide.md"
    "docs/QUALITY_ASSURANCE_PLAN.md"
    "docs/VIRAL_FEATURES_TEST_PLAN.md"
)

# Verify each documentation file
for doc in "${docs[@]}"; do
    if [ -f "$doc" ]; then
        verify_markdown_examples "$doc"
        verify_cli_examples "$doc"
    else
        echo -e "${RED}Warning: Documentation file $doc not found${NC}"
    fi
done

echo -e "\n${GREEN}Documentation verification completed!${NC}"
