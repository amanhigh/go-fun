#!/bin/zsh

# Exit immediately if a command exits with a non-zero status.
set -e

# --- UAT: Single-Year Tax Parser (2022) ---
# Note: Testing with 2022 only (DriveWealth vested_2022.xlsx)
# No IBKR data for 2022 yet

echo "=== UAT: Parse & Compute 2022 Taxes ==="
echo ""

# Determine the project root directory based on the script's location
SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
PROJECT_ROOT=$(cd "$SCRIPT_DIR/../../.." && pwd)
FA_COMPUTE_DIR=~/Downloads/FACompute

# Step 0: Clean previous test outputs (preserve broker files, refresh reference data)
echo "--- Cleaning previous test outputs ---"
rm -rf "$FA_COMPUTE_DIR/Input/Parsed"
rm -rf "$FA_COMPUTE_DIR/Output"
rm -f "$FA_COMPUTE_DIR/Data/Reference/sbi_rates.csv"  # Force refresh of exchange rates
mkdir -p "$FA_COMPUTE_DIR/Input/Parsed"
mkdir -p "$FA_COMPUTE_DIR/Output/YearEndBalance"
mkdir -p "$FA_COMPUTE_DIR/Output/Computed"
mkdir -p "$FA_COMPUTE_DIR/Output/Reports"
mkdir -p "$FA_COMPUTE_DIR/Data/Reference"
echo "✅ Cleaned outputs (preserved: Input/Brokerage/, refreshed exchange rates)"
echo "-----------------------------------"
echo ""

# Step 0B: UAT Configuration Documentation
echo "--- UAT Configuration ---"
echo "UAT Workflow (Single-Year Test - 2022):"
echo "  📁 Input/Brokerage/  ← Reads: vested_2022.xlsx (DriveWealth)"
echo "  📁 Input/Parsed/     ← Writes: trades.csv, dividends.csv, interest.csv"
echo "  📁 Data/             ← Uses: sbi_rates.csv, Tickers/ (reference files)"
echo "  📁 Output/           ← Writes: accounts_2022.csv, tax_summary_2022.xlsx"
echo "Note: E2E script cleans Input/Parsed and Output, but never touches Input/Brokerage/"
echo "-----------------------------------"
echo ""

# Step 1: Parse broker files (2022)
echo "--- Step 1: Parse Broker Files (2022) ---"
echo "Parsing: vested_2022.xlsx (DriveWealth)"
echo "Running: go run ./components/kohan apps tax parse 2022"
(cd "$PROJECT_ROOT" && go run ./components/kohan apps tax parse 2022)
echo "✅ Broker files parsed successfully"
echo "-----------------------------------"
echo ""

# Verify parsed files
echo "--- Verify Parsed CSVs (2022) ---"
if [ -f "$FA_COMPUTE_DIR/Input/Parsed/trades.csv" ]; then
    echo "✅ trades.csv created ($(wc -l < "$FA_COMPUTE_DIR/Input/Parsed/trades.csv") lines)"
else
    echo "❌ trades.csv NOT FOUND"
fi

if [ -f "$FA_COMPUTE_DIR/Input/Parsed/dividends.csv" ]; then
    echo "✅ dividends.csv created ($(wc -l < "$FA_COMPUTE_DIR/Input/Parsed/dividends.csv") lines)"
else
    echo "❌ dividends.csv NOT FOUND"
fi

if [ -f "$FA_COMPUTE_DIR/Input/Parsed/interest.csv" ]; then
    echo "✅ interest.csv created ($(wc -l < "$FA_COMPUTE_DIR/Input/Parsed/interest.csv") lines)"
else
    echo "❌ interest.csv NOT FOUND"
fi

if [ -f "$FA_COMPUTE_DIR/Output/Computed/gains.csv" ]; then
    echo "✅ gains.csv created ($(wc -l < "$FA_COMPUTE_DIR/Output/Computed/gains.csv") lines)"
else
    echo "❌ gains.csv NOT FOUND"
fi
echo "-----------------------------------"
echo ""

# Step 2: Compute 2022 taxes
echo "--- Step 2: Compute 2022 Taxes ---"
echo "Running: go run ./components/kohan apps tax compute 2022"
(cd "$PROJECT_ROOT" && go run ./components/kohan apps tax compute 2022)
echo "✅ 2022 tax computation successful"
echo "-----------------------------------"
echo ""

# Verify outputs
echo "--- Verify Tax Outputs (2022) ---"
if [ -f "$FA_COMPUTE_DIR/Output/YearEndBalance/accounts_2022.csv" ]; then
    echo "✅ accounts_2022.csv created ($(wc -l < "$FA_COMPUTE_DIR/Output/YearEndBalance/accounts_2022.csv") lines)"
else
    echo "❌ accounts_2022.csv NOT FOUND"
fi

if [ -f "$FA_COMPUTE_DIR/Output/Reports/tax_summary_2022.xlsx" ]; then
    echo "✅ tax_summary_2022.xlsx created"
else
    echo "❌ tax_summary_2022.xlsx NOT FOUND"
fi
echo "-----------------------------------"
echo ""

# Validate Excel sheets for 2022
echo "--- Validating 2022 Excel Sheets ---"
if [ -f "$FA_COMPUTE_DIR/Output/Reports/tax_summary_2022.xlsx" ]; then
    EXCEL_OUTPUT=$("$SCRIPT_DIR/read_excel.zsh" "$FA_COMPUTE_DIR/Output/Reports/tax_summary_2022.xlsx" 2>&1)
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
else
    echo "❌ tax_summary_2022.xlsx not found - skipping validation"
fi
echo "-----------------------------------"
echo ""

echo "✅ UAT COMPLETE - 2022 Single-Year Tax Workflow Validated!"
echo ""
echo "Summary (2022 Only):"
echo "  1. Cleaned previous outputs (Input/Parsed and Output directories)"
echo "  2. Parsed broker files (vested_2022.xlsx from DriveWealth)"
echo "  3. Generated parsed CSVs (trades, dividends, interest)"
echo "  4. Computed 2022 tax summary"
echo "  5. Created accounts_2022.csv and tax_summary_2022.xlsx"
echo "  6. Validated Excel sheets (Gains, Dividends, Valuations, Interest)"
echo ""
echo "Next Steps:"
echo "  • Verify 2022 outputs are correct"
echo "  • Add 2023 data and test year-over-year carry-forward"
echo "  • Eventually add IBKR data for multi-broker support"
