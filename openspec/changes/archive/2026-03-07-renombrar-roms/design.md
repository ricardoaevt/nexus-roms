## Contexto

El proyecto ROM Renamer tiene como objetivo resolver el problema de las colecciones de juegos retro desorganizadas, donde los archivos tienen nombres técnicos o inconsistentes. El sistema automatizará la identificación de estos archivos calculando sus hashes y consultando APIs de metadatos (como ScreenScraper) para obtener los títulos oficiales. Se trata de una nueva aplicación de escritorio independiente construida con **Go** y **Wails**, que requiere un puente eficiente entre las operaciones de bajo nivel del sistema de archivos y una interfaz de usuario reactiva.

## Metas / No-Metas

**Metas:**
- Proporcionar un sistema robusto de identificación utilizando hashes MD5/SHA1/CRC32.
- Soportar el escaneo profundo dentro de archivos comprimidos (zip, rar, 7z) sin necesidad de extracción completa al disco.
- Implementar un "Orquestador" con estado que permita al usuario pausar, detener y reanudar tareas de scraping de larga duración.
- Proveer un motor de plantillas de nombres para soportar formatos personalizados (ej: `Nombre (Región) [Hash]`).
- Persistir la configuración y el progreso de la sesión entre reinicios de la aplicación.
- Garantizar que la interfaz de usuario siga siendo receptiva durante operaciones intensivas de E/S.

**No-Metas:**
- Jugar o emular juegos (el enfoque se limita a la identificación y el renombrado).
- Modificar los datos internos de la ROM (los metadatos son solo para renombrar el archivo).
- Descarga masiva de medios (imágenes/videos) en la versión inicial (el enfoque es el renombrado).

## Decisiones

- **Capa de Persistencia Cifrada: SQLite + AES-GCM**
  - *Justificación*: Necesitamos almacenar historial y el estado de miles de ROMs. Para las credenciales de APIs (ScreenScraper/TheGamesDB), se implementó cifrado `AES-256-GCM` en lado local (escritura/lectura) para una mayor protección. SQLite nos provee base de datos ligera, sin CGO integrándose perfectamente en Go.
  - *Alternativas consideradas*: Archivos JSON YAML en texto plano (rechazados por falta de indexación y seguridad).
- **Concurrencia: Pool de trabajadores para Hashing y llamadas a la API**
  - *Justificación*: El cálculo de hashes es intensivo en CPU e intensivo en E/S. Implementaremos un pool de trabajadores controlado por semáforos para limitar el acceso concurrente al sistema de archivos y evitar exceder los límites de frecuencia de las APIs.
  - *Alternativas consideradas*: Procesamiento lineal (demasiado lento para colecciones grandes).
- **Arquitectura del Puente: Eventos de Wails**
  - *Justificación*: Para el proceso de scraping de larga duración, utilizaremos los eventos de Wails para enviar actualizaciones en tiempo real (archivo actual, porcentaje, tiempo estimado) al frontend, en lugar de realizar consultas frecuentes desde JavaScript.
- **Manejo de Archivos Comprimidos: Librerías nativas de Go y especializadas**
  - *Justificación*: Uso de `archive/zip` para zip y paquetes especializados para rar/7z. Para 7z, utilizaremos decodificadores de transmisión (streaming) para calcular hashes sin extracción completa cuando sea posible.
- **Interfaz del Scraper: Patrón de Proveedor (Provider-Pattern)**
  - *Justificación*: Se definió una interfaz `Scraper` con métodos `SearchByHash` y `SearchByName`. El sistema ya incluye clientes oficiales y funcionales para **ScreenScraper** y **TheGamesDB**, ejecutables de forma individual o encadenada.
- **Motor de Renombrado Extendido: Plantillas Dinámicas**
  - *Justificación*: Se utiliza un sistema masivo de etiquetas (tokens): `{Name}, {Region}, {Company}, {Year}, {Hash}, {Languages}, {Developer}, {Genre}, {Players}, {Rating}, {RomType}` para personalización profunda.
  - *Implementación*: El parser de Go reemplaza estos tokens detectando fallbacks automáticos antes de actualizar la tabla Wails en la UI.

## Riesgos / Compensaciones

- **[Riesgo] Alta latencia de red en NAS** &rarr; *Mitigación*: Implementación de tiempos de espera (timeouts) ajustables y reintentos de conexión en el escáner del sistema de archivos.
- **[Riesgo] Límites de frecuencia de API y credenciales** &rarr; *Mitigación*: Limitador de frecuencia global en el orquestador. El usuario debe proporcionar sus propias credenciales para APIs privadas como ScreenScraper.
- **[Riesgo] Contención de archivos durante el renombrado** &rarr; *Mitigación*: Asegurar que el orquestador detenga toda la actividad de escaneo en un directorio antes de ejecutar la operación de renombrado masivo.

## Plan de Migración

- No aplicable para esta versión inicial (Proyecto Nuevo). Las versiones futuras incluirán migraciones de esquema para la base de datos SQLite.

## Decisiones Resueltas

- **Búsqueda Resiliente**: Se implementó una lógica dual donde las listas concurrentes siempre evalúan primero por Hash exacto y, si no hay resultados y las APIs lo permiten, realizan un fallback inteligente cruzando búsquedas por Nombre (filename).
- **Framework de UI**: Se escogió **Svelte** para el frontend y como controlador de Wails 2.0. Altere la UI un dashboard reactivo sin Virtual DOM que maneja `EventsOn("progress")` sin fugas de memoria con miles de actualizaciones DOM.
