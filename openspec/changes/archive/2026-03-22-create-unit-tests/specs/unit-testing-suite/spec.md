<artifact id="specs/unit-testing-suite">

# Spec: Unit Testing Suite

## ADDED Requirements

### Requirement: Package-Level Coverage
The test suite must provide comprehensive coverage for critical backend modules to ensure logic integrity and prevent regressions.

#### Scenario: Verify Renamer Logic
- **WHEN** template parsing or file renaming operations are invoked in `renamer`
- **THEN** at least 90% of lines within the `renamer` package must be executed by tests.
- **THEN** use `t.TempDir()` for all filesystem manipulations to ensure idempotency.

#### Scenario: Verify Scraper Retry & Mock Providers
- **WHEN** the `scraper` encounters network errors or API rate limits
- **THEN** the `retry` logic must follow exponential backoff.
- **THEN** `screenscraper` and `thegamesdb` providers must be tested using generated mocks to simulate successful and failed API responses.
- **THEN** at least 90% of `scraper/retry` lines must be covered.

#### Scenario: Verify Orchestrator Concurrency
- **WHEN** multiple workers are spawned by the `orchestrator`
- **THEN** race conditions must be absent (validated by `-race` flag).
- **THEN** worker events (progress updates, completion) must be emitted correctly without deadlocks.
- **THEN** at least 90% of `orchestrator` lines must be covered.

#### Scenario: Verify Database Operations
- **WHEN** CRUD operations are performed on sessions or configuration
- **THEN** the `db` package must validate state changes using an in-memory SQLite database.
- **THEN** credentials and session persistence must be verified via unit tests.

### Requirement: Performance and Isolation
Tests must be efficient and self-contained to encourage frequent execution.

#### Scenario: Test Suite Execution Time
- **WHEN** `go test ./...` is executed locally or in CI
- **THEN** the total execution time for all unit tests must not exceed 10 seconds.
- **THEN** no external network calls or physical database dependencies are permitted.

</artifact>
