## ADDED Requirements

### Requirement: Manejo de colisiones en renombramiento
El sistema DEBE detectar colisiones de nombres cuando el archivo destino renombrado ya existe en la ubicación destino.
El sistema DEBE mover el archivo duplicado a un directorio especial llamado `duplicados` ubicado en relación al destino.
El sistema DEBE recrear la estructura de directorios relativa original del archivo dentro del directorio `duplicados` antes de moverlo, para preservar el contexto de dónde venía.

#### Scenario: Colisión de nombres detectada
- **WHEN** el archivo A va a ser renombrado a "Juego.zip" pero "Juego.zip" ya existe en la raíz destino
- **THEN** el archivo A es movido a `[destino]/duplicados/.../Juego.zip`, donde `...` representa su jerarquía de rutas original relativa
