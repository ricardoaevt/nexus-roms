## ADDED Requirements

### Requirement: Concurrent File Hashing and Identification
The system SHALL orchestrate file identification using a worker pool, calculating streams (MD5, SHA1, CRC32) incrementally without fully extracting zipped containers.

#### Scenario: Deep-scanning a ROM within a ZIP
- **WHEN** the scanner encounters a .zip file containing a ROM
- **THEN** it calculates the hash in-stream without extracting the file to disk

### Requirement: Session Pause and Resume Control
The system SHALL allow background scraping tasks to be paused, stopped, and resumed gracefully, preserving exact progress in a persistent state.

#### Scenario: Reinstating a broken session
- **WHEN** the application starts and discovers an unfinished session for the selected root path
- **THEN** a modal is displayed asking whether to resume progress from where it left off or restart completely
