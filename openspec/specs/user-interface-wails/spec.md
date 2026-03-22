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

### Requirement: Indicadores de estado visuales no bloqueantes
El sistema DEBE mostrar un indicador de carga visual (mata-angustias) durante las operaciones prolongadas (ej. escaneo, scraping).
El estado del indicador de carga NO DEBE bloquear ni deshabilitar los controles de "Pausar" y "Detener" en la interfaz de usuario.

#### Scenario: El usuario inicia un rasgado (scraping) por lotes largo
- **WHEN** inicia una operación de scraping por lotes
- **THEN** la interfaz muestra una capa superpuesta con el indicador de carga activa sobre la configuración y los inputs

#### Scenario: El usuario intenta pausar el scraping
- **WHEN** la operación de scraping está en ejecución y el indicador de carga es visible
- **THEN** el usuario aún puede hacer clic en los botones de Pausa o Detener gracias a su mayor z-index y estado activo
