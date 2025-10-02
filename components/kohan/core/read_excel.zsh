#!/bin/zsh

# Script to read Excel files and print all sheets
# Usage: ./read_excel.zsh <path_to_excel_file>

set -e

if [ $# -eq 0 ]; then
    echo "Usage: $0 <path_to_excel_file>"
    echo "Example: $0 ~/Downloads/FACompute/tax_summary_2022.xlsx"
    exit 1
fi

EXCEL_FILE="$1"

if [ ! -f "$EXCEL_FILE" ]; then
    echo "Error: File not found: $EXCEL_FILE"
    exit 1
fi

python3 - "$EXCEL_FILE" << 'EOF'
import sys
import os
from openpyxl import load_workbook

excel_file = sys.argv[1]
excel_file = os.path.expanduser(excel_file)

try:
    print(f"Reading Excel file: {excel_file}\n")
    wb = load_workbook(excel_file, data_only=True)
    
    print(f"Available sheets: {wb.sheetnames}\n")
    print("=" * 100)
    
    for sheet_name in wb.sheetnames:
        ws = wb[sheet_name]
        print(f"\n{'=' * 100}")
        print(f"SHEET: {sheet_name}")
        print(f"{'=' * 100}\n")
        
        # Print all rows
        for row_idx, row in enumerate(ws.iter_rows(values_only=True), start=1):
            row_data = [str(cell) if cell is not None else "" for cell in row]
            print(f"Row {row_idx:3d}: {' | '.join(row_data)}")
        
        print(f"\nTotal rows in {sheet_name}: {ws.max_row}")
        print(f"Total columns in {sheet_name}: {ws.max_column}")
    
    wb.close()
    print(f"\n{'=' * 100}")
    print("Excel file read successfully")
    
except ImportError:
    print("Error: openpyxl not installed")
    print("Install with: pip install openpyxl")
    sys.exit(1)
except Exception as e:
    print(f"Error reading Excel file: {e}")
    import traceback
    traceback.print_exc()
    sys.exit(1)
EOF
