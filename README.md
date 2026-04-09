# Seedbank

TUI fake data generator. Browse generators, pick fields, preview output, export to file. Built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

Generate realistic, internally-consistent seed data for databases, APIs, and testing — from the terminal.

## Quick Install

Recommended (installs to `~/.local/bin`):

```bash
curl -fsSL https://raw.githubusercontent.com/LFroesch/seedbank/main/install.sh | bash
```

Or download a binary from [GitHub Releases](https://github.com/LFroesch/seedbank/releases).

Or install with Go:

```bash
go install github.com/LFroesch/seedbank@latest
```

Or build from source:

```bash
make install
```

Command:

```bash
seedbank
```
## Install

```bash
go install github.com/LFroesch/seedbank@latest
```

Or build from source:

```bash
cd apps/seedbank
make install
```

## Usage

```bash
seedbank              # interactive TUI
echo "" | seedbank    # pipe mode (CLI, no TUI)
```

### TUI Flow

1. Pick a generator (Names, Addresses, Person, Product, etc.)
2. Toggle fields with `space`
3. Set count with `+/-`
4. Preview output, switch format with `f`
5. Export with `e` or copy with `y`

## Generators

| Generator | What it produces |
|-----------|-----------------|
| Names | Gendered first/last names with prefixes |
| Lorem Ipsum | Words, sentences, paragraphs, titles |
| Addresses | Coherent US addresses (state/city/zip cross-referenced) |
| Emails | Addresses derived from gendered name pools |
| Phone Numbers | US numbers with state-matched area codes |
| Photos | Placeholder image URLs (picsum.photos) |
| Numbers | IDs, ints, floats, currency, booleans |
| Companies | Company names, departments, job titles |
| Dates | Internally consistent dates/timestamps/ages |
| Colors | Color names mapped to real RGB/hex values |
| Person (Linked) | Fully coherent person records across all fields |
| Product (Linked) | Products with category-driven pricing |
| Mix | Combine fields from multiple generators |

## Output Formats

JSON, JSONL, CSV, Markdown table, SQL INSERT statements.

## Keybindings

| Key | Action |
|-----|--------|
| `j/k`, `up/down` | Navigate |
| `enter` | Select generator / confirm |
| `space` | Toggle field |
| `r` | Re-roll data (new seed) |
| `+/-` | Adjust record count |
| `f` | Cycle output format |
| `v` | Toggle pretty view |
| `e` | Export to file |
| `y` | Copy to clipboard |
| `?` | Help |
| `esc` | Back |
| `q` | Quit |

## Configuration

Config file: `~/.config/seedbank/config.json`

## Data Quality

- 360+ gendered first names, 200+ last names
- 20 US states with real cities, zip prefixes, area codes — all cross-referenced
- Color names mapped to actual RGB/hex values
- RFC 2606 safe domains (`.example.com`) for company/work emails
- Date fields derived from single timestamp for internal consistency

## Platform Support

Linux, macOS, WSL. Tested on Linux/WSL.

## License

[AGPL-3.0](LICENSE)
