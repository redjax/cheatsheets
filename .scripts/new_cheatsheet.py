
import sys
import datetime
import argparse
from pathlib import Path

def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Create a new cheatsheet from a template.",
        formatter_class=argparse.RawTextHelpFormatter,
        epilog=(
            "Examples:\n"
            "  python new_cheatsheet.py apps my-new-app\n"
            "  python new_cheatsheet.py --title 'My App' --description 'Cheatsheet' --tags cli,app\n"
            "  python new_cheatsheet.py  # prompts for inputs\n"
        ),
    )

    # Optional positional arguments for category and name
    parser.add_argument("category", nargs="?", help="Cheatsheet category (e.g., 'apps', 'languages').")
    parser.add_argument("name", nargs="?", help="New cheatsheet name (without extension).")

    # Optional metadata flags (skips prompts if provided)
    parser.add_argument("--title", dest="title", help="Title to use in the cheatsheet.")
    parser.add_argument("--description", dest="description", help="Description to use in the cheatsheet.")
    parser.add_argument("--tags", dest="tags", help="Comma-separated additional tags (e.g., cli,tool).")

    args = parser.parse_args()
    
    return args

def create_new_cheatsheet(category: str, name: str, title: str | None = None, description: str | None = None, tags_input: str | None = None):
    """Creates a new cheatsheet from a template.

    Params:
        category (str): The category of the cheatsheet (e.g., 'apps', 'languages').
        name (str): The name of the new cheatsheet (without extension).
        title (str | None): Optional title; if not provided, will be prompted.
        description (str | None): Optional description; if not provided, will be prompted.
        tags_input (str | None): Optional comma-separated tags; if not provided, will be prompted.
    """
    cheatsheet_root = Path(__file__).parent
    template_path = cheatsheet_root / ".templates" / f"{category}.md"
    target_dir = cheatsheet_root / category
    target_file = target_dir / f"{name}.md"

    if not template_path.exists():
        print(f"Error: Template not found at {template_path}")
        sys.exit(1)

    if target_file.exists():
        print(f"Error: Cheatsheet already exists at {target_file}")
        sys.exit(1)

    ## Gather inputs (prompt if missing)
    print(f"Creating new cheatsheet: {name} in {category}")
    if title is None:
        title = input("Enter title: ")
    if description is None:
        description = input("Enter description: ")
    if tags_input is None:
        tags_input = input("Enter additional tags (comma-separated): ")

    ## Prepare template replacements
    replacements = {
        "{{title}}": title,
        "{{description}}": description,
        "{{last_updated}}": datetime.date.today().isoformat(),
        "{{tags}}": ", ".join(f'"{tag.strip()}"' for tag in tags_input.split(",") if tag.strip())
    }

    ## Read template and replace placeholders
    template_content = template_path.read_text(encoding="utf-8")
    new_content = template_content
    for placeholder, value in replacements.items():
        new_content = new_content.replace(placeholder, value)

    ## Create directory and write the new file
    target_dir.mkdir(exist_ok=True)
    target_file.write_text(new_content, encoding="utf-8")

    print(f"\nSuccessfully created {target_file}")

def main():
    args = parse_args()

    category = args.category or input("Enter category: ")
    name = args.name or input("Enter name: ")

    create_new_cheatsheet(
        category=category,
        name=name,
        title=args.title,
        description=args.description,
        tags_input=args.tags,
    )

if __name__ == "__main__":
    main()
