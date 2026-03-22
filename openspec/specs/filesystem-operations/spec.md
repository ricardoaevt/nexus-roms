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

### Requirement: Filtrado de archivos comprimidos complejos
El sistema DEBE descartar los archivos `.zip`, `.rar` y `.7z` que contengan múltiples ROMs u archivos no relacionados (ej. romsets/recopilaciones).
El sistema NO DEBE descartar los archivos que, aunque contengan múltiples archivos, estén relacionados a un título válido único (ej. DLCs, parches, expansiones).

#### Scenario: Archivo con recopilación no relacionada
- **WHEN** el sistema escanea un archivo comprimido que contiene múltiples archivos ROM con nombres base distintos
- **THEN** el sistema descarta el archivo completo sin extraer ni procesar las ROMs

#### Scenario: Archivo con DLCs o discos múltiples
- **WHEN** el sistema escanea un archivo comprimido que contiene múltiples archivos que comparten un nombre base muy similar o idéntico (ej. "Juego (Disco 1)", "Juego (Disco 2)")
- **THEN** el sistema considera el archivo válido y procede a procesarlo

### Requirement: Manejo de colisiones en renombramiento
El sistema DEBE detectar colisiones de nombres cuando el archivo destino renombrado ya existe en la ubicación destino.
El sistema DEBE mover el archivo duplicado a un directorio especial llamado `duplicados` ubicado en relación al destino.
El sistema DEBE recrear la estructura de directorios relativa original del archivo dentro del directorio `duplicados` antes de moverlo, para preservar el contexto de dónde venía.

#### Scenario: Colisión de nombres detectada
- **WHEN** el archivo A va a ser renombrado a "Juego.zip" pero "Juego.zip" ya existe en la raíz destino
- **THEN** el archivo A es movido a `[destino]/duplicados/.../Juego.zip`, donde `...` representa su jerarquía de rutas original relativa

### Requirement: Reporte detallado de fallos en renombramiento
El sistema DEBE recolectar y almacenar las razones de cualquier fallo que ocurra al intentar renombrar o mover un archivo ROM.
El sistema DEBE mostrar una tabla de resumen o modal al final del proceso por lotes detallando qué archivos fallaron y por qué.

#### Scenario: El archivo está bloqueado por el sistema operativo
- **WHEN** el sistema intenta renombrar un archivo que está en uso/bloqueado por el SO, causando un fallo
- **THEN** el sistema registra la ruta del archivo y la razón del error, y continúa con el siguiente archivo

#### Scenario: El renombramiento por lotes finaliza con errores
- **WHEN** el proceso de renombramiento por lotes termina y se registraron errores
- **THEN** la interfaz presenta un reporte mostrando la lista de archivos que no se pudieron renombrar y sus correspondientes mensajes de error
