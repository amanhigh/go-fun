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
rm -f $FA_COMPUTE_DIR/sbi_rates.csv

if [ ! -d "$FA_COMPUTE_DIR" ]; then
    echo "Creating directory and copying test data to $FA_COMPUTE_DIR..."
    mkdir -p "$FA_COMPUTE_DIR/Tickers"
    cp "$TEST_DATA_DIR/trades.csv" "$FA_COMPUTE_DIR/"
    cp "$TEST_DATA_DIR/dividends.csv" "$FA_COMPUTE_DIR/"
    cp "$TEST_DATA_DIR/interest.csv" "$FA_COMPUTE_DIR/"
    cp "$TEST_DATA_DIR/gains.csv" "$FA_COMPUTE_DIR/"
    
    cp "$TEST_DATA_DIR/accounts.csv" "$FA_COMPUTE_DIR/"
    cp "$TEST_DATA_DIR/AAPL.json" "$FA_COMPUTE_DIR/Tickers/"
    
    # Add NVDA ticker to trades.csv for auto-download testing
    echo "NVDA,2024-06-15,BUY,25,300.00,7500.00,2.50" >> "$FA_COMPUTE_DIR/trades.csv"
else
    echo "Directory $FA_COMPUTE_DIR already exists, updating test data..."
    # Ensure we have the base trades.csv and add NVDA ticker for auto-download testing
    cp "$TEST_DATA_DIR/trades.csv" "$FA_COMPUTE_DIR/"
    if ! grep -q "NVDA" "$FA_COMPUTE_DIR/trades.csv"; then
        echo "NVDA,2024-06-15,BUY,25,300.00,7500.00,2.50" >> "$FA_COMPUTE_DIR/trades.csv"
    fi
fi

# 4. Print environment for debugging
echo "--- Running with the following configuration ---"
echo "PROJECT_ROOT: $PROJECT_ROOT"
echo "FA_COMPUTE_DIR: $FA_COMPUTE_DIR"
echo "------------------------------------------------"

# 5. Run the application's tax command from the project root
echo "Executing 'go run ./components/kohan apps tax 2024' from $PROJECT_ROOT..."
(cd "$PROJECT_ROOT" && go run ./components/kohan apps tax compute 2024) || echo "Application returned non-zero exit code, continuing for verification..."

# 6. Verify ticker auto-download functionality
echo "--- Verifying Ticker Auto-Download ---"
if [ -f "$FA_COMPUTE_DIR/Tickers/NVDA.json" ]; then
    echo "✅ SUCCESS: NVDA.json was auto-downloaded"
    echo "File size: $(wc -c < "$FA_COMPUTE_DIR/Tickers/NVDA.json") bytes"
else
    echo "❌ FAILURE: NVDA.json was NOT auto-downloaded"
    echo "Available tickers: $(ls -la "$FA_COMPUTE_DIR/Tickers/" || echo "No ticker directory")"
    exit 1
fi
echo "-----------------------------------"

# 7. Verify that the output files were created
echo "Verifying output..."

echo "--- Checking sbi_rates.csv ---"
if [ -f "$FA_COMPUTE_DIR/sbi_rates.csv" ]; then
    echo "✅ sbi_rates.csv was created."
    echo "Line count:"
    wc -l "$FA_COMPUTE_DIR/sbi_rates.csv"
else
    echo "❌ FAILURE: sbi_rates.csv was NOT created."
fi
echo "-----------------------------------"

echo "--- Checking accounts_2024.csv ---"
if [ -f "$FA_COMPUTE_DIR/accounts_2024.csv" ]; then
    echo "✅ accounts_2024.csv was created."
    echo "Line count:"
    wc -l "$FA_COMPUTE_DIR/accounts_2024.csv"
else
    echo "❌ FAILURE: accounts_2024.csv was NOT created."
    exit 1
fi
echo "-----------------------------------"

if [ -f "$FA_COMPUTE_DIR/tax_summary.xlsx" ]; then
  echo "✅ SUCCESS: Tax summary Excel file was created at $FA_COMPUTE_DIR/tax_summary.xlsx"
else
  echo "❌ FAILURE: Tax summary Excel file was NOT created."
  exit 1
fi

# 8. Cleanup auto-downloaded files for clean test environment
echo "--- Cleaning up auto-downloaded files ---"
if [ -f "$FA_COMPUTE_DIR/Tickers/NVDA.json" ]; then
    rm -f "$FA_COMPUTE_DIR/Tickers/NVDA.json"
    echo "✅ Cleaned up NVDA.json"
else
    echo "ℹ️  No NVDA.json to clean up"
fi

echo "✅ E2E Test with Ticker Auto-Download PASSED!"
