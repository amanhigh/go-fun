#!/usr/bin/env python3
"""
Comprehensive Journal Analysis Script
Analyzes all aspects of journal data in processed directory
Usage: python3 journal_analysis.py
"""

import os
import re
import glob
from collections import Counter
from datetime import datetime

def analyze_journal_data():
    """Perform comprehensive analysis of journal data."""
    
    print("=== Comprehensive Journal Analysis ===")
    print("Scanning directory: /home/aman/Projects/go-fun/processed/*.md")
    print()
    
    # Get all markdown files
    pattern = "/home/aman/Projects/go-fun/processed/*.md"
    files = glob.glob(pattern)
    
    # Basic file statistics
    total_files = len(files)
    total_lines = 0
    all_content = []
    
    for file in files:
        with open(file, 'r', encoding='utf-8') as f:
            content = f.readlines()
            total_lines += len(content)
            all_content.extend(content)
    
    print("=== File Statistics ===")
    print(f"Total journal files: {total_files}")
    print(f"Total lines across all files: {total_lines}")
    print()
    
    # Content analysis
    print("=== Content Analysis ===")
    
    # Count images (lines containing ![[...]] )
    image_pattern = r'!\[.*?\]'
    image_count = sum(1 for line in all_content if re.search(image_pattern, line))
    print(f"Images (![[...]]): {image_count}")
    
    # Count notes (lines starting with - but not SNF)
    note_count = sum(1 for line in all_content 
                   if line.startswith('- ') and 'SNF' not in line)
    print(f"Notes (- ...): {note_count}")
    
    # Count journal rows (lines containing |)
    journal_row_count = sum(1 for line in all_content if '|' in line)
    print(f"Journal rows (| ...): {journal_row_count}")
    
    # Extract tickers from journal rows
    ticker_pattern = r'\|\s*`([^`]+)`'
    tickers = set()
    for line in all_content:
        match = re.search(ticker_pattern, line)
        if match:
            tickers.add(match.group(1))
    
    ticker_count = len(tickers)
    print(f"Unique tickers: {ticker_count}")
    
    # Count SNF rows (exclude from analysis)
    snf_count = sum(1 for line in all_content if 'SNF' in line)
    print(f"SNF rows (to exclude): {snf_count}")
    print()
    
    # Legacy tags analysis
    print("=== Legacy Tags Analysis ===")
    
    # Extract all tags
    tag_pattern = r'#([tmr])\.([a-zA-Z0-9_-]+)'
    all_tags = []
    
    for line in all_content:
        matches = re.findall(tag_pattern, line)
        for tag_type, tag_value in matches:
            all_tags.append(f"{tag_type}.{tag_value}")
    
    # Separate by tag type
    trade_tags = [tag for tag in all_tags if tag.startswith('t.')]
    reason_tags = [tag for tag in all_tags if tag.startswith('r.')]
    management_tags = [tag for tag in all_tags if tag.startswith('m.')]
    
    trade_count = len(trade_tags)
    reason_count = len(reason_tags)
    management_count = len(management_tags)
    total_tag_count = len(all_tags)
    
    print(f"Trade Tags (#t.*): {trade_count} occurrences ({len(set(trade_tags))} unique)")
    print(f"Reason Tags (#r.*): {reason_count} occurrences ({len(set(reason_tags))} unique)")
    print(f"Management Tags (#m.*): {management_count} occurrences ({len(set(management_tags))} unique)")
    print(f"Total Tags: {total_tag_count} occurrences")
    print()
    
    # Detailed tag analysis per PRD 4.8.6.3
    print("=== Detailed Tag Analysis (PRD 4.8.6.3) ===")
    
    # Trade tags by frequency - ALL of them for migration verification
    trade_counter = Counter(trade_tags)
    print("Trade Tags (#t.*) - Complete List:")
    for tag, count in trade_counter.most_common():
        print(f"   {count:4d} {tag}")
    print(f"   Total: {trade_count} occurrences, {len(trade_counter)} unique")
    print()
    
    # ALL reason tags for migration verification
    reason_counter = Counter(reason_tags)
    print("Reason Tags (#r.*) - Complete List:")
    for tag, count in reason_counter.most_common():
        print(f"    {count:3d} {tag}")
    print(f"   Total: {reason_count} occurrences, {len(reason_counter)} unique")
    print()
    
    # ALL management tags
    management_counter = Counter(management_tags)
    print("Management Tags (#m.*) - Complete List:")
    for tag, count in management_counter.most_common():
        print(f"    {count:3d} {tag}")
    print(f"   Total: {management_count} occurrences, {len(management_counter)} unique")
    print()
    
    # Tag mapping verification for migration
    print("=== Tag Mapping Verification (PRD 4.8.6.3) ===")
    print("\nTrade Tag Mappings:")
    trade_mappings = {
        't.trend': 'tags -> tag: trend, type: DIRECTION',
        't.ctrend': 'tags -> tag: ctrend, type: DIRECTION',
        't.mwd': 'sequence -> MWD',
        't.yr': 'sequence -> YR',
        't.wdh': 'sequence -> WDH',
        't.rejected': 'type -> REJECTED',
        't.set': 'type -> SET',
        't.taken': 'type -> TAKEN',
        't.fail': 'status -> FAIL',
        't.success': 'status -> SUCCESS',
        't.broken': 'status -> BROKEN',
        't.miss': 'status -> MISSED',
        't.justloss': 'status -> JUST_LOSS',
        't.running': 'status -> RUNNING',
        't.full': 'tags -> tag: double, type: MANAGEMENT',
    }
    for tag, count in trade_counter.most_common():
        mapping = trade_mappings.get(tag, 'UNKNOWN - needs mapping')
        print(f"   {count:4d} {tag} -> {mapping}")
    print()
    
    print("Reason Tag Mappings (all map to tags with type: REASON):")
    for tag, count in reason_counter.most_common():
        # Check for override pattern (e.g., r.dep-loc)
        tag_value = tag[2:]  # Remove 'r.' prefix
        if '-' in tag_value:
            parts = tag_value.split('-', 1)
            print(f"    {count:3d} {tag} -> tag: {parts[0]}, override: {parts[1]}")
        else:
            print(f"    {count:3d} {tag} -> tag: {tag_value}")
    print()
    
    print("Management Tag Mappings (all map to tags with type: MANAGEMENT):")
    for tag, count in management_counter.most_common():
        tag_value = tag[2:]  # Remove 'm.' prefix
        print(f"    {count:3d} {tag} -> tag: {tag_value}")
    print()
    
    # Date range analysis
    print("=== Date Range Analysis ===")
    
    # Extract dates from filenames
    dates = []
    date_pattern = r'(\d{4}_\d{2}_\d{2})'
    
    for file in files:
        match = re.search(date_pattern, file)
        if match:
            dates.append(match.group(1))
    
    dates.sort()
    earliest_date = dates[0] if dates else "N/A"
    latest_date = dates[-1] if dates else "N/A"
    date_count = len(dates)
    
    # Format dates for better readability
    def format_date(date_str):
        if date_str == "N/A":
            return date_str
        year, month, day = date_str.split('_')
        return f"{year}-{month}-{day}"
    
    print(f"Date range: {format_date(earliest_date)} to {format_date(latest_date)}")
    print(f"Total dates with entries: {date_count}")
    print()
    
    # Comprehensive ticker analysis
    print("=== Ticker Analysis ===")
    
    # Count ticker occurrences
    ticker_occurrences = []
    for line in all_content:
        match = re.search(ticker_pattern, line)
        if match:
            ticker_occurrences.append(match.group(1))
    
    ticker_counter = Counter(ticker_occurrences)
    
    print(f"Total ticker occurrences: {len(ticker_occurrences)}")
    print(f"Unique tickers: {len(ticker_counter)}")
    print()
    
    print("Top 20 Most Active Tickers:")
    for ticker, count in ticker_counter.most_common(20):
        print(f"    {count:3d} {ticker}")
    print()
    
    # Ticker frequency distribution
    print("Ticker Frequency Distribution:")
    freq_dist = Counter(ticker_counter.values())
    for freq, count in sorted(freq_dist.items(), reverse=True):
        print(f"    {count:3d} tickers appear {freq} times")
    print()
    
    # Exclusion patterns analysis
    print("=== Exclusion Patterns Analysis ===")
    
    # SNF pattern analysis
    snf_lines = [line for line in all_content if 'SNF' in line]
    print(f"SNF pattern occurrences: {len(snf_lines)}")
    
    # Analyze SNF line types
    snf_headers = [line for line in snf_lines if line.strip().startswith('- SNF')]
    snf_collapsed = [line for line in snf_lines if 'collapsed:: true' in line]
    snf_other = len(snf_lines) - len(snf_headers) - len(snf_collapsed)
    
    print("SNF line breakdown:")
    print(f"  SNF headers: {len(snf_headers)}")
    print(f"  SNF collapsed lines: {len(snf_collapsed)}")
    print(f"  Other SNF lines: {snf_other}")
    print()
    
    # Other potential exclusion patterns
    exclusion_patterns = {
        'collapsed:: true': r'collapsed::\s*true',
        'TODO lines': r'TODO',
        'NOTE lines': r'NOTE',
        'DEBUG lines': r'DEBUG',
        'Empty lines': r'^\s*$'
    }
    
    print("Other Exclusion Patterns:")
    for pattern_name, pattern_regex in exclusion_patterns.items():
        count = sum(1 for line in all_content if re.search(pattern_regex, line))
        print(f"  {pattern_name}: {count}")
    print()
    
    # Content distribution
    print("=== Content Distribution ===")
    
    # Count rows with each tag type
    rows_with_trade_tags = sum(1 for line in all_content if re.search(r'#t\.', line))
    rows_with_reason_tags = sum(1 for line in all_content if re.search(r'#r\.', line))
    rows_with_management_tags = sum(1 for line in all_content if re.search(r'#m\.', line))
    
    print("Journal rows per tag type:")
    print(f"  Rows with trade tags: {rows_with_trade_tags}")
    print(f"  Rows with reason tags: {rows_with_reason_tags}")
    print(f"  Rows with management tags: {rows_with_management_tags}")
    print()
    
    print("Content density:")
    print(f"  Average lines per file: {total_lines // total_files if total_files > 0 else 0}")
    print(f"  Average images per file: {image_count // total_files if total_files > 0 else 0}")
    print(f"  Average notes per file: {note_count // total_files if total_files > 0 else 0}")
    print(f"  Average journal rows per file: {journal_row_count // total_files if total_files > 0 else 0}")
    print()
    
    # Migration summary
    print("=== Migration Summary ===")
    print(f"Total legacy tag occurrences to migrate: {total_tag_count}")
    print(f"Total images to migrate: {image_count}")
    print(f"Total notes to migrate: {note_count}")
    print(f"Total journal rows to migrate: {journal_row_count}")
    print(f"Total unique tickers to migrate: {ticker_count}")
    print(f"SNF rows to exclude: {snf_count}")
    print()
    
    # Comprehensive line categorization
    print("=== Comprehensive Line Categorization ===")
    
    # Count different line types
    line_types = {
        'Image lines': 0,
        'Plan notes': 0,
        'Post-set notes': 0,
        'Other notes': 0,
        'Journal rows': 0,
        'SNF lines': 0,
        'Collapsed lines': 0,
        'Code block markers': 0,  # Just the ``` lines
        'Code block content': 0,  # Content inside code blocks (part of notes)
        'Empty lines': 0,
        'Other lines': 0
    }
    
    # Collect samples for analysis with detailed file/line info
    post_set_samples = []  # (filename, line_num, content, context)
    plan_note_samples = []
    other_note_samples = []
    code_block_samples = []
    other_line_samples = []
    
    # Track potential notes after images (user's concern)
    potential_post_image_notes = []
    
    # Process files and track context
    file_contents = {}
    for file_path in files:
        with open(file_path, 'r', encoding='utf-8') as f:
            file_contents[file_path] = f.readlines()
    
    for file_path, lines in file_contents.items():
        current_file = file_path.split('/')[-1]  # Get just filename
        prev_was_image = False
        in_code_block = False  # Track if we're inside a code block
        
        for i, line in enumerate(lines):
            line_stripped = line.strip()
            
            # Check for code block markers first
            if '```' in line:
                line_types['Code block markers'] += 1
                in_code_block = not in_code_block
                if len(code_block_samples) < 10:
                    code_block_samples.append({
                        'file': current_file,
                        'line': i + 1,
                        'content': line.strip()
                    })
                continue
            
            # If inside code block, count as code block content (part of notes)
            if in_code_block:
                line_types['Code block content'] += 1
                # Also check for Plan: inside code blocks
                if 'Plan:' in line:
                    line_types['Plan notes'] += 1
                    if len(plan_note_samples) < 10:
                        plan_note_samples.append({
                            'file': current_file,
                            'line': i + 1,
                            'content': line.strip()
                        })
                continue
            
            # Check if this line could be a post-image note
            if prev_was_image and line_stripped:
                # Not an image, not a table row, not collapsed, not code block
                if (not re.search(image_pattern, line) and 
                    not line_stripped.startswith('|') and 
                    'collapsed::' not in line and
                    '```' not in line and
                    'SNF' not in line and
                    not line_stripped.startswith('-')):
                    potential_post_image_notes.append({
                        'file': current_file,
                        'line': i + 1,
                        'content': line_stripped,
                        'prev_line': lines[i-1].strip() if i > 0 else ''
                    })
            
            prev_was_image = bool(re.search(image_pattern, line))
            
            if not line_stripped:
                line_types['Empty lines'] += 1
            elif 'SNF' in line:
                line_types['SNF lines'] += 1
            elif re.search(image_pattern, line):
                line_types['Image lines'] += 1
                # Check for post-set notes after images
                if i + 1 < len(lines):
                    next_line = lines[i + 1].strip()
                    if next_line and not next_line.startswith('-') and not next_line.startswith('|') and not re.search(image_pattern, next_line) and 'SNF' not in next_line and 'collapsed::' not in next_line and '```' not in next_line:
                        line_types['Post-set notes'] += 1
                        if len(post_set_samples) < 10:
                            post_set_samples.append({
                                'file': current_file,
                                'line': i + 2,  # next line
                                'content': next_line,
                                'image_line': line.strip()
                            })
            elif 'Plan:' in line:
                line_types['Plan notes'] += 1
                if len(plan_note_samples) < 10:
                    plan_note_samples.append({
                        'file': current_file,
                        'line': i + 1,
                        'content': line.strip()
                    })
            elif line.startswith('- ') and 'SNF' not in line:
                line_types['Other notes'] += 1
                if len(other_note_samples) < 10:
                    other_note_samples.append({
                        'file': current_file,
                        'line': i + 1,
                        'content': line.strip()
                    })
            elif 'collapsed::' in line:
                line_types['Collapsed lines'] += 1
            elif '|' in line:
                line_types['Journal rows'] += 1
            else:
                line_types['Other lines'] += 1
                if len(other_line_samples) < 30:
                    other_line_samples.append({
                        'file': current_file,
                        'line': i + 1,
                        'content': line.rstrip()
                    })
    
    print("Line type distribution:")
    for line_type, count in line_types.items():
        print(f"  {line_type}: {count}")
    
    # Calculate total properly - code block content is part of notes
    total_excluding_content = sum(line_types.values()) - line_types['Code block content']
    print(f"\nLine total verification: {total_excluding_content} + {line_types['Code block content']} (code block content) = {total_excluding_content + line_types['Code block content']} (should equal {total_lines})")
    print(f"Match: {'✓' if total_excluding_content + line_types['Code block content'] == total_lines else '✗ MISMATCH'}")
    print()
    
    # Update notes total to include code block content (which are part of notes)
    total_notes = line_types['Plan notes'] + line_types['Post-set notes'] + line_types['Other notes']
    print("=== Updated Notes Analysis ===")
    print(f"  Plan notes: {line_types['Plan notes']}")
    print(f"  Post-set notes: {line_types['Post-set notes']}")
    print(f"  Other notes: {line_types['Other notes']}")
    print(f"  Code block content (part of notes): {line_types['Code block content']}")
    print(f"  Total notes (including code block content): {total_notes + line_types['Code block content']}")
    print()
    
    # Show samples of different note types with file/line info
    print("=== Note Type Analysis (with file:line) ===")
    
    print("Plan note samples:")
    for sample in plan_note_samples:
        print(f"  {sample['file']}:{sample['line']}: {sample['content']}")
    print()
    
    if post_set_samples:
        print(f"Post-set note samples (notes after images) - {len(post_set_samples)} shown of {line_types['Post-set notes']} total:")
        for sample in post_set_samples:
            print(f"  {sample['file']}:{sample['line']}: {sample['content']}")
            print(f"    (after image: {sample['image_line'][:60]}...)")
        print()
    
    if other_note_samples:
        print(f"Other note samples (- lines) - {len(other_note_samples)} shown of {line_types['Other notes']} total:")
        for sample in other_note_samples:
            print(f"  {sample['file']}:{sample['line']}: {sample['content']}")
        print()
    
    # Show potential post-image notes that might be missed
    if potential_post_image_notes:
        print(f"=== POTENTIAL POST-IMAGE NOTES (User Concern) ===")
        print(f"Found {len(potential_post_image_notes)} lines that appear after images and might be notes:")
        for i, note in enumerate(potential_post_image_notes[:20]):
            print(f"  {i+1}. {note['file']}:{note['line']}: {note['content']}")
        if len(potential_post_image_notes) > 20:
            print(f"  ... and {len(potential_post_image_notes) - 20} more")
        print()
    
    # Show uncategorized content
    print("=== Uncategorized Content Analysis ===")
    
    if code_block_samples:
        print("Code block samples (with file:line):")
        for i, sample in enumerate(code_block_samples, 1):
            print(f"  {i}. {sample['file']}:{sample['line']}: {sample['content']}")
        print()
    
    if other_line_samples:
        print(f"=== OTHER/UNEXPLAINED LINES (showing {len(other_line_samples)} of {line_types['Other lines']} total) ===")
        print("These lines are NOT categorized and need review:")
        for i, sample in enumerate(other_line_samples, 1):
            print(f"  {i:2d}. {sample['file']}:{sample['line']}: {sample['content']}")
        print()
        
        # Show more context around other lines
        print("=== Other Line Context Analysis ===")
        print("Showing 3 lines before and after each 'other line' for context:")
        print()
        
        context_samples = []
        for file_path, lines in file_contents.items():
            filename = file_path.split('/')[-1]
            for i, line in enumerate(lines):
                line_stripped = line.strip()
                if (line_stripped and 
                    not any(['SNF' in line,
                           re.search(image_pattern, line),
                           'Plan:' in line,
                           line.startswith('- '),
                           '|' in line,
                           'collapsed::' in line,
                           '```' in line])):
                    
                    # Get context (3 lines before and after)
                    start_idx = max(0, i - 3)
                    end_idx = min(len(lines), i + 4)
                    context = lines[start_idx:end_idx]
                    
                    context_str = "".join(context)
                    if len(context_samples) < 5 and context_str not in [s[1] for s in context_samples]:
                        context_samples.append((filename, i+1, context_str))
        
        for filename, line_num, context in context_samples:
            print(f"File: {filename}, Line {line_num}:")
            print("Context:")
            for j, ctx_line in enumerate(context.split('\n')):
                marker = " >>> " if j == 3 else "     "  # Mark the target line
                print(f"{marker}{ctx_line}")
            print()
        
        # Analyze patterns in other lines
        print("Other line pattern analysis:")
        patterns = {
            'Lines starting with spaces': 0,
            'Lines with special chars': 0,
            'Lines with numbers only': 0,
            'Lines with URLs': 0,
            'Lines with dates': 0,
            'Lines with brackets': 0,
            'Lines with arrows': 0,
            'Lines with colons': 0,
            'Lines with dashes': 0,
            'Other patterns': 0
        }
        
        for file_path, lines in file_contents.items():
            for line in lines:
                if line.strip() and not any([
                    'SNF' in line,
                    re.search(image_pattern, line),
                    'Plan:' in line,
                    line.startswith('- '),
                    '|' in line,
                    'collapsed::' in line,
                    '```' in line
                ]):
                    if line.startswith('    ') or line.startswith('\t'):
                        patterns['Lines starting with spaces'] += 1
                    elif re.search(r'[<>*+={}~|]', line):
                        patterns['Lines with special chars'] += 1
                    elif re.search(r'^\s*\d+\s*$', line):
                        patterns['Lines with numbers only'] += 1
                    elif re.search(r'http[s]?://', line):
                        patterns['Lines with URLs'] += 1
                    elif re.search(r'\d{4}[-_]\d{2}[-_]\d{2}', line):
                        patterns['Lines with dates'] += 1
                    elif re.search(r'[\[\]{}()]', line):
                        patterns['Lines with brackets'] += 1
                    elif re.search(r'[→←↑↓]', line):
                        patterns['Lines with arrows'] += 1
                    elif ':' in line:
                        patterns['Lines with colons'] += 1
                    elif '-' in line and not line.startswith('-'):
                        patterns['Lines with dashes'] += 1
                    else:
                        patterns['Other patterns'] += 1
        
        for pattern, count in patterns.items():
            if count > 0:
                print(f"  {pattern}: {count}")
        print()
    
    # Detailed notes breakdown
    total_notes = line_types['Plan notes'] + line_types['Post-set notes'] + line_types['Other notes']
    print("Notes Summary:")
    print(f"  Plan notes: {line_types['Plan notes']}")
    print(f"  Post-set notes: {line_types['Post-set notes']}")
    print(f"  Other notes: {line_types['Other notes']}")
    print(f"  Total notes: {total_notes}")
    print()
    
    # Migration strategy validation
    print("=== Migration Strategy Validation ===")
    
    # Calculate migration counts
    migrate_content = line_types['Image lines'] + line_types['Plan notes'] + line_types['Post-set notes'] + line_types['Other notes'] + line_types['Journal rows']
    exclude_content = line_types['SNF lines'] + line_types['Collapsed lines'] + line_types['Empty lines']
    review_content = line_types['Code blocks'] + line_types['Other lines']
    
    print("Migration Counts:")
    print(f"  Images to migrate: {line_types['Image lines']}")
    print(f"  Plan notes to migrate: {line_types['Plan notes']}")
    print(f"  Post-set notes to migrate: {line_types['Post-set notes']}")
    print(f"  Other notes to migrate: {line_types['Other notes']}")
    print(f"  Journal rows to migrate: {line_types['Journal rows']}")
    print(f"  Total to migrate: {migrate_content}")
    print()
    
    print("Exclusion Counts:")
    print(f"  SNF lines to exclude: {line_types['SNF lines']}")
    print(f"  Collapsed lines to exclude: {line_types['Collapsed lines']}")
    print(f"  Empty lines to exclude: {line_types['Empty lines']}")
    print(f"  Total to exclude: {exclude_content}")
    print()
    
    print("Review Counts:")
    print(f"  Code blocks to review: {line_types['Code blocks']}")
    print(f"  Other lines to review: {line_types['Other lines']}")
    print(f"  Total to review: {review_content}")
    print()
    
    # Validation against total
    total_calculated = migrate_content + exclude_content + review_content
    print("Validation:")
    print(f"  Calculated total: {total_calculated}")
    print(f"  Actual total lines: {total_lines}")
    print(f"  Match: {'✓' if total_calculated == total_lines else '✗ MISMATCH'}")
    print()
    
    # PRD validation counts
    print("=== PRD Validation Counts ===")
    print("These counts should match section 4.8.6 in the PRD:")
    print()
    
    print("Journal Statistics (4.8.6.1):")
    print(f"  Files: {total_files}")
    print(f"  Lines: {total_lines}")
    print(f"  Images: {line_types['Image lines']}")
    print(f"  Notes: {line_types['Plan notes'] + line_types['Post-set notes'] + line_types['Other notes']}")
    print(f"  Journal Rows: {line_types['Journal rows']}")
    print(f"  Tickers: {ticker_count}")
    print(f"  Tags: {total_tag_count}")
    print(f"  SNF Rows: {line_types['SNF lines']}")
    print()
    
    print("Notes Analysis (4.8.6.4):")
    print(f"  Plan notes: {line_types['Plan notes']}")
    print(f"  Post-set notes: {line_types['Post-set notes']}")
    print(f"  Other notes: {line_types['Other notes']}")
    print(f"  Total notes: {line_types['Plan notes'] + line_types['Post-set notes'] + line_types['Other notes']}")
    print()
    
    print("Complete Line Breakdown (4.8.6.5):")
    print(f"  Images: {line_types['Image lines']} ({line_types['Image lines']/total_lines*100:.1f}%)")
    print(f"  Plan notes: {line_types['Plan notes']} ({line_types['Plan notes']/total_lines*100:.1f}%)")
    print(f"  Post-set notes: {line_types['Post-set notes']} ({line_types['Post-set notes']/total_lines*100:.1f}%)")
    print(f"  Other notes: {line_types['Other notes']} ({line_types['Other notes']/total_lines*100:.1f}%)")
    print(f"  Journal rows: {line_types['Journal rows']} ({line_types['Journal rows']/total_lines*100:.1f}%)")
    print(f"  SNF lines: {line_types['SNF lines']} ({line_types['SNF lines']/total_lines*100:.1f}%)")
    print(f"  Collapsed lines: {line_types['Collapsed lines']} ({line_types['Collapsed lines']/total_lines*100:.1f}%)")
    print(f"  Code blocks: {line_types['Code blocks']} ({line_types['Code blocks']/total_lines*100:.1f}%)")
    print(f"  Empty lines: {line_types['Empty lines']} ({line_types['Empty lines']/total_lines*100:.1f}%)")
    print(f"  Other lines: {line_types['Other lines']} ({line_types['Other lines']/total_lines*100:.1f}%)")
    print()
    
    print("Exclusion Patterns (4.8.6.6):")
    print(f"  SNF headers: {line_types['SNF lines']}")
    print(f"  collapsed:: true: {line_types['Collapsed lines']}")
    print(f"  Empty lines: {line_types['Empty lines']}")
    print(f"  Total exclusions: {exclude_content} ({exclude_content/total_lines*100:.1f}%)")
    print()
    
    print("Ticker Analysis (4.8.6.7):")
    print(f"  Total ticker occurrences: {len(ticker_occurrences)}")
    print(f"  Unique tickers: {ticker_count}")
    print(f"  Average occurrences per ticker: {len(ticker_occurrences)/ticker_count:.1f}")
    print()
    
    print("=== Final Summary ===")
    print(f"Files: {total_files}")
    print(f"Lines: {total_lines}")
    print(f"Images: {line_types['Image lines']}")
    print(f"Notes: {line_types['Plan notes'] + line_types['Post-set notes'] + line_types['Other notes']}")
    print(f"Journal Rows: {line_types['Journal rows']}")
    print(f"Tickers: {ticker_count}")
    print(f"Tags: {total_tag_count}")
    print(f"SNF Exclusions: {line_types['SNF lines']}")
    print()
    print("Analysis complete.")

if __name__ == "__main__":
    analyze_journal_data()
