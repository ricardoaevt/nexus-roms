## ADDED Requirements

### Requirement: Indicadores de estado visuales no bloqueantes
El sistema DEBE mostrar un indicador de carga visual (mata-angustias) durante las operaciones prolongadas (ej. escaneo, scraping).
El estado del indicador de carga NO DEBE bloquear ni deshabilitar los controles de "Pausar" y "Detener" en la interfaz de usuario.

#### Scenario: El usuario inicia un rasgado (scraping) por lotes largo
- **WHEN** inicia una operación de scraping por lotes
- **THEN** la interfaz muestra una capa superpuesta con el indicador de carga activa sobre la configuración y los inputs

#### Scenario: El usuario intenta pausar el scraping
- **WHEN** la operación de scraping está en ejecución y el indicador de carga es visible
- **THEN** el usuario aún puede hacer clic en los botones de Pausa o Detener gracias a su mayor z-index y estado activo
