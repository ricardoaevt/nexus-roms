## Why

La gestión de colecciones de ROMs de juegos retro es una tarea tediosa debido a que los archivos suelen tener nombres inconsistentes o técnicos. Este software automatiza el proceso de identificación y renombrado masivo conectándose a APIs de metadatos, permitiendo a los usuarios tener colecciones organizadas y legibles con un esfuerzo mínimo.

## What Changes

Se desarrollará una aplicación de escritorio completa utilizando **Go + Wails** con las siguientes capacidades:
- **Selección de Directorio**: Soporte para rutas locales y de red.
- **Soporte de Archivos**: Procesamiento de ROMs tanto en formato directo como comprimidos (**zip, rar, 7z**).
- **Identificación por Hash**: Cálculo automático de hashes para una identificación precisa mediante APIs.
- **Scraping Multi-API**: Integración configurable con múltiples proveedores de metadatos simultáneos (ScreenScraper, TheGamesDB).
- **Control de Flujo**: Funcionalidad para pausar, reanudar o detener el proceso de scraping.
- **Persistencia Inteligente**: Modal interactivo para retomar sesiones de scraping incompletas, calcular avances o iniciar desde cero.
- **Previsualización de Renombrado**: Interfaz interactiva para revisar y filtrar cambios antes de aplicarlos.
- **Formato Personalizable**: Los usuarios podrán definir el patrón de renombrado extendido (ej: `{Name} ({Region}) [{Hash}]`, además de variables como `{Languages}, {Genre}, {Developer}, {Rating}, {RomType}`).
- **Configuración Avanzada Segmentada**: Gestión de credenciales (cifradas en disco local con núcleo AES-256-GCM), activación independiente de APIs y estrategias de búsqueda (Hash vs Nombre).
- **Sistema de Logs**: Ventana dedicada para el rastreo de errores de conexión y fallos en el sistema.

## Capabilities

### New Capabilities
- `api-management`: Gestión centralizada de múltiples APIs (ScreenScraper, TheGamesDB), incluyendo autenticación cifrada por AES, URLs base y activación de estrategias de búsqueda mixtas por hash/nombre.
- `scraping-orchestration`: Lógica para el control del ciclo de vida del proceso en background, cálculo stream de hashes nativos y recuperación heurística de sesiones rotas.
- `filesystem-operations`: Escaneo de directorios, soporte deep-scan para **zip, rar y 7z**, motor avanzado de plantillas de tokens de metadatos, y ejecución paralela de renombrado masivo.
- `user-interface-wails`: Dashboard Svelte que integra la vista de tabla, modales de control de estados, configuración en cascada de APIs y consola de logs en vivo.

### Modified Capabilities
- Ninguna (Proyecto nuevo).

## Impact

- **Tecnología**: Uso de Go para el backend y Wails para la interfaz de usuario.
- **Dependencias**: Conexión con APIs externas de metadatos (requiere gestión de secretos y configuración).
- **Almacenamiento**: Creación de un sistema de persistencia local para configuraciones y estado de sesiones.
- **Red**: Capacidad de operar sobre sistemas de archivos en red (SMB/NFS).
