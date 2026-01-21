#!/bin/zsh

set -e

echo "=== UAT: Parse & Compute 2022-2024 Multi-Year Taxes ==="
echo ""

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
PROJECT_ROOT=$(cd "$SCRIPT_DIR/../../.." && pwd)
FA_COMPUTE_DIR=~/Downloads/FACompute

echo "Cleaning previous test outputs..."
rm -rf "$FA_COMPUTE_DIR/Input/Parsed"
rm -rf "$FA_COMPUTE_DIR/Output"
rm -f "$FA_COMPUTE_DIR/Data/Reference/sbi_rates.csv"
mkdir -p "$FA_COMPUTE_DIR/Input/Parsed"
mkdir -p "$FA_COMPUTE_DIR/Output/YearEndBalance"
mkdir -p "$FA_COMPUTE_DIR/Output/Computed"
mkdir -p "$FA_COMPUTE_DIR/Output/Reports"
mkdir -p "$FA_COMPUTE_DIR/Data/Reference"
echo "✅ Cleaned"
echo ""

echo "Step 1: Parse 2022 Broker Files"
(cd "$PROJECT_ROOT" && go run ./components/kohan apps tax parse 2022)
echo "✅ Parsed 2022"
echo ""

if [ -f "$FA_COMPUTE_DIR/Input/Parsed/trades.csv" ]; then
    echo "✅ trades.csv created"
else
    echo "❌ trades.csv NOT FOUND"
    exit 1
fi
echo ""

echo "Step 2: Compute 2022 Taxes"
(cd "$PROJECT_ROOT" && go run ./components/kohan apps tax compute 2022)
echo "✅ Computed 2022"
echo ""

if [ -f "$FA_COMPUTE_DIR/Output/YearEndBalance/accounts_2022.csv" ]; then
    echo "✅ accounts_2022.csv created"
else
    echo "❌ accounts_2022.csv NOT FOUND"
    exit 1
fi

if [ -f "$FA_COMPUTE_DIR/Output/Reports/tax_summary_2022.xlsx" ]; then
    echo "✅ tax_summary_2022.xlsx created"
else
    echo "❌ tax_summary_2022.xlsx NOT FOUND"
    exit 1
fi
echo ""

echo "Step 3: Parse 2023 Broker Files"
(cd "$PROJECT_ROOT" && go run ./components/kohan apps tax parse 2023)
echo "✅ Parsed 2023"
echo ""

if [ -f "$FA_COMPUTE_DIR/Input/Parsed/trades.csv" ]; then
    TRADES_2022=$(grep -c "^[^,]*,2022-" "$FA_COMPUTE_DIR/Input/Parsed/trades.csv" || true)
    TRADES_2023=$(grep -c "^[^,]*,2023-" "$FA_COMPUTE_DIR/Input/Parsed/trades.csv" || true)
    echo "✅ trades.csv updated ($TRADES_2022 2022 trades + $TRADES_2023 2023 trades)"
else
    echo "❌ trades.csv NOT FOUND"
    exit 1
fi
echo ""

echo "Step 4: Compute 2023 Taxes (with carry-forward)"
if [ ! -f "$FA_COMPUTE_DIR/Output/YearEndBalance/accounts_2022.csv" ]; then
    echo "❌ accounts_2022.csv NOT FOUND - cannot proceed"
    exit 1
fi
(cd "$PROJECT_ROOT" && go run ./components/kohan apps tax compute 2023)
echo "✅ Computed 2023"
echo ""

if [ -f "$FA_COMPUTE_DIR/Output/YearEndBalance/accounts_2023.csv" ]; then
    echo "✅ accounts_2023.csv created"
else
    echo "❌ accounts_2023.csv NOT FOUND"
    exit 1
fi

if [ -f "$FA_COMPUTE_DIR/Output/Reports/tax_summary_2023.xlsx" ]; then
	echo "✅ tax_summary_2023.xlsx created"
else
	echo "❌ tax_summary_2023.xlsx NOT FOUND"
	exit 1
fi
echo ""

echo "Step 5: Parse 2024 Broker Files"
(cd "$PROJECT_ROOT" && go run ./components/kohan apps tax parse 2024)
echo "✅ Parsed 2024"
echo ""

if [ -f "$FA_COMPUTE_DIR/Input/Parsed/trades.csv" ]; then
	TRADES_2022=$(grep -c "^[^,]*,2022-" "$FA_COMPUTE_DIR/Input/Parsed/trades.csv" || true)
	TRADES_2023=$(grep -c "^[^,]*,2023-" "$FA_COMPUTE_DIR/Input/Parsed/trades.csv" || true)
	TRADES_2024=$(grep -c "^[^,]*,2024-" "$FA_COMPUTE_DIR/Input/Parsed/trades.csv" || true)
	echo "✅ trades.csv updated ($TRADES_2022 2022 + $TRADES_2023 2023 + $TRADES_2024 2024 trades)"
else
	echo "❌ trades.csv NOT FOUND"
	exit 1
fi
echo ""

echo "Step 6: Compute 2024 Taxes (with carry-forward from 2023)"
if [ ! -f "$FA_COMPUTE_DIR/Output/YearEndBalance/accounts_2023.csv" ]; then
	echo "❌ accounts_2023.csv NOT FOUND - cannot proceed"
	exit 1
fi
(cd "$PROJECT_ROOT" && go run ./components/kohan apps tax compute 2024)
echo "✅ Computed 2024"
echo ""

if [ -f "$FA_COMPUTE_DIR/Output/YearEndBalance/accounts_2024.csv" ]; then
	echo "✅ accounts_2024.csv created"
else
	echo "❌ accounts_2024.csv NOT FOUND"
	exit 1
fi

if [ -f "$FA_COMPUTE_DIR/Output/Reports/tax_summary_2024.xlsx" ]; then
	echo "✅ tax_summary_2024.xlsx created"
else
	echo "❌ tax_summary_2024.xlsx NOT FOUND"
	exit 1
fi
echo ""

echo "✅ UAT COMPLETE"
echo ""
echo "Generated Files:"
echo "  ✅ Input/Parsed/trades.csv"
echo "  ✅ Output/YearEndBalance/accounts_2022.csv"
echo "  ✅ Output/YearEndBalance/accounts_2023.csv"
echo "  ✅ Output/YearEndBalance/accounts_2024.csv"
echo "  ✅ Output/Reports/tax_summary_2022.xlsx"
echo "  ✅ Output/Reports/tax_summary_2023.xlsx"
echo "  ✅ Output/Reports/tax_summary_2024.xlsx"
