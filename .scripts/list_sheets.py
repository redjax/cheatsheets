from pathlib import Path
import sys
import json
from dataclasses import dataclass, field


repo_root = Path(__file__).parent.parent
cheatsheets_dir = repo_root / "cheatsheets"


def path_from_root(path: Path, root: str) -> Path:
    path = Path(path).resolve()
    parts = path.parts

    try:
        ## Find the index of the root directory in the path
        root_index = parts.index(root)
    except ValueError:
        raise ValueError(f"Root '{root}' not found in path '{path}'")

    ## Build the relative path from the part after root
    relative_parts = parts[root_index + 1 :]

    return Path(*relative_parts)


def main():
    if not cheatsheets_dir.exists():
        raise FileNotFoundError(
            f"Could not find cheatsheets at path: {cheatsheets_dir}"
        )

    print(f"Cheatsheets path: {cheatsheets_dir}")

    cheatsheets: list[Path] = []

    for p in cheatsheets_dir.rglob("**/*"):
        if p.is_file():
            if p.suffix == ".md":
                cheatsheets.append(p)

    print(
        f"\nFound {len(cheatsheets)} cheatsheet{'s' if len(cheatsheets) > 1 else ''}\n"
    )

    for sheet in cheatsheets:
        _path = path_from_root(sheet, root="cheatsheets")
        print(f"  - {sheet.name} ({_path})")

    print()


if __name__ == "__main__":
    try:
        main()
    except Exception as exc:
        print(f"[ERROR] Failed listing cheatsheets: {exc}")
        sys.exit(1)
