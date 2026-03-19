#!/usr/bin/env python3
"""
Comprehensive Journal Validation Script - STRICT LINE ACCOUNTING
Every line must be either:
1. PROCESSED (understood and will be migrated)
2. SKIPPED (understood but intentionally not migrated)
3. FLAGGED (unexplained - needs investigation)

Usage: python3 journal_validation.py
"""

import os
import re
import glob
from collections import Counter, defaultdict
from datetime import datetime

# Line categories
PROCESSED = "PROCESSED"  # Data that will be migrated
SKIPPED = "SKIPPED"      # Understood but not migrated (e.g., SNF header, collapsed::)
FLAGGED = "FLAGGED"      # Unexplained - needs investigation

class LineAccountant:
    """Tracks every line and ensures nothing is missed."""
    
    def __init__(self):
        self.total_lines = 0
        self.processed_lines = 0
        self.skipped_lines = 0
        self.flagged_lines = 0
        
        # Detailed tracking
        self.processed_details = defaultdict(list)  # category -> [(file, line_num, content)]
        self.skipped_details = defaultdict(list)
        self.flagged_details = defaultdict(list)
        
        # Data extraction
        self.tickers = []
        self.images = []
        self.tags = {'trade': [], 'reason': [], 'management': [], 'other': []}
        self.notes = {'code_block': [], 'simple': [], 'plan': []}
        self.important_tags = []  # #important tags
        
    def process_line(self, file_name, line_num, line, category, subcategory, status):
        """Record a line's status."""
        self.total_lines += 1
        
        entry = {
            'file': file_name,
            'line': line_num,
            'content': line.rstrip()[:100],  # Truncate for display
            'full_content': line.rstrip()
        }
        
        if status == PROCESSED:
            self.processed_lines += 1
            self.processed_details[subcategory].append(entry)
        elif status == SKIPPED:
            self.skipped_lines += 1
            self.skipped_details[subcategory].append(entry)
        else:  # FLAGGED
            self.flagged_lines += 1
            self.flagged_details[subcategory].append(entry)
    
    def verify_totals(self):
        """Verify all lines are accounted for."""
        accounted = self.processed_lines + self.skipped_lines + self.flagged_lines
        return accounted == self.total_lines, accounted, self.total_lines


def analyze_journal_line(line, line_stripped, in_code_block, prev_line_type, file_name, line_num, accountant):
    """
    Analyze a single line and categorize it.
    Returns: (status, category, new_in_code_block, line_type)
    """
    
    # Patterns
    image_pattern = r'!\[.*?\]\(([^)]+)\)'
    ticker_pattern = r'\|\s*`([A-Z0-9_!]+)`\s*\|'
    tag_pattern = r'#([tmr])\.([a-zA-Z0-9_-]+)'
    important_pattern = r'#important'
    
    # === CODE BLOCK HANDLING ===
    if '```' in line:
        return PROCESSED, 'code_block_marker', not in_code_block, 'code_block_marker'
    
    if in_code_block:
        # Check for Plan: inside code blocks
        if 'Plan:' in line:
            accountant.notes['plan'].append({
                'file': file_name, 'line': line_num, 'content': line_stripped,
                'type': 'plan_in_code_block'
            })
            return PROCESSED, 'plan_note_in_code_block', True, 'plan_note'
        return PROCESSED, 'code_block_content', True, 'code_block_content'
    
    # === EMPTY LINES ===
    if not line_stripped:
        return SKIPPED, 'empty_line', False, 'empty'
    
    # === SNF HEADER ===
    if 'SNF' in line and 'Journal' in line:
        return SKIPPED, 'snf_header', False, 'snf_header'
    
    # === COLLAPSED LINES ===
    if 'collapsed::' in line:
        return SKIPPED, 'collapsed_marker', False, 'collapsed'
    
    # === JOURNAL ENTRY ROW (ticker + tags) ===
    ticker_match = re.search(ticker_pattern, line)
    if ticker_match and '|' in line:
        ticker = ticker_match.group(1)
        accountant.tickers.append({'file': file_name, 'line': line_num, 'ticker': ticker})
        
        # Extract all tags from this line
        tags = re.findall(tag_pattern, line)
        for tag_type, tag_value in tags:
            if tag_type == 't':
                accountant.tags['trade'].append({'file': file_name, 'line': line_num, 'tag': f't.{tag_value}'})
            elif tag_type == 'r':
                accountant.tags['reason'].append({'file': file_name, 'line': line_num, 'tag': f'r.{tag_value}'})
            elif tag_type == 'm':
                accountant.tags['management'].append({'file': file_name, 'line': line_num, 'tag': f'm.{tag_value}'})
        
        # Check for #important tag
        if re.search(important_pattern, line):
            accountant.important_tags.append({'file': file_name, 'line': line_num, 'content': line_stripped})
        
        # Check for any OTHER tags not matching #t./#r./#m. pattern
        # Use negative lookahead to exclude #t.xxx, #r.xxx, #m.xxx patterns
        other_tag_pattern = r'#([a-zA-Z][a-zA-Z0-9_-]*)(?![.\w])'
        other_tags = re.findall(other_tag_pattern, line)
        for tag in other_tags:
            # Skip known non-data tags and the base t/r/m (which are part of #t.xxx patterns)
            if tag not in ['trading-tome', 't', 'r', 'm'] and not tag.startswith('t.') and not tag.startswith('r.') and not tag.startswith('m.'):
                accountant.tags['other'].append({'file': file_name, 'line': line_num, 'tag': f'#{tag}'})
        
        return PROCESSED, 'journal_entry_row', False, 'journal_row'
    
    # === IMAGE LINES ===
    image_match = re.search(image_pattern, line)
    if image_match:
        image_path = image_match.group(1)
        accountant.images.append({'file': file_name, 'line': line_num, 'path': image_path})
        
        # Check if there's text BEFORE or AFTER the image on the same line
        line_without_image = re.sub(image_pattern, '', line).strip()
        line_without_image = re.sub(r'^-\s*', '', line_without_image).strip()  # Remove leading -
        # Remove Logseq image dimension metadata {:height X, :width Y}
        line_without_image = re.sub(r'\{:height\s+\d+,?\s*:width\s+\d+\}', '', line_without_image).strip()
        
        if line_without_image:
            # There's additional content on the image line that's NOT just metadata!
            return FLAGGED, 'image_line_with_extra_content', False, 'image_with_extra'
        
        return PROCESSED, 'image_line', False, 'image'
    
    # === SIMPLE NOTES (lines starting with - but not images) ===
    if line_stripped.startswith('-'):
        content = line_stripped[1:].strip()
        
        # Check if it's just a dash with no content
        if not content:
            return SKIPPED, 'empty_dash_line', False, 'empty_dash'
        
        # Check for background-color:: (Logseq property)
        if 'background-color::' in content:
            return SKIPPED, 'logseq_property', False, 'logseq_property'
        
        # This is a SIMPLE NOTE - needs to be migrated!
        accountant.notes['simple'].append({
            'file': file_name, 
            'line': line_num, 
            'content': content,
            'prev_line_type': prev_line_type
        })
        return PROCESSED, 'simple_note', False, 'simple_note'
    
    # === LOGSEQ PROPERTIES ===
    if '::' in line and not line_stripped.startswith('|'):
        # Logseq property like background-color:: yellow
        return SKIPPED, 'logseq_property', False, 'logseq_property'
    
    # === ANYTHING ELSE IS FLAGGED ===
    return FLAGGED, 'unexplained_line', False, 'unexplained'


def validate_journals():
    """Main validation function."""
    
    print("=" * 80)
    print("COMPREHENSIVE JOURNAL VALIDATION - STRICT LINE ACCOUNTING")
    print("=" * 80)
    print()
    
    # Get all markdown files
    pattern = "/home/aman/Projects/go-fun/processed/*.md"
    files = sorted(glob.glob(pattern))
    
    accountant = LineAccountant()
    
    print(f"Scanning {len(files)} files...")
    print()
    
    # Process each file
    for file_path in files:
        file_name = os.path.basename(file_path)
        
        with open(file_path, 'r', encoding='utf-8') as f:
            lines = f.readlines()
        
        in_code_block = False
        prev_line_type = None
        
        for line_num, line in enumerate(lines, 1):
            line_stripped = line.strip()
            
            status, category, in_code_block, line_type = analyze_journal_line(
                line, line_stripped, in_code_block, prev_line_type,
                file_name, line_num, accountant
            )
            
            accountant.process_line(file_name, line_num, line, category, category, status)
            prev_line_type = line_type
    
    # === VERIFICATION ===
    print("=" * 80)
    print("LINE ACCOUNTING VERIFICATION")
    print("=" * 80)
    
    is_valid, accounted, total = accountant.verify_totals()
    
    print(f"\nTotal lines scanned: {total}")
    print(f"  PROCESSED (will migrate): {accountant.processed_lines}")
    print(f"  SKIPPED (understood, not migrated): {accountant.skipped_lines}")
    print(f"  FLAGGED (needs investigation): {accountant.flagged_lines}")
    print(f"  Accounted total: {accounted}")
    print(f"\nVerification: {'✓ ALL LINES ACCOUNTED' if is_valid else '✗ MISMATCH - LINES MISSING'}")
    
    # === PROCESSED BREAKDOWN ===
    print("\n" + "=" * 80)
    print("PROCESSED DATA (Will Be Migrated)")
    print("=" * 80)
    
    print("\n--- Processed Line Categories ---")
    for category, entries in sorted(accountant.processed_details.items()):
        print(f"  {category}: {len(entries)}")
    
    # === DATA EXTRACTION SUMMARY ===
    print("\n--- Extracted Data Summary ---")
    print(f"  Tickers: {len(accountant.tickers)}")
    print(f"  Images: {len(accountant.images)}")
    print(f"  Trade Tags (#t.*): {len(accountant.tags['trade'])}")
    print(f"  Reason Tags (#r.*): {len(accountant.tags['reason'])}")
    print(f"  Management Tags (#m.*): {len(accountant.tags['management'])}")
    print(f"  Other Tags: {len(accountant.tags['other'])}")
    print(f"  Important Tags (#important): {len(accountant.important_tags)}")
    print(f"  Code Block Notes: {len(accountant.notes['code_block'])}")
    print(f"  Plan Notes: {len(accountant.notes['plan'])}")
    print(f"  Simple Notes (CRITICAL): {len(accountant.notes['simple'])}")
    
    # === SIMPLE NOTES DETAIL (CRITICAL) ===
    if accountant.notes['simple']:
        print("\n" + "=" * 80)
        print("SIMPLE NOTES - MUST BE MIGRATED (Outside Code Blocks)")
        print("=" * 80)
        print(f"\nFound {len(accountant.notes['simple'])} simple notes that need migration:")
        for i, note in enumerate(accountant.notes['simple'], 1):
            print(f"\n  {i}. {note['file']}:{note['line']}")
            print(f"     Content: {note['content'][:80]}{'...' if len(note['content']) > 80 else ''}")
            print(f"     Previous line type: {note['prev_line_type']}")
    
    # === IMPORTANT TAGS ===
    if accountant.important_tags:
        print("\n" + "=" * 80)
        print("#IMPORTANT TAGS - MUST BE CAPTURED")
        print("=" * 80)
        print(f"\nFound {len(accountant.important_tags)} entries with #important tag:")
        for i, entry in enumerate(accountant.important_tags, 1):
            print(f"  {i}. {entry['file']}:{entry['line']}: {entry['content'][:60]}...")
    
    # === OTHER TAGS (Non-standard) ===
    if accountant.tags['other']:
        print("\n" + "=" * 80)
        print("OTHER TAGS (Non-standard #t./#r./#m. tags)")
        print("=" * 80)
        other_tag_counter = Counter([t['tag'] for t in accountant.tags['other']])
        print(f"\nFound {len(accountant.tags['other'])} non-standard tags:")
        for tag, count in other_tag_counter.most_common():
            print(f"  {tag}: {count}")
    
    # === SKIPPED BREAKDOWN ===
    print("\n" + "=" * 80)
    print("SKIPPED DATA (Understood, Not Migrated)")
    print("=" * 80)
    
    for category, entries in sorted(accountant.skipped_details.items()):
        print(f"  {category}: {len(entries)}")
    
    # === FLAGGED (CRITICAL) ===
    print("\n" + "=" * 80)
    print("⚠️  FLAGGED LINES - NEEDS INVESTIGATION ⚠️")
    print("=" * 80)
    
    if accountant.flagged_lines == 0:
        print("\n✓ No flagged lines - all content is accounted for!")
    else:
        print(f"\n⚠️  {accountant.flagged_lines} LINES NEED INVESTIGATION:")
        
        for category, entries in sorted(accountant.flagged_details.items()):
            print(f"\n--- {category} ({len(entries)} lines) ---")
            for entry in entries[:20]:  # Show first 20
                print(f"  {entry['file']}:{entry['line']}: {entry['content']}")
            if len(entries) > 20:
                print(f"  ... and {len(entries) - 20} more")
    
    # === TAG ANALYSIS (PRD 4.8.6.3) ===
    print("\n" + "=" * 80)
    print("TAG ANALYSIS (PRD 4.8.6.3 Compliance)")
    print("=" * 80)
    
    print("\n--- Trade Tags (#t.*) ---")
    trade_counter = Counter([t['tag'] for t in accountant.tags['trade']])
    for tag, count in trade_counter.most_common():
        print(f"  {tag}: {count}")
    
    print("\n--- Reason Tags (#r.*) ---")
    reason_counter = Counter([t['tag'] for t in accountant.tags['reason']])
    for tag, count in reason_counter.most_common(20):
        print(f"  {tag}: {count}")
    if len(reason_counter) > 20:
        print(f"  ... and {len(reason_counter) - 20} more unique tags")
    
    print("\n--- Management Tags (#m.*) ---")
    mgmt_counter = Counter([t['tag'] for t in accountant.tags['management']])
    for tag, count in mgmt_counter.most_common():
        print(f"  {tag}: {count}")
    
    # === FINAL SUMMARY ===
    print("\n" + "=" * 80)
    print("FINAL VALIDATION SUMMARY")
    print("=" * 80)
    
    issues = []
    
    if accountant.flagged_lines > 0:
        issues.append(f"⚠️  {accountant.flagged_lines} unexplained lines need investigation")
    
    if accountant.notes['simple']:
        issues.append(f"⚠️  {len(accountant.notes['simple'])} simple notes must be migrated (currently may be lost)")
    
    if accountant.important_tags:
        issues.append(f"ℹ️  {len(accountant.important_tags)} #important tags found - ensure they're captured")
    
    if accountant.tags['other']:
        issues.append(f"ℹ️  {len(accountant.tags['other'])} non-standard tags found - verify handling")
    
    if issues:
        print("\nISSUES FOUND:")
        for issue in issues:
            print(f"  {issue}")
    else:
        print("\n✓ All validation checks passed!")
    
    print("\n--- Migration Counts (for comparison with migration script) ---")
    print(f"  Files: {len(files)}")
    print(f"  Journal Entries (tickers): {len(accountant.tickers)}")
    print(f"  Images: {len(accountant.images)}")
    print(f"  Total Tags: {len(accountant.tags['trade']) + len(accountant.tags['reason']) + len(accountant.tags['management'])}")
    print(f"  Notes (code block + simple + plan): {len(accountant.notes['code_block']) + len(accountant.notes['simple']) + len(accountant.notes['plan'])}")
    
    return accountant


if __name__ == "__main__":
    validate_journals()
