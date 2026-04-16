## DevLog
### 2026-04-14: `--schema` flag implementation
Implemented `--schema` CLI mode that parses a `CREATE TABLE` statement from SQL, infers column mappings from SQL types + column-name hints, and generates output with exact schema column names. Added support for key backend mappings including `UUID`/`*_id` via `Identifiers`, network-centric columns (`ip`, `url`, `user_agent`, `mac`, `hostname`) via `Network`, and numeric/date/boolean mappings through `Numbers`/`Dates`. Added parser/mapping unit tests and preserved existing `--fmt`/`--out` flow for schema mode.
Files: main.go, schema.go, schema_test.go, main_test.go, README.md, WORK.md.

### 2026-04-14: Generator breadth + determinism pass
Expanded schema-friendly field coverage by adding two new field-source generators: `Identifiers` (UUID variants and compact IDs) and `Network` (IPv4/IPv6/MAC/hostname/URL/user-agent). Upgraded `Numbers` to include `bigint`, `decimal`, `currency_amount`, and `currency_code` while keeping legacy formatted `currency`. Also removed wall-clock drift from seeded outputs by calculating age-derived fields against a fixed reference date, and added regression tests to enforce DOB/date age consistency.
Files: internal/generator/registry.go, internal/generator/identifier.go, internal/generator/network.go, internal/generator/number.go, internal/generator/person.go, internal/generator/date.go, internal/generator/time_ref.go, internal/generator/generator_test.go, README.md, WORK.md.

### 2026-04-14: Taxonomy and CLI export pass
Reduced generator overlap confusion by classifying generators as either coherent record builders or field sources, then surfaced that grouping in the TUI, generator descriptions, and `--list` output. Also polished CLI export behavior with config-backed defaults, alias-based generator lookup, and a new `--out` flag that writes to files directly and creates parent directories when needed.
Files: main.go, main_test.go, internal/generator/registry.go, internal/generator/*.go, update.go, view.go, README.md, WORK.md.

### 2026-04-14: UX polish pass
Polished the TUI workflow around shorter terminals and faster navigation: added numbered step tabs (`1-6`), pane focus switching with `tab`, independent scrolling for overflowed side panels, a denser `./sb`-style footer, and safer export behavior for relative paths plus directory creation on write. Also tightened generator guidance in the UI, clarified the photo generator as placeholder/WIP output, and added regression tests for step jumping, footer rendering, output-path resolution, and the new photo metadata.
Files: model.go, update.go, view.go, update_test.go, view_test.go, internal/config/config.go, internal/generator/photo.go, internal/generator/generator_test.go, README.md, WORK.md.

### 2026-04-13: Shell layout bugfix pass
Fixed the refreshed Seedbank shell so it actually fits the terminal: panel rendering now respects Lip Gloss border sizing, clips panel content to the available interior space, preserves the header/footer in both wide and stacked layouts, and keeps the left/right workflow panels balanced. Also fixed contextual help to return to the previous screen instead of dumping back to Browse, and aligned export/back navigation with the footer hints. Added regression tests for panel sizing and full-view height so the shell does not silently overflow again.
Files: view.go, update.go, model.go, view_test.go, WORK.md.

### 2026-04-13: TUI shell refresh
Reworked the Seedbank interface to match the sibling `sb` app more closely: full-width header/footer shell, active step tabs, separator lines, two-panel layouts for generator/field/preview flows, and a persistent footer with mode-aware key hints. Followed that up with equal-width responsive panels plus a compact stacked fallback so the shell stays visible on small terminals. Also aligned docs with the actual copy key (`c`) and noted remaining v1 polish work.
Files: view.go, README.md, WORK.md.

### 2026-03-23: Doc suite refresh
Updated README to scout standard. Added LICENSE (AGPL-3.0). Updated WORK.md with feature ideas.

### 2026-03-20: Audit
Code reviewed. Clean architecture, good data quality (gendered names, cross-referenced locations, RFC-safe domains). Has tests. UI functional but flat — no full-width header, no responsive hints, no left/right status bar. v1 generators are solid, biggest interview-impact addition is `--schema` flag.

### 2026-03-18: Data coherence overhaul
Rewrote all 12 generators for internal consistency:
- Gendered names (180+ M/F each) with gender-matched prefixes
- 20 US states with real cities, zip prefixes, area codes — all cross-referenced
- Color names mapped to actual RGB/hex values
- Date fields (year/age/date) all derived from single timestamp
- Product descriptions use 6 varied templates (no repeated adjective)
- Company/work domains use `.example.com` (RFC 2606 safe)
- Images switched to picsum.photos (deterministic, reliable)
- New `data.go` shared layer with `pickGendered()`, `pickLocation()`, `pickStreet()`, `genUUID()`
- Expanded pools: ~200 last names, ~360 first names, 64 streets, 40 company prefixes
Files: internal/generator/data.go (new), names.go, address.go, person.go, email.go, phone.go, company.go, product.go, date.go, color.go, photo.go, number.go.

### 2026-03-17: Daily driver features
Clipboard copy, pipe/CLI mode, mix mode, email linking, SQL table names, tests.

### 2026-03-17: Linked generators + view toggle
Person (Linked) and Product (Linked) generators. Pretty view toggle with `v` key.

### 2026-03-17: Initial scaffolding
10 generators, 5 output formats, full TUI flow.
