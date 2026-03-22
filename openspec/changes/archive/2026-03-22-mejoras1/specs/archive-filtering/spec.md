## ADDED Requirements

### Requirement: Filtrado de archivos comprimidos complejos
El sistema DEBE descartar los archivos `.zip`, `.rar` y `.7z` que contengan múltiples ROMs u archivos no relacionados (ej. romsets/recopilaciones).
El sistema NO DEBE descartar los archivos que, aunque contengan múltiples archivos, estén relacionados a un título válido único (ej. DLCs, parches, expansiones).

#### Scenario: Archivo con recopilación no relacionada
- **WHEN** el sistema escanea un archivo comprimido que contiene múltiples archivos ROM con nombres base distintos
- **THEN** el sistema descarta el archivo completo sin extraer ni procesar las ROMs

#### Scenario: Archivo con DLCs o discos múltiples
- **WHEN** el sistema escanea un archivo comprimido que contiene múltiples archivos que comparten un nombre base muy similar o idéntico (ej. "Juego (Disco 1)", "Juego (Disco 2)")
- **THEN** el sistema considera el archivo válido y procede a procesarlo
