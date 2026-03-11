## ADDED Requirements

### Requirement: Reactive Desktop Dashboard
The system SHALL render a Wails-bridged responsive dashboard using Svelte, providing real-time log outputs, scraping statistics, and a match preview table.

#### Scenario: Displaying real time scanning progress
- **WHEN** the backend orchestrator emits a 'progress' event
- **THEN** the frontend updates the live stream file display and increments the progress bar without blocking the UI

### Requirement: Match Approval and Execution
The UI SHALL provide a table previewing the original filename next to the proposed template output, allowing the user to select manually which files to rename.

#### Scenario: Committing Rename Operations
- **WHEN** the user selects rows in the live identification table and clicks 'Rename Selected'
- **THEN** the UI delegates the physical file rename to the Go backend and marks the rows as 'Renamed' upon success
