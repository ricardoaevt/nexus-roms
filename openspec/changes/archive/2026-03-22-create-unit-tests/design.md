<artifact id="design">

# Design: Backend Unit Testing Strategy

## Context

The Nexus ROMs backend is currently lacking automated unit tests for its core logic. The architecture consists of several internal packages: `renamer` for file organization, `scraper` for metadata retrieval, `orchestrator` for coordination, and `db` for local state. Following **Product Backlog Item 14**, we need to implement a testing framework that ensures 90% coverage on these critical paths.

## Goals / Non-Goals

**Goals:**
- Reach 90% LCOV coverage in `renamer`, `orchestrator`, and `scraper/retry`.
- Isolate units from network and external filesystem state using mocks and temporary directories.
- Ensure thread safety in concurrent components via race detection.
- Fast execution (under 10 seconds for the entire suite).

**Non-Goals:**
- Integration tests with real cloud providers or physical databases.
- UI/Wails testing.
- Manual verification of scraper results against live APIs (use mocks instead).

## Decisions

- **Test Framework**: Use the built-in Go `testing` package for standard behavior and `github.com/stretchr/testify/assert` for fluent assertions.
- **Mocking**: Use `mockery` to generate mocks for internal interfaces (especially the `Scraper` interface). This allows predictable testing of retry logic and error handling without actual network requests.
- **Filesystem Isolation**: Tests in the `renamer` package will use `t.TempDir()` to create unique, isolated environments for file operations.
- **Database Testing**: Use SQLite in-mem mode (`file::memory:?cache=shared`) to test the `db` package logic without side effects.
- **CI/CD Integration**: Tests must follow the `go test -v ./...` structure to be easily integrated into common pipeline runners.

## Risks / Trade-offs

- **Mock Maintenance**: As the `Scraper` or `Renamer` interfaces evolve, mocks will need regeneration.
- **Over-mocking**: High unit test coverage doesn't guarantee integration success. We must ensure mocks reflect current production behavior patterns (e.g., error types).
- **Concurrency Complexity**: Testing race conditions can be flaky if timing is not properly handled in tests. We choose to prioritize race detection over complex stress testing.

</artifact>
