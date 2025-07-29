#!/bin/zsh

# Exit immediately if a command exits with a non-zero status.
set -e

# --- E2E Test for the 'tax' command ---

echo "Setting up E2E test environment..."

# Determine the project root directory based on the script's location
SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
PROJECT_ROOT=$(cd "$SCRIPT_DIR/../../.." && pwd)
TEST_DATA_DIR="$PROJECT_ROOT/components/kohan/testdata/tax"

# 1. Create a temporary directory for outputs
TEMP_DIR=$(mktemp -d)
OUTPUT_DIR="$TEMP_DIR/outputs"
mkdir -p "$OUTPUT_DIR"

# 2. Schedule cleanup to run when the script exits
cleanup() {
  echo "Cleaning up E2E test environment..."
  rm -rf "$TEMP_DIR"
}
trap cleanup EXIT

# 3. Set environment variables with ABSOLUTE paths
echo "Configuring environment with absolute paths..."
export FA_DOWNLOADS_DIR="$TEST_DATA_DIR"
export FA_BROKER_STATEMENT_PATH="$TEST_DATA_DIR/trades.csv"
export FA_DIVIDEND_FILE_PATH="$TEST_DATA_DIR/dividends.csv"
export FA_INTEREST_FILE_PATH="$TEST_DATA_DIR/interest.csv"
export FA_GAINS_FILE_PATH="$TEST_DATA_DIR/gains.csv"
export SBI_FILE_PATH="$TEST_DATA_DIR/sbi_rates.csv"
export ACCOUNT_FILE_PATH="$TEST_DATA_DIR/accounts.csv"
export YEARLY_SUMMARY_PATH="$OUTPUT_DIR/tax_summary.xlsx"
export ALPHA_API_KEY="DUMMY_KEY_FOR_E2E" # Not used, but required by config

# 4. Print environment for debugging
echo "--- Running with the following configuration ---"
echo "PROJECT_ROOT: $PROJECT_ROOT"
echo "YEARLY_SUMMARY_PATH: $YEARLY_SUMMARY_PATH"
echo "FA_BROKER_STATEMENT_PATH: $FA_BROKER_STATEMENT_PATH"
echo "------------------------------------------------"

# 5. Run the application's tax command from the project root
echo "Executing 'go run ./components/kohan apps tax 2024' from $PROJECT_ROOT..."
(cd "$PROJECT_ROOT" && go run ./components/kohan apps tax 2024)

# 6. Verify that the output file was created
echo "Verifying output..."
if [ -f "$YEARLY_SUMMARY_PATH" ]; then
  echo "✅ SUCCESS: Tax summary Excel file was created at $YEARLY_SUMMARY_PATH"
else
  echo "❌ FAILURE: Tax summary Excel file was NOT created."
  exit 1
fi

echo "E2E Test Passed!"