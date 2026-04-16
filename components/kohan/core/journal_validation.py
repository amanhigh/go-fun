#!/usr/bin/env python3
"""
Comprehensive Journal Validation and Analysis Script - STRICT LINE ACCOUNTING
Every line must be either:
1. PROCESSED (understood and will be migrated)
2. SKIPPED (understood but intentionally not migrated)
3. FLAGGED (unexplained - needs investigation)

Provides both validation (line accounting) and analysis (statistics, mapping verification)
Usage: python3 journal_validation.py
"""

import os
import re
import glob
from collections import Counter, defaultdict
from datetime import datetime

# Environment variable configuration
BARKAT_PATH = os.environ.get("BARKAT_PATH")

# Image filenames embed a date in multiple legacy formats.
IMAGE_DASH_DATE_PATTERN = re.compile(
    r"--(\d{4})-(\d{2})-(\d{2})-\d+\.(?:jpg|jpeg|png)$", re.IGNORECASE
)
IMAGE_DOUBLE_UNDERSCORE_DATE_PATTERN = re.compile(
    r"__(\d{4})(\d{2})(\d{2})__\d{6}\.(?:jpg|jpeg|png)$", re.IGNORECASE
)
IMAGE_SINGLE_UNDERSCORE_DATE_PATTERN = re.compile(
    r"_(\d{4})(\d{2})(\d{2})_\d{6}\.(?:jpg|jpeg|png)$", re.IGNORECASE
)

ALLOWED_SEQUENCES = {"MWD", "YR", "WDH"}
ALLOWED_TYPES = {"REJECTED", "TAKEN"}
ALLOWED_STATUSES = {
    "SET",
    "RUNNING",
    "SUCCESS",
    "FAIL",
    "MISSED",
    "JUST_LOSS",
    "BROKEN",
}
ALLOWED_IMAGE_TIMEFRAMES = {"DL", "WK", "MN", "TMN", "SMN", "YR"}
ALLOWED_TAG_TYPES = {"REASON", "MANAGEMENT", "DIRECTION"}
IMAGE_TIMEFRAME_ORDER = ["DL", "WK", "MN", "TMN"]
NOTE_CONTENT_MAX_LENGTH = 5000
FILENAME_DATE_FORMATS = ("%Y_%m_%d", "%Y-%m-%d")
STATUS_TAGS = {
    "fail": "FAIL",
    "success": "SUCCESS",
    "running": "RUNNING",
    "broken": "BROKEN",
    "missed": "MISSED",
    "miss": "MISSED",
    "justloss": "JUST_LOSS",
}
TRADE_TAG_MAPPINGS = {
    "t.trend": "tags -> tag: trend, type: DIRECTION",
    "t.ctrend": "tags -> tag: ctrend, type: DIRECTION",
    "t.mwd": "sequence -> MWD",
    "t.yr": "sequence -> YR",
    "t.wdh": "sequence -> WDH",
    "t.rejected": "type -> REJECTED",
    "t.set": "status -> SET",
    "t.taken": "type -> TAKEN",
    "t.fail": "status -> FAIL",
    "t.success": "status -> SUCCESS",
    "t.broken": "status -> BROKEN",
    "t.running": "status -> RUNNING",
    "t.miss": "status -> MISSED",
    "t.missed": "status -> MISSED",
    "t.justloss": "status -> JUST_LOSS",
    "t.full": "raw_markdown note content -> preserved source context",
}

# Line categories
PROCESSED = "PROCESSED"  # Data that will be migrated
SKIPPED = "SKIPPED"  # Understood but not migrated (e.g., SNF header, collapsed::)
FLAGGED = "FLAGGED"  # Unexplained - needs investigation


def extract_image_year_month(basename):
    for pattern in (
        IMAGE_DASH_DATE_PATTERN,
        IMAGE_DOUBLE_UNDERSCORE_DATE_PATTERN,
        IMAGE_SINGLE_UNDERSCORE_DATE_PATTERN,
    ):
        match = pattern.search(basename)
        if match:
            return match.group(1), match.group(2)
    return None


def resolve_image_candidates(image_path):
    normalized = image_path.replace("\\", "/")
    normalized = re.sub(r"/+/", "/", normalized)
    basename = os.path.basename(normalized)
    candidates = []

    trading_match = re.search(r"assets/trading/(.+)$", normalized)
    if trading_match:
        relative_path = trading_match.group(1).lstrip("/")
        candidates.append(("source_path", os.path.join(BARKAT_PATH, relative_path)))

    date_parts = extract_image_year_month(basename)
    if date_parts:
        year, month = date_parts
        candidates.append(
            (
                "filename_date",
                os.path.join(BARKAT_PATH, year, month, basename),
            )
        )

    candidates.append(("basename_fallback", os.path.join(BARKAT_PATH, basename)))

    unique_candidates = []
    seen = set()
    for source, candidate in candidates:
        normalized_candidate = os.path.normpath(candidate)
        if normalized_candidate in seen:
            continue
        seen.add(normalized_candidate)
        unique_candidates.append((source, normalized_candidate))

    return unique_candidates


class LineAccountant:
    """Tracks every line and ensures nothing is missed."""

    def __init__(self):
        self.total_lines = 0
        self.processed_lines = 0
        self.skipped_lines = 0
        self.flagged_lines = 0

        # Detailed tracking
        self.processed_details = defaultdict(
            list
        )  # category -> [(file, line_num, content)]
        self.skipped_details = defaultdict(list)
        self.flagged_details = defaultdict(list)

        # Data extraction
        self.tickers = []
        self.images = []
        self.tags = {"trade": [], "reason": [], "management": [], "other": []}
        self.notes = {"code_block": [], "simple": [], "plan": []}
        self.important_tags = []  # #important tags

        # Analysis data
        self.ticker_occurrences = []
        self.date_range = {"earliest": None, "latest": None}
        self.file_count = 0
        self.projection_issues = []
        self.projection_stats = defaultdict(int)
        self.unknown_trade_tags = Counter()
        self.source_only_trade_tags = Counter()
        self.entry_count = 0

        # Image existence tracking
        self.image_exists_results = {
            "found": 0,
            "missing": 0,
            "skipped_clipboard": 0,
            "found_by_source_path": 0,
            "found_by_filename_date": 0,
            "missing_paths": [],
        }
        self._seen_image_paths = set()  # Track unique full paths to avoid duplicates
        self.image_projection = defaultdict(int)

    def process_line(self, file_name, line_num, line, category, status):
        """Record a line's status."""
        self.total_lines += 1

        entry = {
            "file": file_name,
            "line": line_num,
            "content": line.rstrip()[:100],  # Truncate for display
            "full_content": line.rstrip(),
        }

        if status == PROCESSED:
            self.processed_lines += 1
            self.processed_details[category].append(entry)
        elif status == SKIPPED:
            self.skipped_lines += 1
            self.skipped_details[category].append(entry)
        else:  # FLAGGED
            self.flagged_lines += 1
            self.flagged_details[category].append(entry)

    def verify_totals(self):
        """Verify all lines are accounted for."""
        accounted = self.processed_lines + self.skipped_lines + self.flagged_lines
        return accounted == self.total_lines, accounted, self.total_lines

    def image_exists(self, file_name, line_num, image_path):
        """
        Check if an image exists on disk and track the result.
        Images are resolved first from the original source-relative path, then by
        filename-derived year/month, then by basename fallback.
        Only tracks unique image paths to avoid counting duplicates.
        Skips clipboard-paste images (image_* pattern) — these were pasted from
        clipboard into Logseq and the actual file was never on disk.
        """
        basename = os.path.basename(image_path)

        # Skip clipboard paste images — these were pasted into Logseq but the
        # actual file was never on disk (filename pattern: image_*)
        if basename.startswith("image_"):
            self.image_exists_results["skipped_clipboard"] = (
                self.image_exists_results.get("skipped_clipboard", 0) + 1
            )
            return

        candidates = resolve_image_candidates(image_path)
        dedupe_key = tuple(candidate for _, candidate in candidates)
        if dedupe_key in self._seen_image_paths:
            return
        self._seen_image_paths.add(dedupe_key)

        for source, candidate in candidates:
            if os.path.exists(candidate):
                self.image_exists_results["found"] += 1
                if source == "source_path":
                    self.image_exists_results["found_by_source_path"] += 1
                elif source == "filename_date":
                    self.image_exists_results["found_by_filename_date"] += 1
                return

        self.image_exists_results["missing"] += 1
        self.image_exists_results["missing_paths"].append(
            {
                "file": file_name,
                "line": line_num,
                "path": image_path,
                "candidates": [candidate for _, candidate in candidates],
                "full_path": candidates[0][1],
            }
        )


def sanitize_ticker(ticker):
    sanitized = ticker
    if sanitized.endswith("!"):
        sanitized = sanitized[:-1]
    return sanitized.replace("_", "")


def sanitize_filename(name):
    return re.sub(r'[!@#$%^&*()+=\[\]{}|;:\'",<>?/\\]', "_", name)


def choose_trade_type(current, candidate):
    priority = {
        "": 0,
        "REJECTED": 1,
        "TAKEN": 2,
    }
    if priority.get(candidate, 0) >= priority.get(current, 0):
        return candidate
    return current


def derive_status_from_type(journal_type):
    if journal_type == "TAKEN":
        return "SET"
    return "FAIL"


def extract_journal_date(file_name):
    base_name = file_name.replace(".md", "")
    for date_format in FILENAME_DATE_FORMATS:
        try:
            return datetime.strptime(base_name, date_format)
        except ValueError:
            continue
    return None


def parse_legacy_tags(tags_part, entry, accountant):
    tags = re.findall(r"#([trm])\.([a-z0-9-]+)", tags_part)
    for prefix, value in tags:
        if prefix == "t":
            if value in ("mwd", "yr"):
                entry["sequence"] = value.upper()
            elif value == "wdh":
                entry["sequence"] = "WDH"
            elif value == "rejected":
                entry["type"] = choose_trade_type(entry["type"], "REJECTED")
            elif value == "taken":
                entry["type"] = choose_trade_type(entry["type"], "TAKEN")
            elif value == "full":
                accountant.source_only_trade_tags["t.full"] += 1
            elif value == "set":
                entry["status"] = "SET"
            elif value in STATUS_TAGS:
                entry["status"] = STATUS_TAGS[value]
            elif value in ("trend", "ctrend"):
                entry["direction"] = value
            else:
                accountant.unknown_trade_tags[f"t.{value}"] += 1
        elif prefix == "r":
            entry["reason_tags"].append(value)
        elif prefix == "m":
            entry["management_tags"].append(value)

    if "#important" in tags_part:
        entry["is_important"] = True


def finalize_entry(entries, current_entry, raw_lines, note_lines):
    if current_entry is None:
        return

    if note_lines and not current_entry["note"]:
        current_entry["note"] = "\n".join(note_lines).strip()

    current_entry["raw_markdown"] = "\n".join(raw_lines).strip()
    entries.append(current_entry)


def parse_legacy_entries(file_path, accountant):
    entries = []
    current_entry = None
    in_code_block = False
    note_lines = []
    raw_lines = []
    entry_pattern = re.compile(r"\|\s*`([A-Z0-9_!]+)`\s*\|(.+)\|")
    image_pattern = re.compile(r"!\[.*?\]\(([^)]+)\)")

    with open(file_path, "r", encoding="utf-8") as handle:
        for line_num, raw_line in enumerate(handle, 1):
            line = raw_line.rstrip("\n")
            matches = entry_pattern.search(line)
            if matches:
                finalize_entry(entries, current_entry, raw_lines, note_lines)
                current_entry = {
                    "ticker": matches.group(1),
                    "sequence": "",
                    "type": "",
                    "status": "",
                    "direction": "",
                    "reason_tags": [],
                    "management_tags": [],
                    "is_important": False,
                    "images": [],
                    "note": "",
                    "simple_notes": [],
                    "raw_markdown": "",
                    "line_number": line_num,
                    "file_name": os.path.basename(file_path),
                }
                raw_lines = [line]
                note_lines = []
                in_code_block = False
                parse_legacy_tags(matches.group(2), current_entry, accountant)
                continue

            if current_entry is None:
                continue

            raw_lines.append(line)

            if "```" in line:
                if in_code_block:
                    current_entry["note"] = "\n".join(note_lines).strip()
                    note_lines = []
                in_code_block = not in_code_block
                continue

            if in_code_block:
                note_lines.append(line)
                continue

            image_match = image_pattern.search(line)
            if image_match:
                current_entry["images"].append(image_match.group(1))

            stripped = line.strip()
            if stripped.startswith("-"):
                content = stripped[1:].strip()
                if content and "::" not in content:
                    current_entry["simple_notes"].append(content)

    finalize_entry(entries, current_entry, raw_lines, note_lines)
    return entries


def extract_image_created_at(image_path, journal_date):
    basename = os.path.basename(image_path)
    for pattern in (
        IMAGE_DASH_DATE_PATTERN,
        IMAGE_DOUBLE_UNDERSCORE_DATE_PATTERN,
        IMAGE_SINGLE_UNDERSCORE_DATE_PATTERN,
    ):
        match = pattern.search(basename)
        if not match:
            continue
        return datetime.strptime("".join(match.groups()[:3]), "%Y%m%d")
    return journal_date


def build_projected_images(image_paths, journal_date, accountant):
    images = []
    for image_path in image_paths:
        basename = os.path.basename(image_path)
        if basename.startswith("image_"):
            accountant.image_projection["skipped_clipboard"] += 1
            continue

        if len(images) >= 16:
            accountant.image_projection["truncated_source_refs"] += 1
            continue

        images.append(
            {
                "timeframe": IMAGE_TIMEFRAME_ORDER[
                    len(images) % len(IMAGE_TIMEFRAME_ORDER)
                ],
                "file_name": sanitize_filename(basename),
                "created_at": extract_image_created_at(
                    image_path, journal_date
                ).strftime("%Y-%m-%d"),
            }
        )
        accountant.image_projection["migrated_source_refs"] += 1

    while len(images) < 4:
        images.append(
            {
                "timeframe": IMAGE_TIMEFRAME_ORDER[len(images)],
                "file_name": f"placeholder_{len(images)}.png",
                "created_at": journal_date.strftime("%Y-%m-%d"),
            }
        )
        accountant.image_projection["placeholder_images"] += 1

    accountant.image_projection["final_image_records"] += len(images)
    accountant.image_projection["source_refs"] += len(image_paths)
    return images


def build_journal_projection(entry, journal_date, accountant):
    ticker = sanitize_ticker(entry["ticker"])
    journal_type = entry["type"] or "REJECTED"
    status = entry["status"] or derive_status_from_type(journal_type)
    images = build_projected_images(entry["images"], journal_date, accountant)

    tags = []
    if entry["direction"]:
        tags.append({"tag": entry["direction"], "type": "DIRECTION", "override": None})

    for reason_tag in entry["reason_tags"]:
        parts = reason_tag.split("-", 1)
        tags.append(
            {
                "tag": parts[0],
                "type": "REASON",
                "override": parts[1] if len(parts) > 1 else None,
            }
        )

    for management_tag in entry["management_tags"]:
        tags.append({"tag": management_tag, "type": "MANAGEMENT", "override": None})

    if entry["is_important"]:
        tags.append({"tag": "important", "type": "MANAGEMENT", "override": None})

    raw_markdown = entry["raw_markdown"].strip()
    plan_note = entry["note"].strip()
    sections = []
    if raw_markdown:
        sections.append(f"=== ORIGINAL MARKDOWN ===\n{raw_markdown}")
    elif plan_note:
        sections.append(f"=== PLAN NOTES ===\n{plan_note}")

    if entry["simple_notes"]:
        review_lines = ["=== REVIEW NOTES ==="]
        review_lines.extend(f"- {note}" for note in entry["simple_notes"])
        sections.append("\n".join(review_lines))

    notes = []
    if sections:
        notes.append(
            {
                "status": status,
                "content": "\n\n".join(sections),
                "format": "MARKDOWN",
            }
        )

    return {
        "ticker": ticker,
        "sequence": entry["sequence"] or "MWD",
        "type": journal_type,
        "status": status,
        "created_at": journal_date.strftime("%Y-%m-%d"),
        "images": images,
        "tags": tags,
        "notes": notes,
    }


def validate_projection(journal, entry, accountant):
    location = f"{entry['file_name']}:{entry['line_number']}:{entry['ticker']}"

    if (
        not journal["ticker"]
        or len(journal["ticker"]) > 10
        or not re.fullmatch(r"[A-Z0-9]+", journal["ticker"])
    ):
        accountant.projection_issues.append(
            f"{location} invalid ticker after sanitization: {journal['ticker']}"
        )

    if journal["sequence"] not in ALLOWED_SEQUENCES:
        accountant.projection_issues.append(
            f"{location} invalid sequence: {journal['sequence']}"
        )

    if journal["type"] not in ALLOWED_TYPES:
        accountant.projection_issues.append(
            f"{location} invalid type: {journal['type']}"
        )

    if journal["status"] not in ALLOWED_STATUSES:
        accountant.projection_issues.append(
            f"{location} invalid status: {journal['status']}"
        )

    if not 4 <= len(journal["images"]) <= 16:
        accountant.projection_issues.append(
            f"{location} invalid image count: {len(journal['images'])}"
        )

    for image in journal["images"]:
        if image["timeframe"] not in ALLOWED_IMAGE_TIMEFRAMES:
            accountant.projection_issues.append(
                f"{location} invalid image timeframe: {image['timeframe']}"
            )
        if not image["file_name"] or len(image["file_name"]) > 255:
            accountant.projection_issues.append(
                f"{location} invalid image filename: {image['file_name']}"
            )

    if len(journal["tags"]) > 10:
        accountant.projection_issues.append(
            f"{location} tag count exceeds model limit: {len(journal['tags'])}"
        )

    seen_tags = set()
    for tag in journal["tags"]:
        key = (tag["tag"], tag["type"])
        if key in seen_tags:
            accountant.projection_issues.append(
                f"{location} duplicate tag/type generated: {tag['tag']}/{tag['type']}"
            )
        seen_tags.add(key)

        if not tag["tag"] or len(tag["tag"]) > 10:
            accountant.projection_issues.append(
                f"{location} invalid tag value: {tag['tag']}"
            )
        if tag["type"] not in ALLOWED_TAG_TYPES:
            accountant.projection_issues.append(
                f"{location} invalid tag type: {tag['type']}"
            )
        if tag["override"] is not None and len(tag["override"]) > 5:
            accountant.projection_issues.append(
                f"{location} invalid tag override length: {tag['override']}"
            )

    if len(journal["notes"]) > 1:
        accountant.projection_issues.append(
            f"{location} note count exceeds model limit: {len(journal['notes'])}"
        )

    if journal["notes"]:
        note = journal["notes"][0]
        content = note["content"].strip()
        if note["status"] not in ALLOWED_STATUSES:
            accountant.projection_issues.append(
                f"{location} invalid note status: {note['status']}"
            )
        if note["format"] not in ("MARKDOWN", "PLAINTEXT"):
            accountant.projection_issues.append(
                f"{location} invalid note format: {note['format']}"
            )
        if not content or len(content) > NOTE_CONTENT_MAX_LENGTH:
            accountant.projection_issues.append(
                f"{location} invalid note content length: {len(content)}"
            )
        if content in ("=== ORIGINAL MARKDOWN ===", "=== PLAN NOTES ==="):
            accountant.projection_issues.append(
                f"{location} placeholder-only note content detected"
            )

    accountant.projection_stats["journals"] += 1
    accountant.projection_stats["images"] += len(journal["images"])
    accountant.projection_stats["tags"] += len(journal["tags"])
    accountant.projection_stats["notes"] += len(journal["notes"])


def run_projection_validation(files, accountant):
    for file_path in files:
        file_name = os.path.basename(file_path)
        journal_date = extract_journal_date(file_name)
        if journal_date is None:
            accountant.projection_issues.append(
                f"{file_name} invalid filename date format"
            )
            continue

        entries = parse_legacy_entries(file_path, accountant)
        accountant.entry_count += len(entries)
        for entry in entries:
            validate_projection(
                build_journal_projection(entry, journal_date, accountant),
                entry,
                accountant,
            )


def analyze_journal_line(
    line, line_stripped, in_code_block, prev_line_type, file_name, line_num, accountant
):
    """
    Analyze a single line and categorize it.
    Returns: (status, category, new_in_code_block, line_type)
    """

    # Patterns
    image_pattern = r"!\[.*?\]\(([^)]+)\)"
    ticker_pattern = r"\|\s*`([A-Z0-9_!]+)`\s*\|"
    tag_pattern = r"#([tmr])\.([a-zA-Z0-9_-]+)"
    important_pattern = r"#important"

    # === CODE BLOCK HANDLING ===
    if "```" in line:
        return PROCESSED, "code_block_marker", not in_code_block, "code_block_marker"

    if in_code_block:
        # Check for Plan: inside code blocks
        if "Plan:" in line:
            accountant.notes["plan"].append(
                {
                    "file": file_name,
                    "line": line_num,
                    "content": line_stripped,
                    "type": "plan_in_code_block",
                }
            )
            return PROCESSED, "plan_note_in_code_block", True, "plan_note"
        return PROCESSED, "code_block_content", True, "code_block_content"

    # === EMPTY LINES ===
    if not line_stripped:
        return SKIPPED, "empty_line", False, "empty"

    # === SNF HEADER ===
    if "SNF" in line and "Journal" in line:
        return SKIPPED, "snf_header", False, "snf_header"

    # === COLLAPSED LINES ===
    if "collapsed::" in line:
        return SKIPPED, "collapsed_marker", False, "collapsed"

    # === JOURNAL ENTRY ROW (ticker + tags) ===
    ticker_match = re.search(ticker_pattern, line)
    if ticker_match and "|" in line:
        ticker = ticker_match.group(1)
        accountant.tickers.append(
            {"file": file_name, "line": line_num, "ticker": ticker}
        )
        accountant.ticker_occurrences.append(ticker)  # Add to analysis data

        # Extract all tags from this line
        tags = re.findall(tag_pattern, line)
        for tag_type, tag_value in tags:
            if tag_type == "t":
                accountant.tags["trade"].append(
                    {"file": file_name, "line": line_num, "tag": f"t.{tag_value}"}
                )
            elif tag_type == "r":
                accountant.tags["reason"].append(
                    {"file": file_name, "line": line_num, "tag": f"r.{tag_value}"}
                )
            elif tag_type == "m":
                accountant.tags["management"].append(
                    {"file": file_name, "line": line_num, "tag": f"m.{tag_value}"}
                )

        # Check for #important tag
        if re.search(important_pattern, line):
            accountant.important_tags.append(
                {"file": file_name, "line": line_num, "content": line_stripped}
            )

        # Check for any OTHER tags not matching #t./#r./#m. pattern
        # Use negative lookahead to exclude #t.xxx, #r.xxx, #m.xxx patterns
        other_tag_pattern = r"#([a-zA-Z][a-zA-Z0-9_-]*)(?![.\w])"
        other_tags = re.findall(other_tag_pattern, line)
        for tag in other_tags:
            # Skip known non-data tags and the base t/r/m (which are part of #t.xxx patterns)
            if (
                tag not in ["trading-tome", "t", "r", "m", "important"]
                and not tag.startswith("t.")
                and not tag.startswith("r.")
                and not tag.startswith("m.")
            ):
                accountant.tags["other"].append(
                    {"file": file_name, "line": line_num, "tag": f"#{tag}"}
                )

        return PROCESSED, "journal_entry_row", False, "journal_row"

    # === IMAGE LINES ===
    image_match = re.search(image_pattern, line)
    if image_match:
        image_path = image_match.group(1)
        accountant.images.append(
            {"file": file_name, "line": line_num, "path": image_path}
        )

        # Track image existence on disk (skips clipboard paste images)
        accountant.image_exists(file_name, line_num, image_path)

        # Check if there's text BEFORE or AFTER the image on the same line
        line_without_image = re.sub(image_pattern, "", line).strip()
        line_without_image = re.sub(
            r"^-\s*", "", line_without_image
        ).strip()  # Remove leading -
        # Remove Logseq image dimension metadata {:height X, :width Y}
        line_without_image = re.sub(
            r"\{:height\s+\d+,?\s*:width\s+\d+\}", "", line_without_image
        ).strip()

        if line_without_image:
            # There's additional content on the image line that's NOT just metadata!
            return FLAGGED, "image_line_with_extra_content", False, "image_with_extra"

        return PROCESSED, "image_line", False, "image"

    # === SIMPLE NOTES (lines starting with - but not images) ===
    if line_stripped.startswith("-"):
        content = line_stripped[1:].strip()

        # Check if it's just a dash with no content
        if not content:
            return SKIPPED, "empty_dash_line", False, "empty_dash"

        # Check for background-color:: (Logseq property)
        if "background-color::" in content:
            return SKIPPED, "logseq_property", False, "logseq_property"

        # This is a SIMPLE NOTE - needs to be migrated!
        accountant.notes["simple"].append(
            {
                "file": file_name,
                "line": line_num,
                "content": content,
                "prev_line_type": prev_line_type,
            }
        )
        return PROCESSED, "simple_note", False, "simple_note"

    # === LOGSEQ PROPERTIES ===
    if "::" in line and not line_stripped.startswith("|"):
        # Logseq property like background-color:: yellow
        return SKIPPED, "logseq_property", False, "logseq_property"

    # === ANYTHING ELSE IS FLAGGED ===
    return FLAGGED, "unexplained_line", False, "unexplained"


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
        accountant.file_count += 1

        # Extract date from filename for analysis
        date_match = re.search(r"(\d{4}_\d{2}_\d{2})", file_name)
        if date_match:
            date_str = date_match.group(1)
            formatted_date = f"{date_str[:4]}-{date_str[5:7]}-{date_str[8:10]}"

            if (
                not accountant.date_range["earliest"]
                or formatted_date < accountant.date_range["earliest"]
            ):
                accountant.date_range["earliest"] = formatted_date
            if (
                not accountant.date_range["latest"]
                or formatted_date > accountant.date_range["latest"]
            ):
                accountant.date_range["latest"] = formatted_date

        with open(file_path, "r", encoding="utf-8") as f:
            lines = f.readlines()

        in_code_block = False
        prev_line_type = None

        for line_num, line in enumerate(lines, 1):
            line_stripped = line.strip()

            status, category, in_code_block, line_type = analyze_journal_line(
                line,
                line_stripped,
                in_code_block,
                prev_line_type,
                file_name,
                line_num,
                accountant,
            )

            accountant.process_line(file_name, line_num, line, category, status)
            prev_line_type = line_type

    run_projection_validation(files, accountant)

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
    print(
        f"\nVerification: {'✓ ALL LINES ACCOUNTED' if is_valid else '✗ MISMATCH - LINES MISSING'}"
    )

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
    print(f"  Simple Notes (review notes): {len(accountant.notes['simple'])}")

    print("\n--- Projected Migration Shape ---")
    print(f"  Projected Journals: {accountant.projection_stats['journals']}")
    print(f"  Projected Images: {accountant.projection_stats['images']}")
    print(f"  Projected Tags: {accountant.projection_stats['tags']}")
    print(f"  Projected Notes: {accountant.projection_stats['notes']}")

    print("\n--- Projected Image Accounting ---")
    print(f"  Source image refs: {accountant.image_projection['source_refs']}")
    print(
        f"  Migrated source refs: {accountant.image_projection['migrated_source_refs']}"
    )
    print(
        f"  Skipped clipboard refs: {accountant.image_projection['skipped_clipboard']}"
    )
    print(
        f"  Truncated source refs: {accountant.image_projection['truncated_source_refs']}"
    )
    print(
        f"  Placeholder images added: {accountant.image_projection['placeholder_images']}"
    )
    print(
        f"  Final DB image records: {accountant.image_projection['final_image_records']}"
    )

    # === SIMPLE NOTES DETAIL (CRITICAL) ===
    if accountant.notes["simple"]:
        print("\n" + "=" * 80)
        print("SIMPLE NOTES CAPTURED FROM SOURCE (Outside Code Blocks)")
        print("=" * 80)
        print(
            f"\nFound {len(accountant.notes['simple'])} simple notes that need migration:"
        )
        for i, note in enumerate(accountant.notes["simple"], 1):
            print(f"\n  {i}. {note['file']}:{note['line']}")
            print(
                f"     Content: {note['content'][:80]}{'...' if len(note['content']) > 80 else ''}"
            )
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
    if accountant.tags["other"]:
        print("\n" + "=" * 80)
        print("OTHER TAGS (Non-standard #t./#r./#m. tags)")
        print("=" * 80)
        other_tag_counter = Counter([t["tag"] for t in accountant.tags["other"]])
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
    trade_counter = Counter([t["tag"] for t in accountant.tags["trade"]])
    for tag, count in trade_counter.most_common():
        print(f"  {tag}: {count}")

    print("\n--- Reason Tags (#r.*) ---")
    reason_counter = Counter([t["tag"] for t in accountant.tags["reason"]])
    for tag, count in reason_counter.most_common(20):
        print(f"  {tag}: {count}")
    if len(reason_counter) > 20:
        print(f"  ... and {len(reason_counter) - 20} more unique tags")

    print("\n--- Management Tags (#m.*) ---")
    mgmt_counter = Counter([t["tag"] for t in accountant.tags["management"]])
    for tag, count in mgmt_counter.most_common():
        print(f"  {tag}: {count}")

    # === COMPREHENSIVE ANALYSIS SECTION ===
    print("\n" + "=" * 80)
    print("COMPREHENSIVE JOURNAL ANALYSIS")
    print("=" * 80)

    # Date range analysis
    print("\n=== Date Range Analysis ===")
    if accountant.date_range["earliest"] and accountant.date_range["latest"]:
        print(
            f"Date range: {accountant.date_range['earliest']} to {accountant.date_range['latest']}"
        )
        print(f"Total dates with entries: {accountant.file_count}")

    # Detailed ticker analysis
    print("\n=== Ticker Analysis ===")
    ticker_counter = Counter(accountant.ticker_occurrences)
    print(f"Total ticker occurrences: {len(accountant.ticker_occurrences)}")
    print(f"Unique tickers: {len(ticker_counter)}")
    print(
        f"Average occurrences per ticker: {len(accountant.ticker_occurrences) / len(ticker_counter):.1f}"
    )

    print("\nTop 20 Most Active Tickers:")
    for ticker, count in ticker_counter.most_common(20):
        print(f"    {count:3d} {ticker}")

    # PRD validation counts
    print("\n=== PRD Validation Counts ===")
    print("These counts should match section 4.8.6 in the PRD:")
    print()

    print("Journal Statistics (4.8.6.1):")
    print(f"  Files: {accountant.file_count}")
    print(f"  Lines: {accountant.total_lines}")
    print(f"  Images: {len(accountant.images)}")
    print(
        f"  Notes: {len(accountant.notes['code_block']) + len(accountant.notes['simple']) + len(accountant.notes['plan'])}"
    )
    print(f"  Journal Rows: {len(accountant.tickers)}")
    print(f"  Tickers: {len(ticker_counter)}")
    print(
        f"  Tags: {len(accountant.tags['trade']) + len(accountant.tags['reason']) + len(accountant.tags['management'])}"
    )
    print(f"  SNF Rows: {len(accountant.skipped_details.get('snf_header', []))}")

    # Tag mapping verification for migration
    print("\n=== Tag Mapping Verification (PRD 4.8.6.3) ===")
    print("\nTrade Tag Mappings:")
    for tag, count in trade_counter.most_common():
        mapping = TRADE_TAG_MAPPINGS.get(tag, "UNKNOWN - needs mapping")
        print(f"   {count:4d} {tag} -> {mapping}")

    if accountant.source_only_trade_tags:
        print("\n--- Source-Only Trade Tags (preserved in raw markdown notes) ---")
        for tag, count in accountant.source_only_trade_tags.most_common():
            print(f"  {tag}: {count}")

    if accountant.unknown_trade_tags:
        print("\n--- Unknown Trade Tags In Source ---")
        for tag, count in accountant.unknown_trade_tags.most_common():
            print(f"  {tag}: {count}")

    print("\n=== Projected Migration Validation ===")
    if accountant.projection_issues:
        print(f"Found {len(accountant.projection_issues)} projected migration issues:")
        for issue in accountant.projection_issues[:50]:
            print(f"  {issue}")
        if len(accountant.projection_issues) > 50:
            print(f"  ... and {len(accountant.projection_issues) - 50} more")
    else:
        print("  ✓ No projection issues detected against current migration rules")

    print("\nReason Tag Mappings (all map to tags with type: REASON):")
    for tag, count in reason_counter.most_common():
        # Check for override pattern (e.g., r.dep-loc)
        tag_value = tag[2:]  # Remove 'r.' prefix
        if "-" in tag_value:
            parts = tag_value.split("-", 1)
            print(f"    {count:3d} {tag} -> tag: {parts[0]}, override: {parts[1]}")
        else:
            print(f"    {count:3d} {tag} -> tag: {tag_value}")

    print("\nManagement Tag Mappings (all map to tags with type: MANAGEMENT):")
    for tag, count in mgmt_counter.most_common():
        tag_value = tag[2:]  # Remove 'm.' prefix
        print(f"    {count:3d} {tag} -> tag: {tag_value}")

    # === IMAGE FILE EXISTENCE VALIDATION ===
    print("\n" + "=" * 80)
    print("IMAGE FILE EXISTENCE VALIDATION")
    print("=" * 80)

    print(f"\nBARKAT_PATH: {BARKAT_PATH}")
    total_checked = (
        accountant.image_exists_results["found"]
        + accountant.image_exists_results["missing"]
    )
    skipped_clip = accountant.image_exists_results.get("skipped_clipboard", 0)
    print(f"Total unique images checked (excl. clipboard pastes): {total_checked}")
    print(f"  Found on disk:     {accountant.image_exists_results['found']}")
    print(
        f"    - resolved from source path: {accountant.image_exists_results['found_by_source_path']}"
    )
    print(
        f"    - resolved from filename date: {accountant.image_exists_results['found_by_filename_date']}"
    )
    print(f"  Missing from disk: {accountant.image_exists_results['missing']}")
    print(f"  Skipped (clipboard paste — never existed on disk): {skipped_clip}")

    if accountant.image_exists_results["missing"] > 0:
        print("\n--- Missing Image Files ---")
        print(f"{'File':<30} {'Line':<6} {'Referenced Path':<50} {'Full Path Tried'}")
        print("-" * 130)
        for entry in accountant.image_exists_results["missing_paths"][:100]:
            print(
                f"{entry['file']:<30} {entry['line']:<6} "
                f"{entry['path']:<50} {entry['full_path']}"
            )
            if len(entry["candidates"]) > 1:
                for candidate in entry["candidates"][1:]:
                    print(f"{'':<30} {'':<6} {'':<50} alt -> {candidate}")
        if len(accountant.image_exists_results["missing_paths"]) > 100:
            print(
                f"... and "
                f"{len(accountant.image_exists_results['missing_paths']) - 100} more"
            )

    # === FINAL SUMMARY ===
    print("\n" + "=" * 80)
    print("FINAL VALIDATION SUMMARY")
    print("=" * 80)

    issues = []

    if accountant.flagged_lines > 0:
        issues.append(
            f"⚠️  {accountant.flagged_lines} unexplained lines need investigation"
        )

    if accountant.tags["other"]:
        issues.append(
            f"ℹ️  {len(accountant.tags['other'])} non-standard tags found - verify handling"
        )

    if accountant.unknown_trade_tags:
        issues.append(
            f"⚠️  {sum(accountant.unknown_trade_tags.values())} trade tags are not mapped by current migration"
        )

    if accountant.projection_issues:
        issues.append(
            f"⚠️  {len(accountant.projection_issues)} projected migration issues violate current Barkat constraints"
        )

    projected_source_accounted = (
        accountant.image_projection["migrated_source_refs"]
        + accountant.image_projection["skipped_clipboard"]
        + accountant.image_projection["truncated_source_refs"]
    )
    if projected_source_accounted != accountant.image_projection["source_refs"]:
        issues.append(
            "⚠️  projected image source accounting does not reconcile to migrated + skipped + truncated"
        )

    if issues:
        print("\nISSUES FOUND:")
        for issue in issues:
            print(f"  {issue}")
    else:
        print("\n✓ All validation checks passed!")

    print("\n--- Migration Counts (for comparison with migration script) ---")
    print(f"  Files: {len(files)}")
    print(f"  Journal Entries (tickers): {len(accountant.tickers)}")
    print(f"  Parsed Entries: {accountant.entry_count}")
    print(f"  Images: {len(accountant.images)}")
    print(
        f"  Total Tags: {len(accountant.tags['trade']) + len(accountant.tags['reason']) + len(accountant.tags['management'])}"
    )
    print(
        f"  Notes (code block + simple + plan): {len(accountant.notes['code_block']) + len(accountant.notes['simple']) + len(accountant.notes['plan'])}"
    )

    return accountant


if __name__ == "__main__":
    validate_journals()
