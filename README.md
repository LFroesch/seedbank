# Seedbank

TUI fake data generator. Browse generators, pick fields, preview output, and export deterministic fixture data from the terminal. Built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

Generate realistic, internally-consistent seed data for databases, APIs, and testing from the terminal.

The TUI uses the same full-width shell pattern as `sb`: active step tabs in the header, panel-based content, transient status line, and a footer with mode-aware key hints. On wider terminals, workflow screens use equal-size side-by-side panels; on tighter terminals they collapse into a stacked compact layout while keeping the header and footer visible.

## Quick Install

Supported platforms: Linux and macOS. On Windows, use WSL.

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
echo "" | seedbank    # pipe mode (CLI, no TUI)
seedbank --gen person --count 25 --fmt json --out fixtures/users.json
```

### TUI Flow

1. Pick a generator (Names, Addresses, Person, Product, etc.)
2. Toggle fields with `space`
3. Set count with `+/-`
4. Preview output, switch format with `f`
5. Export with `e` or copy with `c`

Coherent record builders are best when you want one ready-to-export record with internally consistent values out of the box. Field sources are best when you want to mix and match specific columns.

## Generators

### Coherent Record Builders

| Generator | What it produces |
|-----------|-----------------|
| Person | Fully coherent person records across identity, contact, address, and company fields |
| Products | Products with category-driven pricing, stock, and ratings |

### Field Sources

| Generator | What it produces |
|-----------|-----------------|
| Names | Gendered first/last names with prefixes |
| Lorem Ipsum | Words, sentences, paragraphs, titles |
| Addresses | Coherent US addresses (state/city/zip cross-referenced) |
| Emails | Addresses derived from fake name pools |
| Phone Numbers | US numbers with state-matched area codes |
| Photos | Deterministic placeholder image asset URLs and metadata |
| Numbers | IDs, ints, bigints, decimals, currency amounts/codes, booleans |
| Identifiers | UUID variants and compact IDs |
| Network | IPv4, IPv6, MAC, hostname, URL, and user-agent values |
| Companies | Company names, departments, job titles |
| Dates | Internally consistent dates/timestamps/ages |
| Colors | Color names mapped to real RGB/hex values |

`Custom Mix` is a TUI-only compose mode. Press `m` from Browse to combine fields from multiple field-source generators into one record shape.

## Output Formats

JSON, JSONL, CSV, Markdown table, SQL INSERT statements.

## Keybindings

| Key | Action |
|-----|--------|
| `j/k`, `up/down` | Navigate |
| `pgup/pgdn`, `g/G` | Faster scroll / jump to top or bottom |
| `enter` | Select generator / confirm |
| `m` | Open custom mix builder from Browse |
| `space` | Toggle field |
| `tab` | Switch active pane on split views |
| `1-6` | Jump between workflow steps |
| `r` | Re-roll data (new seed) |
| `+/-` | Adjust record count |
| `f` | Cycle output format |
| `v` | Toggle pretty view |
| `e` | Export to file |
| `c` | Copy to clipboard |
| `?` | Toggle help |
| `esc` | Back |
| `q` | Back / quit from Browse |

## Configuration

Config file: `~/.config/seedbank/config.json`

By default, relative TUI exports write into the directory where you launched `seedbank`. Set `output_dir` in config if you want relative exports rooted somewhere else.

## CLI

Use `--gen` to skip the TUI and write to stdout or a file directly.

```bash
seedbank --list
seedbank --gen person --count 10 --fmt json
seedbank --gen identifiers --count 10 --fields uuid,short_id --fmt csv
seedbank --gen network --count 20 --fmt jsonl
seedbank --gen phone --fields phone,state --fmt csv --out fixtures/phones.csv
seedbank --gen products --count 50 --fmt sql --table products --out db/seed/products.sql
seedbank --gen person --count 3 --fmt json --seed 42
```

Notes:
- `--list` now groups generators into coherent record builders and field sources.
- `--out -` writes to stdout. Any other path writes to that file and creates parent directories if needed.
- Relative `--out` paths are resolved from the directory where you run `seedbank`.
- `--seed` makes outputs reproducible across runs, including age/date-derived fields.
- `--schema <file.sql>` is a preview feature: it parses a `CREATE TABLE` statement and generates rows with exact column names using type/name heuristics.

### Schema Mode (`--schema`)

Schema mode is a preview feature for generating seed data that matches a table definition directly. It is useful for simple single-table `CREATE TABLE` inputs, but it is still heuristic-driven rather than a full schema seeding workflow.

```bash
seedbank --schema db/schema/users.sql --count 100 --fmt json
seedbank --schema db/schema/orders.sql --count 500 --fmt sql --out db/seed/orders.sql
```

Heuristic mapping examples:
- `UUID`, `id`, `*_id` -> `identifiers.uuid` (or numeric IDs for integer types)
- `VARCHAR/TEXT` with name hints (`email`, `first_name`, `last_name`, etc.) -> matching field generators
- `INT/BIGINT` -> `numbers.integer` / `numbers.bigint`
- `DECIMAL/NUMERIC/FLOAT` -> `numbers.decimal` or `numbers.currency_amount` for price-like names
- `BOOLEAN` -> `numbers.boolean`
- `TIMESTAMP/DATE/TIME` -> `dates.datetime` / `dates.date` / `dates.time`
- `ip_address`, `user_agent`, `url`, `hostname`, `mac` -> `network.*`

### Placeholder Image Assets

The `Photos` generator is for deterministic placeholder image asset data, not text-to-image generation. It emits stable placeholder URLs plus dimensions, aspect ratio, a category tag, a deterministic seed, and a simple asset label you can use in fixtures or mock content pipelines.

## Data Quality

- 360+ gendered first names, 200+ last names
- 20 US states with real cities, zip prefixes, area codes — all cross-referenced
- Color names mapped to actual RGB/hex values
- RFC 2606 safe domains (`.example.com`) for company/work emails
- Date fields derived from single timestamp for internal consistency
- Age-derived fields use a fixed reference date to keep seeded output stable

## Platform Support

Linux, macOS, WSL. Tested on Linux/WSL.

## License

[AGPL-3.0](LICENSE)
