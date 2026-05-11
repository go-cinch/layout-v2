# Internal Layering Guidelines

## Service Layer

- `service` package only handles request parameter validation, request/response conversion, and routing into different business branches
- `service` package must not directly operate the database
- `service` package should not contain business orchestration that belongs in `biz`

## Biz Layer

- `biz` package owns business logic and orchestration across one or more repos or external components
- `biz` package decides which repo paths to use for a scenario
- `biz` package does not need to repeatedly validate whether request parameters are legal once `service` has accepted them

## Data Layer

- `data` package defines repo implementations and directly operates the database
- `data` package should focus on persistence concerns and SQL/query assembly
- `data` package does not need to repeatedly validate whether upstream request parameters are legal once they have entered the repo

## Test Placement

- tests for the `biz` package should be placed under `internal/tests/biz`
- tests for the `data` package should be placed under `internal/tests/data`
- tests for the `service` package should be placed under `internal/tests/service`
