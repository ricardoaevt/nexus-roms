## ADDED Requirements

### Requirement: Reporte detallado de fallos en renombramiento
El sistema DEBE recolectar y almacenar las razones de cualquier fallo que ocurra al intentar renombrar o mover un archivo ROM.
El sistema DEBE mostrar una tabla de resumen o modal al final del proceso por lotes detallando qué archivos fallaron y por qué.

#### Scenario: El archivo está bloqueado por el sistema operativo
- **WHEN** el sistema intenta renombrar un archivo que está en uso/bloqueado por el SO, causando un fallo
- **THEN** el sistema registra la ruta del archivo y la razón del error, y continúa con el siguiente archivo

#### Scenario: El renombramiento por lotes finaliza con errores
- **WHEN** el proceso de renombramiento por lotes termina y se registraron errores
- **THEN** la interfaz presenta un reporte mostrando la lista de archivos que no se pudieron renombrar y sus correspondientes mensajes de error
