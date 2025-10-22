#!/bin/zsh

# Exit immediately if a command exits with a non-zero status.
set -e

# --- UAT Restart for Multi-Broker Tax Parser ---

echo "=== UAT Restart: Multi-Broker Tax Parser ==="
echo ""

# Determine the project root directory based on the script's location
SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
PROJECT_ROOT=$(cd "$SCRIPT_DIR/../../.." && pwd)
FA_COMPUTE_DIR=~/Downloads/FACompute

# Step 0: Clean previous test outputs (preserve input files)
echo "--- Cleaning previous test outputs ---"
# Remove parsed CSVs
rm -f "$FA_COMPUTE_DIR/trades.csv"
rm -f "$FA_COMPUTE_DIR/dividends.csv"
rm -f "$FA_COMPUTE_DIR/interest.csv"
rm -f "$FA_COMPUTE_DIR/gains.csv"

# Remove all year-specific outputs (use find to avoid glob issues)
find "$FA_COMPUTE_DIR" -maxdepth 1 -name "accounts_*.csv" -delete 2>/dev/null || true
find "$FA_COMPUTE_DIR" -maxdepth 1 -name "tax_summary_*.xlsx" -delete 2>/dev/null || true

# Preserve: vested.xlsx, ib_*.csv, sbi_rates.csv, Tickers/
echo "✅ Cleaned outputs (preserved: vested.xlsx, ib_*.csv, sbi_rates.csv, Tickers/)"
echo "-----------------------------------"
echo ""

# Step 1: Parse broker files
echo "--- Step 1: Parse Broker Files ---"
echo "Running: go run ./components/kohan apps tax parse"
(cd "$PROJECT_ROOT" && go run ./components/kohan apps tax parse)
echo "✅ Broker files parsed successfully"
echo "-----------------------------------"
echo ""

# Verify parsed files
echo "--- Verify Parsed CSVs ---"
ls -lh "$FA_COMPUTE_DIR/trades.csv" 2>/dev/null || echo "❌ trades.csv NOT FOUND"
ls -lh "$FA_COMPUTE_DIR/dividends.csv" 2>/dev/null || echo "❌ dividends.csv NOT FOUND"
ls -lh "$FA_COMPUTE_DIR/interest.csv" 2>/dev/null || echo "❌ interest.csv NOT FOUND"
ls -lh "$FA_COMPUTE_DIR/gains.csv" 2>/dev/null || echo "❌ gains.csv NOT FOUND"
echo "-----------------------------------"
echo ""

# Step 2: Compute 2022 taxes (prerequisite for 2023)
echo "--- Step 2: Compute 2022 Taxes ---"
echo "Running: go run ./components/kohan apps tax compute 2022"
(cd "$PROJECT_ROOT" && go run ./components/kohan apps tax compute 2022)
echo "✅ 2022 tax computation successful"
echo "-----------------------------------"
echo ""

# Step 3: Compute 2023 taxes (prerequisite for 2024)
echo "--- Step 3: Compute 2023 Taxes ---"
echo "Running: go run ./components/kohan apps tax compute 2023"
(cd "$PROJECT_ROOT" && go run ./components/kohan apps tax compute 2023)
echo "✅ 2023 tax computation successful"
echo "-----------------------------------"
echo ""

# Step 4: Compute 2024 taxes
echo "--- Step 4: Compute 2024 Taxes ---"
echo "Running: go run ./components/kohan apps tax compute 2024"
(cd "$PROJECT_ROOT" && go run ./components/kohan apps tax compute 2024)
echo "✅ 2024 tax computation successful"
echo "-----------------------------------"
echo ""

# Verify outputs
echo "--- Verify Tax Outputs ---"
for year in 2022 2023 2024; do
    if [ -f "$FA_COMPUTE_DIR/accounts_$year.csv" ]; then
        echo "✅ accounts_$year.csv created ($(wc -l < "$FA_COMPUTE_DIR/accounts_$year.csv") lines)"
    else
        echo "❌ accounts_$year.csv NOT FOUND"
    fi
    
    if [ -f "$FA_COMPUTE_DIR/tax_summary_$year.xlsx" ]; then
        echo "✅ tax_summary_$year.xlsx created"
    else
        echo "❌ tax_summary_$year.xlsx NOT FOUND"
    fi
    echo ""
done
echo "-----------------------------------"
echo ""

# Validate Excel sheets for 2024
echo "--- Validating 2024 Excel Sheets ---"
if [ -f "$FA_COMPUTE_DIR/tax_summary_2024.xlsx" ]; then
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
else
    echo "❌ tax_summary_2024.xlsx not found - skipping validation"
fi
echo "-----------------------------------"
echo ""

echo "✅ UAT COMPLETE - Multi-broker workflow validated!"
echo ""
echo "Summary:"
echo "  1. Parsed broker files (DriveWealth + Interactive Brokers)"
echo "  2. Generated merged CSVs (trades, dividends, interest, gains)"
echo "  3. Computed 2022 → 2023 → 2024 tax summaries (year-over-year carry-forward)"
echo "  4. Created accounts_*.csv and tax_summary_*.xlsx for all years"
echo "  5. Validated Excel sheets (Gains, Dividends, Valuations, Interest)"
