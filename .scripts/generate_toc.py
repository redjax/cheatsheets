"""Generate or update Table of Contents in Markdown files using GitHub-style anchors.

This script scans specified Markdown files for headings, generates a Table of Contents
(ToC) with links to those headings, and inserts or updates the ToC in the files.

Usage:
    python generate_toc.py [options] <file_or_directory>

Options:
    -h, --help           Show this help message and exit.
    --min-level N        Minimum heading level to include (default: 2).
    --max-level N        Maximum heading level to include (default: 6).
    --heading-text TEXT  Heading line that marks the ToC section.
    --write              Write changes in-place. Without this, runs in dry mode.
    --no-diff            Suppress diff output on dry runs.
"""

import argparse
import difflib
import re
import sys
from pathlib import Path

## Regex to find <!-- omit in toc --> lines to skip generating ToC entry
HEADING_RE = re.compile(r"^(#{1,6})\s+(.*?)(\s*<!--\s*omit in toc\s*-->\s*)?$")
## Regex to find start of frontmatter (--- at the top)
FRONT_MATTER_START = re.compile(r"^\s*---\s*$")
## Regex to find end of frontmatter (--- at the end of frontmatter)
FRONT_MATTER_END = re.compile(r"^\s*---\s*$")


def parse_args() -> argparse.Namespace:
    """Parse command-line arguments."""
    parser = argparse.ArgumentParser(
        description="Update Markdown Table of Contents using GitHub-style anchors."
    )

    parser.add_argument(
        "paths", nargs="*", help="Markdown files or directories to update."
    )
    parser.add_argument(
        "--min-level",
        type=int,
        default=2,
        help="Minimum heading level to include (default: 2).",
    )
    parser.add_argument(
        "--max-level",
        type=int,
        default=6,
        help="Maximum heading level to include (default: 6).",
    )
    parser.add_argument(
        "--heading-text",
        default="## Table of Contents <!-- omit in toc -->",
        help="Heading line that marks the ToC section.",
    )
    parser.add_argument("--write", action="store_true", help="Write changes in-place.")
    parser.add_argument(
        "--no-diff", action="store_true", help="Suppress diff output on dry runs."
    )

    args = parser.parse_args()

    return args


def github_slugify(text: str) -> str:
    """Generate GitHub-style slug from heading text.

    Params:
        text (str): Heading text to slugify.

    Returns:
        str: GitHub-style slug.
    """
    text = re.sub(r"`+([^`]*)`+", r"\1", text)
    text = re.sub(r"<[^>]+>", "", text)

    text = text.lower()

    text = re.sub(r"\s+", "-", text)
    text = re.sub(r"[^\w\-]", "", text)
    text = re.sub(r"-{2,}", "-", text)

    text = text.strip("-_")

    return text


def parse_headings(lines: list[str]) -> list[tuple[int, str, str]]:
    """Parse headings from markdown lines, returning list of (level, text, unique_slug).

    Params:
        lines (list[str]): Lines of the markdown file.

    Returns:
        list[tuple[int, str, str]]: List of (level, text, unique_slug).
    """
    headings: list[tuple[int, str, str]] = []
    in_front_matter = False

    for i, line in enumerate(lines):
        ## Skip front matter
        if i == 0 and FRONT_MATTER_START.match(line):
            in_front_matter = True
            continue

        if in_front_matter:
            ## Find end of front matter
            if FRONT_MATTER_END.match(line):
                in_front_matter = False
            continue

        ## Match headings with regex
        m: re.Match | None = HEADING_RE.match(line.rstrip())

        if not m:
            continue

        ## Extract data from match
        hashes, text, omit = m.groups()

        if omit is not None:
            continue

        ## Determine heading level and slug
        level: int = len(hashes)
        slug: str = github_slugify(text)

        headings.append((level, text, slug))

    counts: dict[str, int] = {}
    uniq: list[tuple[int, str, str]] = []

    ## Iterate over headings to ensure unique slugs
    for level, text, slug in headings:
        ## Count occurrences of each slug
        n = counts.get(slug, 0)
        ## Update count
        counts[slug] = n + 1
        ## Create unique slug
        uniq_slug = slug if n == 0 else f"{slug}-{n}"

        uniq.append((level, text, uniq_slug))

    return uniq


def build_toc(
    headings: list[tuple[int, str, str]], min_level: int = 2, max_level: int = 6
) -> str:
    """Build ToC markdown from headings within specified levels.

    Params:
        headings (list[tuple[int, str, str]]): List of (level, text, slug).
        min_level (int): Minimum heading level to include.
        max_level (int): Maximum heading level to include.

    Returns:
        str: Generated ToC markdown.
    """
    ## Filter headings by level
    filtered = [
        (lvl, txt, slug)
        for (lvl, txt, slug) in headings
        if min_level <= lvl <= max_level
    ]

    if not filtered:
        return ""

    ## Determine base level for indentation
    base: int = min(lvl for (lvl, _, _) in filtered)
    ## Build ToC lines
    out: list[str] = []

    ## Create markdown list entries
    for lvl, txt, slug in filtered:
        indent = "  " * (lvl - base)
        out.append(f"{indent}- [{txt}](#{slug})")

    return "\n".join(out) + "\n"


def find_toc_section(
    lines, heading_text="## Table of Contents <!-- omit in toc -->"
) -> tuple[int, int] | None:
    """Find the range of lines for the ToC section.

    Params:
        lines (list[str]): Lines of the markdown file.
        heading_text (str): Heading line that marks the ToC section.

    Returns:
        tuple[int, int] | None: (start_idx, end_idx) of ToC section, or None if not found.
    """
    idx = None

    ## Find the heading line
    for i, line in enumerate(lines):
        ## Exact match for heading text
        if line.strip() == heading_text.strip():
            idx = i
            break

    if idx is None:
        return None

    ## Determine heading level
    m: re.Match | None = re.match(r"^(#{1,6})\s+", lines[idx])
    ## Default to level 2 if not found
    level = len(m.group(1)) if m else 2

    end: int = len(lines)

    ## Find the end of the ToC section
    for j in range(idx + 1, len(lines)):
        m2 = HEADING_RE.match(lines[j].rstrip())

        if m2:
            lvl2 = len(m2.group(1))

            ## Stop if a heading of same or higher level is found
            if lvl2 <= level:
                end = j
                break

    return (idx, end)


def insert_toc(
    lines: list[str],
    toc_md: str,
    heading_text: str = "## Table of Contents <!-- omit in toc -->",
) -> list[str]:
    """Insert or replace ToC under the specified heading.

    Params:
        lines (list[str]): Lines of the markdown file.
        toc_md (str): Generated ToC markdown.
        heading_text (str): Heading line that marks the ToC section.

    Returns:
        list[str]: Updated lines with ToC inserted or replaced.
    """
    ## Find existing ToC section
    rng: tuple[int, int] | None = find_toc_section(lines, heading_text)

    ## Ensure one blank line before and after the ToC list
    toc_block = ["\n", toc_md, "\n"] if toc_md else ["\n"]

    if rng is None:
        ## Insert after front matter if present, else after first H1, else at top
        insert_at = 0

        ## Check for front matter
        if len(lines) >= 2 and FRONT_MATTER_START.match(lines[0]):
            for i in range(1, len(lines)):
                ## Find end of front matter
                if FRONT_MATTER_END.match(lines[i]):
                    ## Insert after front matter
                    insert_at = i + 1
                    break

        ## Check for first H1 heading
        else:
            ## Look for first H1 heading
            for i, line in enumerate(lines):
                ## Insert after first H1
                if re.match(r"^#\s+", line):
                    insert_at = i + 1
                    break

        return lines[:insert_at] + [heading_text + "\n"] + toc_block + lines[insert_at:]

    ## Replace existing ToC section
    else:
        start, end = rng

        return lines[: start + 1] + toc_block + lines[end:]


def normalize_newlines(s: str) -> str:
    """Normalize newlines to Unix style and ensure single trailing newline.

    Params:
        s (str): Input string.

    Returns:
        str: Normalized string.
    """
    s: str = s.replace("\r\n", "\n").replace("\r", "\n")
    s: str = re.sub(r"\n{3,}", "\n\n", s)

    if not s.endswith("\n"):
        s += "\n"

    return s


def process_file(path: Path, min_level: int, max_level: int, heading_text: str):
    """Process a single markdown file to update its ToC.

    Params:
        path (Path): Path to the markdown file.
        min_level (int): Minimum heading level to include in ToC.
        max_level (int): Maximum heading level to include in ToC.
        heading_text (str): Heading line that marks the ToC section.

    Returns:
        tuple[bool, str, str]: (changed, old_text, new_text)
    """
    text: str = path.read_text(encoding="utf-8")
    ## Split text into lines
    lines: list[str] = text.splitlines(keepends=True)
    ## Parse headings from lines
    headings: list[tuple[int, str]] = parse_headings(lines)
    ## Build ToC markdown
    toc_md: str = build_toc(headings, min_level=min_level, max_level=max_level)
    ## Insert or replace ToC in lines
    new_lines: list[str] = insert_toc(lines, toc_md, heading_text=heading_text)
    ## Join lines and normalize newlines
    new_text: str = normalize_newlines("".join(new_lines))

    ## Determine if changes were made
    changed = new_text != text

    return changed, text, new_text


def main():
    args = parse_args()

    if not args.paths:
        print("No paths provided. Pass files or directories.", file=sys.stderr)
        sys.exit(2)

    targets: list[Path] = []

    ## Collect markdown files from provided paths
    for p in args.paths:
        path = Path(p)

        if path.is_dir():
            targets.extend(sorted(path.rglob("*.md")))

        elif path.suffix.lower() == ".md":
            targets.append(path)

    if not targets:
        print("No Markdown files found.", file=sys.stderr)
        sys.exit(1)

    ## Process each markdown file
    for md in targets:
        ## Process the file to update ToC
        changed, old_text, new_text = process_file(
            md, args.min_level, args.max_level, args.heading_text
        )

        if args.write:
            ## Write changes in-place if --write is specified
            if changed:
                Path(md).write_text(new_text, encoding="utf-8", newline="\n")
                print(f"updated: {md}")

            ## Report unchanged files
            else:
                print(f"ok: {md}")

        ## Dry run mode
        else:
            ## Show what would be changed
            if changed:
                print(f"would update: {md}")

                if not args.no_diff:
                    diff = difflib.unified_diff(
                        old_text.splitlines(),
                        new_text.splitlines(),
                        fromfile=str(md),
                        tofile=str(md),
                        lineterm="",
                    )

                    for line in diff:
                        print(line)

            ## Report unchanged files
            else:
                print(f"ok: {md}")

    sys.exit(0)


if __name__ == "__main__":
    try:
        main()
    except Exception as exc:
        print(f"Error: {exc}", file=sys.stderr)
        sys.exit(1)
