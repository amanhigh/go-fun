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

echo "Copying test data to $FA_COMPUTE_DIR..."

# Ensure base directory exists
mkdir -p "$FA_COMPUTE_DIR"

# Copy entire test data directory structure (Input, Data, Output layers)
# The -r flag recursively copies directories
cp -r "$TEST_DATA_DIR/Input" "$FA_COMPUTE_DIR/"
cp -r "$TEST_DATA_DIR/Data" "$FA_COMPUTE_DIR/"
cp -r "$TEST_DATA_DIR/Output" "$FA_COMPUTE_DIR/"

# Ensure Reports directory exists (it will be empty initially, filled by application)
mkdir -p "$FA_COMPUTE_DIR/Output/Reports"

# Add NVDA ticker to trades.csv for auto-download testing (if not already present)
if ! grep -q "NVDA" "$FA_COMPUTE_DIR/Input/Parsed/trades.csv" 2>/dev/null; then
    echo "NVDA,2024-06-15,BUY,25,300.00,7500.00,2.50" >> "$FA_COMPUTE_DIR/Input/Parsed/trades.csv"
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
if [ -f "$FA_COMPUTE_DIR/Data/Tickers/NVDA.json" ]; then
    echo "✅ SUCCESS: NVDA.json was auto-downloaded"
    echo "File size: $(wc -c < "$FA_COMPUTE_DIR/Data/Tickers/NVDA.json") bytes"
else
    echo "❌ FAILURE: NVDA.json was NOT auto-downloaded"
    echo "Available tickers: $(ls -la "$FA_COMPUTE_DIR/Data/Tickers/" || echo "No ticker directory")"
    exit 1
fi
echo "-----------------------------------"

# 7. Verify that the output files were created
echo "Verifying output..."

echo "--- Checking sbi_rates.csv (Layer 2: Data/Reference) ---"
if [ -f "$FA_COMPUTE_DIR/Data/Reference/sbi_rates.csv" ]; then
    echo "✅ sbi_rates.csv exists in Data/Reference/"
    echo "Line count:"
    wc -l "$FA_COMPUTE_DIR/Data/Reference/sbi_rates.csv"
else
    echo "❌ FAILURE: sbi_rates.csv NOT found in Data/Reference/"
fi
echo "-----------------------------------"

echo "--- Checking accounts_2024.csv (Layer 3: Output/YearEndBalance) ---"
if [ -f "$FA_COMPUTE_DIR/Output/YearEndBalance/accounts_2024.csv" ]; then
    echo "✅ accounts_2024.csv was created in Output/YearEndBalance/"
    echo "Line count:"
    wc -l "$FA_COMPUTE_DIR/Output/YearEndBalance/accounts_2024.csv"
else
    echo "❌ FAILURE: accounts_2024.csv was NOT created in Output/YearEndBalance/"
    exit 1
fi
echo "-----------------------------------"

echo "--- Checking tax_summary_2024.xlsx (Layer 3: Output/Reports) ---"
if [ -f "$FA_COMPUTE_DIR/Output/Reports/tax_summary_2024.xlsx" ]; then
  echo "✅ SUCCESS: Tax summary Excel file was created at Output/Reports/tax_summary_2024.xlsx"
else
  echo "❌ FAILURE: Tax summary Excel file was NOT created in Output/Reports/"
  exit 1
fi
echo "-----------------------------------"

# 8. Cleanup auto-downloaded files for clean test environment
echo "--- Cleaning up auto-downloaded files ---"
if [ -f "$FA_COMPUTE_DIR/Data/Tickers/NVDA.json" ]; then
    rm -f "$FA_COMPUTE_DIR/Data/Tickers/NVDA.json"
    echo "✅ Cleaned up NVDA.json"
else
    echo "ℹ️  No NVDA.json to clean up"
fi

echo "✅ E2E Test with Ticker Auto-Download PASSED!"
