## DevLog
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
