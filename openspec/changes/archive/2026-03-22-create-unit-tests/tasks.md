<artifact id="tasks">

# Tasks: Implement Unit Testing Suite for Backend

## 1. Infrastructure Setup

- [x] 1.1 Add `github.com/stretchr/testify/assert` to `go.mod`.
- [x] 1.2 Initialize `mockery` for the `scraper.Scraper` interface.
- [x] 1.3 Add test execution scripts or `Makefile` targets for coverage and race detection.

## 2. Renamer Testing

- [x] 2.1 Implement unit tests for `internal/renamer` (template parsing).
- [x] 2.2 Add tests for filesystem operations using `t.TempDir()`.
- [x] 2.3 Verify 90% LCOV coverage in `renamer`.

## 3. Scraper & Retry Testing

- [x] 3.1 Implement unit tests for exponential backoff in `internal/scraper/retry.go`.
- [x] 3.2 Mock `screenscraper` and `thegamesdb` providers for isolated testing.
- [x] 3.3 Verify 90% LCOV coverage in `scraper/retry`.
- [x] 4.2 Use `go test -race` to validate concurrency logic.
- [x] 4.3 Verify 90% LCOV coverage in `orchestrator`.
- [x] 6.1 Ensure total test execution time is under 10 seconds.
- [x] 6.2 Confirm all backend tests pass locally and follow Go naming standards.

</artifact>
