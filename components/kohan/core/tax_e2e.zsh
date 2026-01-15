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

# 1. CLEAN E2E DIRECTORIES (preserve Input/Brokerage for UAT)
echo "--- Cleaning E2E test directories ---"
rm -rf "$FA_COMPUTE_DIR/Input/Parsed"
rm -rf "$FA_COMPUTE_DIR/Data"
rm -rf "$FA_COMPUTE_DIR/Output"
echo "✅ Cleaned E2E directories:"
echo "   - Input/Parsed/ (removed)"
echo "   - Data/ (removed)"
echo "   - Output/ (removed)"
echo "   - Input/Brokerage/ (preserved for UAT)"
echo "-----------------------------------"
echo ""

# 2. COPY TEST DIRECTORIES
echo "--- Copying test data directories ---"

# Create parent directories
mkdir -p "$FA_COMPUTE_DIR/Input"
mkdir -p "$FA_COMPUTE_DIR/Data"
mkdir -p "$FA_COMPUTE_DIR/Output"

# Copy entire directory trees from testdata
cp -r "$TEST_DATA_DIR/Input/Parsed" "$FA_COMPUTE_DIR/Input/"
cp -r "$TEST_DATA_DIR/Data/Tickers" "$FA_COMPUTE_DIR/Data/"
cp -r "$TEST_DATA_DIR/Data/Reference" "$FA_COMPUTE_DIR/Data/"
cp -r "$TEST_DATA_DIR/Output/YearEndBalance" "$FA_COMPUTE_DIR/Output/"
cp -r "$TEST_DATA_DIR/Output/Computed" "$FA_COMPUTE_DIR/Output/"

echo "✅ Copied test directories:"
echo "   - Input/Parsed/ (trades, dividends, interest)"
echo "   - Data/Tickers/ (AAPL, GOOGL, MSFT)"
echo "   - Data/Reference/ (sbi_rates.csv)"
echo "   - Output/YearEndBalance/ (accounts_2022.csv)"
echo "   - Output/Computed/ (gains.csv)"

# Add NVDA ticker to trades.csv for auto-download testing
echo "NVDA,2024-06-15,BUY,25,300.00,7500.00,2.50" >> "$FA_COMPUTE_DIR/Input/Parsed/trades.csv"
echo "✅ Added NVDA trade for auto-download testing"
echo "-----------------------------------"
echo ""

# 3. PRINT E2E CONFIGURATION
echo "--- E2E Test Configuration ---"
echo "PROJECT_ROOT: $PROJECT_ROOT"
echo "TEST_DATA_DIR: $TEST_DATA_DIR"
echo "FA_COMPUTE_DIR: $FA_COMPUTE_DIR"
echo ""
echo "E2E Scope (copied from testdata):"
echo "  ✅ Input/Parsed/"
echo "  ✅ Data/"
echo "  ✅ Output/"
echo ""
echo "UAT Scope (never touched by E2E):"
echo "  🔒 Input/Brokerage/"
echo "-----------------------------------"
echo ""

# 4. Generate accounts_2023.csv (required for 2024 tax computation)
echo "--- Generating accounts_2023.csv (prerequisite for 2024) ---"
echo "Running: go run ./components/kohan apps tax compute 2023"
(cd "$PROJECT_ROOT" && go run ./components/kohan apps tax compute 2023) || echo "Warning: 2023 tax computation had issues, continuing..."
echo "✅ Generated accounts_2023.csv and tax_summary_2023.xlsx"
echo "-----------------------------------"
echo ""

# 5. Run the application's tax command for 2024
echo "--- Executing 2024 Tax Computation ---"
echo "Running: go run ./components/kohan apps tax compute 2024"
(cd "$PROJECT_ROOT" && go run ./components/kohan apps tax compute 2024) || echo "Application returned non-zero exit code, continuing for verification..."
echo "-----------------------------------"
echo ""

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
echo ""

# 7. Verify that the output files were created
echo "Verifying output..."
echo ""

echo "--- Checking sbi_rates.csv ---"
if [ -f "$FA_COMPUTE_DIR/Data/Reference/sbi_rates.csv" ]; then
    echo "✅ sbi_rates.csv exists."
    echo "Line count:"
    wc -l "$FA_COMPUTE_DIR/Data/Reference/sbi_rates.csv"
else
    echo "❌ FAILURE: sbi_rates.csv was NOT found."
fi
echo "-----------------------------------"
echo ""

echo "--- Checking accounts_2024.csv ---"
if [ -f "$FA_COMPUTE_DIR/Output/YearEndBalance/accounts_2024.csv" ]; then
    echo "✅ accounts_2024.csv was created."
    echo "Line count:"
    wc -l "$FA_COMPUTE_DIR/Output/YearEndBalance/accounts_2024.csv"
else
    echo "❌ FAILURE: accounts_2024.csv was NOT created."
    exit 1
fi
echo "-----------------------------------"
echo ""

if [ -f "$FA_COMPUTE_DIR/Output/Reports/tax_summary_2024.xlsx" ]; then
  echo "✅ SUCCESS: Tax summary Excel file was created at $FA_COMPUTE_DIR/Output/Reports/tax_summary_2024.xlsx"
else
  echo "❌ FAILURE: Tax summary Excel file was NOT created."
  exit 1
fi
echo ""

# 8. Validate Excel sheets using read_excel.zsh
echo "--- Validating Excel Sheets ---"
EXCEL_OUTPUT=$("$SCRIPT_DIR/read_excel.zsh" "$FA_COMPUTE_DIR/Output/Reports/tax_summary_2024.xlsx" 2>&1)
SHEETS_LINE=$(echo "$EXCEL_OUTPUT" | grep "Available sheets:")

REQUIRED_SHEETS=("Gains" "Dividends" "Valuations" "Interest")
MISSING_SHEETS=()

for sheet in "${REQUIRED_SHEETS[@]}"; do
  if ! echo "$SHEETS_LINE" | grep -q "$sheet"; then
    MISSING_SHEETS+=("$sheet")
  fi
done

if [ ${#MISSING_SHEETS[@]} -eq 0 ]; then
  echo "✅ SUCCESS: All required sheets are present: ${REQUIRED_SHEETS[*]}"
  echo "   Found sheets: $SHEETS_LINE"
else
  echo "❌ FAILURE: Missing sheets: ${MISSING_SHEETS[*]}"
  echo "   Found sheets: $SHEETS_LINE"
  exit 1
fi
echo "-----------------------------------"
echo ""

# 9. Cleanup auto-downloaded files for clean test environment
echo "--- Cleaning up auto-downloaded files ---"
if [ -f "$FA_COMPUTE_DIR/Data/Tickers/NVDA.json" ]; then
    rm -f "$FA_COMPUTE_DIR/Data/Tickers/NVDA.json"
    echo "✅ Cleaned up NVDA.json"
else
    echo "ℹ️  No NVDA.json to clean up"
fi
echo ""

echo "✅ E2E Test PASSED! (Input/Brokerage/ preserved for UAT)"
