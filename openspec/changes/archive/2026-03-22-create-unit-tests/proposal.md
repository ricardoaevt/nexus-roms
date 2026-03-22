<artifact id="proposal">

# Proposal: Implement Unit Testing Suite for Backend Components

## Why

Currently, the backend components lack a robust suite of unit tests, making the code harder to maintain and refactor safely. There is a need to ensure the integrity of the core logic, especially for the renamer, scraper, and orchestrator modules. This proposal addresses **Product Backlog Item 14** (crear test unitarios) to reach a 90% code coverage in critical backend packages.

## What Changes

We will implement a comprehensive suite of unit tests using the Go standard `testing` package and `testify/assert` for improved readability. This includes mocking external APIs for the screen-scraper and TheGamesDB providers, as well as simulating concurrent scenarios in the orchestrator.

## Capabilities

### New Capabilities
- `unit-testing-suite`: A robust suite of automated tests for the Nexus ROMs backend, ensuring that core modules meet specific reliability and performance criteria.

### Modified Capabilities
- (None - this change focuses on adding test coverage to existing logic without changing their core functional requirements)

## Impact

- **internal/renamer**: Tests for template parsing and file renaming logic.
- **internal/scraper**: Unit tests for `retry.go` (backoff logic), mocks for `screenscraper.go` and `thegamesdb.go`.
- **internal/orchestrator**: Concurrency tests and worker management validation.
- **internal/db**: CRUD validation using SQLite in memory.
- **Development Workflow**: A new testing stage in the CI/CD pipeline or local pre-commit hooks to ensure quality standards.

</artifact>
