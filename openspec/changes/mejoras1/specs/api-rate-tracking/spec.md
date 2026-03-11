## ADDED Requirements

### Requirement: Recuento persistente de peticiones de API
El sistema DEBE realizar un seguimiento del número de peticiones HTTP realizadas a cada proveedor de API de metadatos (ej. ScreenScraper).
El sistema DEBE persistir este conteo de uso entre reinicios de la aplicación.
El sistema DEBE reiniciar automáticamente el contador a 1 cuando se realiza una petición en un mes calendario nuevo.

#### Scenario: Petición en el mismo mes
- **WHEN** la aplicación realiza una petición a la API y el mes calendario actual coincide con el mes guardado
- **THEN** el sistema incrementa el contador guardado en 1 y lo persiste

#### Scenario: Petición en un nuevo mes
- **WHEN** la aplicación realiza una petición a la API y el mes calendario actual es diferente al mes guardado
- **THEN** el sistema reinicia el contador a 1, actualiza el mes guardado al mes actual, y persiste los nuevos valores
