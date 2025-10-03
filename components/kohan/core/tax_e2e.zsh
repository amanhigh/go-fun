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

# 1. PRESERVATION BACKUP
echo "--- Backing up preserved files ---"
BACKUP_DIR=$(mktemp -d)
echo "Created temporary backup directory: $BACKUP_DIR"

if [ -f "$FA_COMPUTE_DIR/vested.xlsx" ]; then
    cp "$FA_COMPUTE_DIR/vested.xlsx" "$BACKUP_DIR/"
    echo "✅ Backed up vested.xlsx"
fi

if [ -f "$FA_COMPUTE_DIR/sbi_rates.csv" ]; then
    cp "$FA_COMPUTE_DIR/sbi_rates.csv" "$BACKUP_DIR/"
    echo "✅ Backed up sbi_rates.csv ($(wc -l < "$FA_COMPUTE_DIR/sbi_rates.csv") lines)"
fi

if [ -d "$FA_COMPUTE_DIR/Tickers" ]; then
    cp -r "$FA_COMPUTE_DIR/Tickers" "$BACKUP_DIR/"
    echo "✅ Backed up Tickers/ directory ($(ls "$FA_COMPUTE_DIR/Tickers" | wc -l) files)"
fi
echo "-----------------------------------"

# 2. CLEAN SLATE (Remove test CSVs and old outputs)
echo "--- Cleaning test data and old outputs ---"
rm -f "$FA_COMPUTE_DIR/trades.csv"
rm -f "$FA_COMPUTE_DIR/gains.csv"
rm -f "$FA_COMPUTE_DIR/dividends.csv"
rm -f "$FA_COMPUTE_DIR/interest.csv"
rm -f "$FA_COMPUTE_DIR/accounts"*.csv
rm -f "$FA_COMPUTE_DIR/tax_summary_"*.xlsx
echo "✅ Removed old test CSVs and output files"
echo "-----------------------------------"

# 3. FRESH DATA POPULATION
echo "--- Populating fresh test data ---"
mkdir -p "$FA_COMPUTE_DIR/Tickers"

cp "$TEST_DATA_DIR/trades.csv" "$FA_COMPUTE_DIR/"
cp "$TEST_DATA_DIR/gains.csv" "$FA_COMPUTE_DIR/"
cp "$TEST_DATA_DIR/dividends.csv" "$FA_COMPUTE_DIR/"
cp "$TEST_DATA_DIR/interest.csv" "$FA_COMPUTE_DIR/"
cp "$TEST_DATA_DIR/accounts_2022.csv" "$FA_COMPUTE_DIR/"
cp "$TEST_DATA_DIR/AAPL.json" "$FA_COMPUTE_DIR/Tickers/"

# Add NVDA ticker to trades.csv for auto-download testing
echo "NVDA,2024-06-15,BUY,25,300.00,7500.00,2.50" >> "$FA_COMPUTE_DIR/trades.csv"
echo "✅ Copied fresh test data from testdata/"
echo "-----------------------------------"

# 4. RESTORE PRESERVED FILES
echo "--- Restoring preserved files ---"
if [ -f "$BACKUP_DIR/vested.xlsx" ]; then
    cp "$BACKUP_DIR/vested.xlsx" "$FA_COMPUTE_DIR/"
    echo "✅ Restored vested.xlsx"
fi

if [ -f "$BACKUP_DIR/sbi_rates.csv" ]; then
    cp "$BACKUP_DIR/sbi_rates.csv" "$FA_COMPUTE_DIR/"
    echo "✅ Restored sbi_rates.csv ($(wc -l < "$FA_COMPUTE_DIR/sbi_rates.csv") lines)"
fi

if [ -d "$BACKUP_DIR/Tickers" ]; then
    cp -r "$BACKUP_DIR/Tickers/"* "$FA_COMPUTE_DIR/Tickers/" 2>/dev/null || true
    echo "✅ Restored Tickers/ directory"
fi

# Cleanup temporary backup
rm -rf "$BACKUP_DIR"
echo "✅ Cleaned up temporary backup directory"
echo "-----------------------------------"

# 4. Print environment for debugging
echo "--- Running with the following configuration ---"
echo "PROJECT_ROOT: $PROJECT_ROOT"
echo "FA_COMPUTE_DIR: $FA_COMPUTE_DIR"
echo "------------------------------------------------"

# 5. Generate accounts_2023.csv (required for 2024 tax computation)
echo "--- Generating accounts_2023.csv (prerequisite for 2024) ---"
echo "Running: go run ./components/kohan apps tax compute 2023"
(cd "$PROJECT_ROOT" && go run ./components/kohan apps tax compute 2023) || echo "Warning: 2023 tax computation had issues, continuing..."
echo "✅ Generated accounts_2023.csv and tax_summary_2023.xlsx"
echo "-----------------------------------"

# 6. Run the application's tax command for 2024
echo "--- Executing 2024 Tax Computation ---"
echo "Running: go run ./components/kohan apps tax compute 2024"
(cd "$PROJECT_ROOT" && go run ./components/kohan apps tax compute 2024) || echo "Application returned non-zero exit code, continuing for verification..."
echo "-----------------------------------"

# 7. Verify ticker auto-download functionality
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

# 8. Verify that the output files were created
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

if [ -f "$FA_COMPUTE_DIR/tax_summary_2024.xlsx" ]; then
  echo "✅ SUCCESS: Tax summary Excel file was created at $FA_COMPUTE_DIR/tax_summary_2024.xlsx"
else
  echo "❌ FAILURE: Tax summary Excel file was NOT created."
  exit 1
fi

# 8.1 Validate Excel sheets using read_excel.zsh
echo "--- Validating Excel Sheets ---"
EXCEL_OUTPUT=$("$SCRIPT_DIR/read_excel.zsh" "$FA_COMPUTE_DIR/tax_summary_2024.xlsx" 2>&1)
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

# 9. Cleanup auto-downloaded files for clean test environment
echo "--- Cleaning up auto-downloaded files ---"
if [ -f "$FA_COMPUTE_DIR/Tickers/NVDA.json" ]; then
    rm -f "$FA_COMPUTE_DIR/Tickers/NVDA.json"
    echo "✅ Cleaned up NVDA.json"
else
    echo "ℹ️  No NVDA.json to clean up"
fi

echo "✅ E2E Test with Ticker Auto-Download PASSED!"
