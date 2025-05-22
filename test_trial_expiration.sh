#!/bin/bash

echo "Nessi Trial Expiration Test"
echo "--------------------------------------------------"
echo "This test verifies the behavior when a trial expires"
echo ""

# Clear existing trial files
echo "Clearing existing trial files..."
rm -f ~/.nessi/trial.json ~/.nessi/trial_registry.json
echo "Done."

# Build the debug tool
echo "Building debug tool..."
cd cmd/debug_trial && go build -o debug_trial
cd ../..
echo "Done."

# Start a new trial
echo "Starting a new trial..."
NESSI_DEV_MODE=true ./cmd/debug_trial/debug_trial
echo "Trial started successfully."

# Test premium feature access
echo "Testing premium feature access with active trial..."
NESSI_DEV_MODE=true ./cmd/test_premium_features/test_premium_features | grep -A 1 "Databricks Integration"
echo ""

# Create a modified trial file with an expired date (1 month and 1 day ago)
echo "Simulating trial expiration..."
mkdir -p ~/.nessi
cat > ~/.nessi/trial.json << EOF
{
  "data": {
    "start_time": "2025-04-21T20:36:01.716116+02:00",
    "active": true,
    "machine_id": "$(NESSI_DEV_MODE=true ./cmd/debug_trial/debug_trial | grep "Current machine ID" | cut -d' ' -f4)",
    "trial_count": 1,
    "last_renewal": "2025-04-21T20:36:01.716116+02:00"
  },
  "signature": "invalid-signature-for-testing"
}
EOF
echo "Trial has been expired for testing purposes."

# Test premium feature access after expiration
echo "Testing premium feature access with expired trial..."
NESSI_DEV_MODE=true ./cmd/test_premium_features/test_premium_features | grep -A 1 "Databricks Integration"

echo ""
echo "Test Complete"
