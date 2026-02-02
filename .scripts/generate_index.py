import argparse
import datetime
import re
from pathlib import Path
from typing import Dict, List, Tuple, Optional

ROOT_DEFAULT = Path(__file__).parent.parent

EXCLUDE_FILE_NAMES = {"README.md", "INDEX.md"}
EXCLUDE_DIR_NAMES = {".git", ".github", ".vscode", ".templates", "__pycache__"}


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Generate INDEX.md for cheatsheets with nested folders."
    )

    parser.add_argument(
        "--root",
        type=Path,
        default=ROOT_DEFAULT,
        help="Root directory of cheatsheets (default: script directory)",
    )
    parser.add_argument(
        "--output",
        type=Path,
        default=None,
        help="Output path for index (default: <root>/INDEX.md)",
    )
    parser.add_argument(
        "--include-dotdirs",
        action="store_true",
        help="Include dot-prefixed directories (e.g., .templates)",
    )
    parser.add_argument(
        "--include-templates",
        action="store_true",
        help="Include .templates directory contents",
    )
    parser.add_argument(
        "--include-readme",
        action="store_true",
        help="Include README.md and existing INDEX.md in listing",
    )

    args = parser.parse_args()

    return args


def read_front_matter(path: Path) -> Dict[str, object]:
    """Parse minimal YAML front matter (title, description, tags) without external libs.

    Returns:
        (dict): A dict with keys: title (str|None), description (str|None), tags (List[str]).
    """
    res: Dict[str, object] = {"title": None, "description": None, "tags": []}

    try:
        text = path.read_text(encoding="utf-8")
    except Exception:
        return res

    if not text.startswith("---"):
        return res

    lines = text.splitlines()

    if not lines:
        return res

    in_front = False
    i = 0

    while i < len(lines):
        line = lines[i]

        if not in_front:
            if line.strip() == "---":
                in_front = True
            i += 1
            continue

        if line.strip() == "---":
            break

        low = line.lower()

        if low.startswith("title:"):
            value = line.split(":", 1)[1].strip()

            if (value.startswith('"') and value.endswith('"')) or (
                value.startswith("'") and value.endswith("'")
            ):
                value = value[1:-1]

            res["title"] = value

        elif low.startswith("description:"):
            value = line.split(":", 1)[1].strip()

            if (value.startswith('"') and value.endswith('"')) or (
                value.startswith("'") and value.endswith("'")
            ):
                value = value[1:-1]

            res["description"] = value

        elif low.startswith("tags:"):
            value = line.split(":", 1)[1].strip()
            tags: List[str] = []

            if value.startswith("["):
                ## Inline list; may span lines
                inner = value

                while not inner.endswith("]") and i + 1 < len(lines):
                    i += 1
                    inner += lines[i].strip()

                inner = inner.lstrip("[").rstrip("]")

                for part in inner.split(","):
                    t = part.strip()

                    if (t.startswith('"') and t.endswith('"')) or (
                        t.startswith("'") and t.endswith("'")
                    ):
                        t = t[1:-1]

                    if t:
                        tags.append(t)
            else:
                ## Attempt to read block list with '- '
                j = i + 1

                while j < len(lines):
                    nxt = lines[j].strip()

                    if nxt.startswith("- "):
                        t = nxt[2:].strip()

                        if (t.startswith('"') and t.endswith('"')) or (
                            t.startswith("'") and t.endswith("'")
                        ):
                            t = t[1:-1]

                        if t:
                            tags.append(t)

                        j += 1
                        continue

                    if nxt == "" or ":" in nxt:
                        break

                    j += 1

            res["tags"] = tags

        i += 1

    return res


def display_name_for_file(path: Path) -> str:
    # Prefer the first H1 heading in the document for previews
    def extract_h1_title(p: Path) -> Optional[str]:
        try:
            text = p.read_text(encoding="utf-8")
        except Exception:
            return None
        lines = text.splitlines()
        i = 0
        # Skip front matter if present
        if i < len(lines) and lines[i].strip() == "---":
            i += 1
            while i < len(lines) and lines[i].strip() != "---":
                i += 1
            if i < len(lines) and lines[i].strip() == "---":
                i += 1
        # Find first top-level heading
        while i < len(lines):
            line = lines[i].strip()
            if line.startswith("# "):
                return line[2:].strip()
            i += 1
        return None

    h1 = extract_h1_title(path)
    if h1:
        return h1

    fm = read_front_matter(path)
    title = fm.get("title")
    if isinstance(title, str) and title.strip():
        return title

    stem = path.stem.replace("_", " ").replace("-", " ")
    return stem.title()


def get_description_for_file(path: Path) -> str:
    fm = read_front_matter(path)
    desc = fm.get("description")

    return desc if isinstance(desc, str) else ""


def get_tags_for_file(path: Path) -> List[str]:
    fm = read_front_matter(path)
    tags = fm.get("tags")

    if isinstance(tags, list):
        return [str(t) for t in tags if str(t).strip()]

    return []


class Tree:
    def __init__(self):
        self.files: List[Path] = []
        self.dirs: Dict[str, "Tree"] = {}

    def add(self, rel_parts: Tuple[str, ...], full_path: Path):
        if not rel_parts:
            self.files.append(full_path)
            return

        head, *tail = rel_parts

        if head not in self.dirs:
            self.dirs[head] = Tree()

        self.dirs[head].add(tuple(tail), full_path)


def build_tree(
    root: Path, include_dotdirs: bool, include_templates: bool, include_readme: bool
) -> Tree:
    tree = Tree()

    for p in root.rglob("*.md"):
        rel = p.relative_to(root)

        ## Exclude files
        name = rel.name

        if not include_readme and name in EXCLUDE_FILE_NAMES:
            continue

        ## Exclude dirs
        parts = rel.parts
        excluded_dirs = EXCLUDE_DIR_NAMES - (
            {".templates"} if include_templates else set()
        )

        if any(part in excluded_dirs for part in parts):
            continue

        if not include_dotdirs:
            if any(
                part.startswith(".") and (not include_templates or part != ".templates")
                for part in parts
            ):
                continue

        tree.add(parts[:-1], p)

    return tree


def render_files_table(files: List[Path], root: Path) -> List[str]:
    lines: List[str] = []

    if not files:
        return lines

    lines.append("| Name | Description | Tags |")
    lines.append("| --- | --- | --- |")

    for f in sorted(files, key=lambda x: display_name_for_file(x).lower()):
        rel_link = f.relative_to(root).as_posix()
        name = display_name_for_file(f)
        desc = get_description_for_file(f)
        tags = ", ".join(get_tags_for_file(f))
        lines.append(f"| [{name}]({rel_link}) | {desc} | {tags} |")

    lines.append("")

    return lines


def render_tree_tables(
    tree: Tree, root: Path, depth: int = 2, section_name: Optional[str] = None
) -> List[str]:
    lines: List[str] = []

    if section_name is not None:
        heading_marks = "#" * max(2, min(6, depth))
        lines.append(f"{heading_marks} {section_name}")
        lines.append("")

    lines.extend(render_files_table(tree.files, root))

    for dname in sorted(tree.dirs.keys(), key=str.lower):
        lines.extend(render_tree_tables(tree.dirs[dname], root, depth + 1, dname))

    return lines


def generate_index_parts(
    root: Path, include_dotdirs: bool, include_templates: bool, include_readme: bool
) -> tuple[str, str]:
    tree = build_tree(root, include_dotdirs, include_templates, include_readme)

    ## Build ToC of sections (no trailing slashes)
    def slugify_heading(text: str) -> str:
        slug = re.sub(r"[^\w\s-]", "", text).strip().lower()
        slug = re.sub(r"[\s]+", "-", slug)

        return slug

    def render_sections_toc(node: Tree) -> List[str]:
        toc_lines: List[str] = []

        for dname in sorted(node.dirs.keys(), key=str.lower):
            anchor = slugify_heading(dname)
            toc_lines.append(f"- [{dname}](#{anchor})")

        return toc_lines

    toc_lines = render_sections_toc(tree)
    toc = "\n".join(toc_lines)

    ## Build body with top-level files and per-folder sections
    body_lines: List[str] = []
    body_lines.extend(render_files_table(tree.files, root))

    for dname in sorted(tree.dirs.keys(), key=str.lower):
        body_lines.extend(
            render_tree_tables(tree.dirs[dname], root, depth=2, section_name=dname)
        )

    body = "\n".join(body_lines) if body_lines else "_No cheatsheets found._"

    return toc, body


def apply_template(root: Path, toc: str, body: str, template_path: Path) -> str:
    if template_path.exists():
        template = template_path.read_text(encoding="utf-8")
        content = template.replace(
            "{{last_updated}}", datetime.date.today().isoformat()
        )

        if "{{toc}}" in content:
            content = content.replace("{{toc}}", toc)

        if "{{body}}" in content:
            content = content.replace("{{body}}", body)

        else:
            ## If no {{body}} placeholder, append body after toc
            content = f"{content}\n\n{body}\n"

        return content

    ## Fallback minimal content
    return (
        f"---\n"
        f'title: "Cheatsheets Index"\n'
        f'last_updated: "{datetime.date.today().isoformat()}"\n'
        f"---\n\n"
        f"## Table of Contents <!-- omit in toc -->\n\n{toc}\n\n{body}\n"
    )


def write_index(output_path: Path, content: str) -> None:
    output_path.write_text(content, encoding="utf-8")


def main():
    args = parse_args()

    root = args.root.resolve()
    output = args.output.resolve() if args.output else (root / "INDEX.md")
    toc, body = generate_index_parts(
        root, args.include_dotdirs, args.include_templates, args.include_readme
    )
    template_path = root / ".templates" / "index.md"
    content = apply_template(root, toc, body, template_path)

    write_index(output, content)
    print(f"Wrote index to {output}")


if __name__ == "__main__":
    main()
