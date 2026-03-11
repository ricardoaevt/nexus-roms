## ADDED Requirements

### Requirement: Deep Directory Traversal
The system SHALL traverse local and network directories to discover supported ROM extensions directly or contained inside compressed archives (.zip, .rar, .7z).

#### Scenario: Scanning mixed files
- **WHEN** point to a folder with isolated roms and nested archives
- **THEN** the scanner builds a unified list of processable files mapping their relative paths

### Requirement: Advanced Template Engine for Renaming
The system SHALL provide a dynamic string replacement engine for mass-renaming that supports tags like `{Name}`, `{Region}`, `{Languages}`, `{Hash}`, removing invalid OS characters.

#### Scenario: Applying the Default Template
- **WHEN** the system resolves a metadata match with name 'Super Mario' and region 'USA'
- **THEN** the `{Name} ({Region})` template yields 'Super Mario (USA)' replacing forbidden characters
