---
name: db-migration
description: Use this skill when adding or changing database tables in finance-data-backend, including SQL migrations, `configs/db.yaml` `db.gen` adjustments, running migration and model generation commands, and verifying generated GORM models.
---

# DB Migration

Use this skill for database migrations and GORM model generation in `finance-data-backend`.

## Step 1

Add a SQL migration file under `internal/db/migrations`.

- The filename prefix must use the `yyyymmddhh` timestamp format
- Prefer the filename format `yyyymmddhh-description.sql`
- If you need multiple migrations within the same hour, append a letter after the timestamp before the description
- Example: `2026041809-secret.sql`
- The file must contain:
  - `-- +migrate Up`
  - `-- +migrate Down`

Also follow the database naming convention from section 13 of `CLAUDE.md`:

- Table names must use the `t_` prefix and must be singular
- Index names must not use the `t_` prefix
- Column names must not use reserved database keywords
- Common column names to avoid: `order`, `group`, `select`, `table`, `index`, `key`, `type`, `status`, `value`

Example:

```sql
-- Good
CREATE TABLE t_order (...);
CREATE INDEX idx_user_name ON t_user(name);
order_type VARCHAR(50);

-- Bad
CREATE TABLE orders (...);
CREATE INDEX idx_t_user_name ON t_user(name);
type VARCHAR(50);
```

## Step 2

Check and update the `db.gen` section in `configs/db.yaml` if needed:

- `customization`
- `association`
- Any other settings related to model generation

Only change these settings when the current migration actually requires it. Do not make unrelated config changes.

## Step 3

Run the following commands in this exact order:

```bash
make gen-up && make gen-down && make gen-up && make gen-model
```

Purpose:

- Verify the migration can be applied
- Verify rollback works correctly
- Verify applying the migration again still works
- Generate the latest GORM models

If you hit a database connection issue, tell the user clearly:

- Check and configure `db.dsn` in `configs/db.yaml`
- Make sure the current environment can connect to the target database before retrying the commands above

## Step 4

Check whether `internal/data/model` contains updated or newly generated models:

- If model files were updated or generated, the migration flow is complete
- If nothing changed, investigate the reason, for example:
  - The migration did not actually run successfully
  - The table name or schema does not match generation expectations
  - `db.gen` settings in `configs/db.yaml` affected generation
  - The connected database is not the expected one

## Step 5

After all migration, generation, and related code changes are finished, run:

```bash
make lint
```

Purpose:

- Apply the repository lint workflow after all changes are in place
- Surface formatting, style, and static analysis issues introduced by the migration work
- Fix any lint errors before considering the task complete

If `make lint` reports issues, continue fixing them until the lint command passes, unless the failure is clearly caused by unrelated pre-existing problems. In that case, tell the user exactly what failed and why it appears unrelated.

Do not bypass lint by adding `//nolint`, weakening lint rules, or suppressing checks just to make the command pass. Prefer real code fixes and refactors that address the reported issue at the source.

## Notes

- Do not hand-write `internal/data/model/*.gen.go`
- Models must be produced through the generation commands in `Makefile`
- After finishing the migration, check these paths first:
  - `internal/db/migrations`
  - `internal/data/model`
  - `configs/db.yaml`
