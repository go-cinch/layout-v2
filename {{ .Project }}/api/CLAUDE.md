# API Proto Naming Conventions

## Message Naming

All message names must be singular (no trailing "s"). Pluralization is handled by the `repeated` keyword in proto3 fields.

Examples:
- `Stock` not `Stocks`
- `Sector` not `Sectors`
- `StockResult` not `StockResults`
- `SectorResult` not `SectorResults`

## Rationale

Proto3 uses `repeated` to indicate collections, so message names should not carry plural semantics. This keeps generated code cleaner and avoids double-pluralization (e.g., `repeated Stocks` → `[]Stocks` in Go would read as "list of list of stocks").
