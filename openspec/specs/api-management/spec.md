## ADDED Requirements

### Requirement: Centralized API Authentication and Configuration
The system SHALL provide a centralized management interface that allows secure configuration of multiple metadata APIs (ScreenScraper, TheGamesDB). Passwords and API Keys MUST be stored using AES-256-GCM encryption.

#### Scenario: User saves ScreenScraper credentials
- **WHEN** the user inputs an ID and Password and clicks save
- **THEN** the system encrypts the password with AES-256-GCM and persists it to the SQLite database

#### Scenario: User enables API search strategies
- **WHEN** the user toggles the Hash or Name strategy for an API
- **THEN** the system updates the configuration state and makes that Scraper eligible for those query types
