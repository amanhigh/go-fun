#!/bin/zsh

# Exit immediately if a command exits with a non-zero status.
set -e

# --- E2E Test for the 'tax' command ---

echo "Setting up E2E test environment..."

# Determine the project root directory based on the script's location
SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
PROJECT_ROOT=$(cd "$SCRIPT_DIR/../../.." && pwd)
TEST_DATA_DIR="$PROJECT_ROOT/components/kohan/testdata/tax"
FA_COMPUTE_DIR=~/Downloads/FACompute

if [ ! -d "$FA_COMPUTE_DIR" ]; then
    echo "Creating directory and copying test data to $FA_COMPUTE_DIR..."
    mkdir -p "$FA_COMPUTE_DIR/Tickers"
    cp "$TEST_DATA_DIR/trades.csv" "$FA_COMPUTE_DIR/"
    cp "$TEST_DATA_DIR/dividends.csv" "$FA_COMPUTE_DIR/"
    cp "$TEST_DATA_DIR/interest.csv" "$FA_COMPUTE_DIR/"
    cp "$TEST_DATA_DIR/gains.csv" "$FA_COMPUTE_DIR/"
    cp "$TEST_DATA_DIR/sbi_rates.csv" "$FA_COMPUTE_DIR/"
    cp "$TEST_DATA_DIR/accounts.csv" "$FA_COMPUTE_DIR/"
    cp "$TEST_DATA_DIR/AAPL.json" "$FA_COMPUTE_DIR/Tickers/"
else
    echo "Directory $FA_COMPUTE_DIR already exists, skipping creation and copy."
fi

# 4. Print environment for debugging
echo "--- Running with the following configuration ---"
echo "PROJECT_ROOT: $PROJECT_ROOT"
echo "FA_COMPUTE_DIR: $FA_COMPUTE_DIR"
echo "------------------------------------------------"

# 5. Run the application's tax command from the project root
echo "Executing 'go run ./components/kohan apps tax 2024' from $PROJECT_ROOT..."
(cd "$PROJECT_ROOT" && go run ./components/kohan apps tax 2024)

# 6. Verify that the output file was created
echo "Verifying output..."
if [ -f "$FA_COMPUTE_DIR/tax_summary.xlsx" ]; then
  echo "✅ SUCCESS: Tax summary Excel file was created at $FA_COMPUTE_DIR/tax_summary.xlsx"
else
  echo "❌ FAILURE: Tax summary Excel file was NOT created."
  exit 1
fi

echo "E2E Test Passed!"
