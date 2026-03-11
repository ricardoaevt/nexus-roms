## 1. Infraestructura y Persistencia

- [x] 1.1 Configurar el esquema de la base de datos SQLite (tablas para configuración, sesiones y archivos).
- [x] 1.2 Implementar la capa de acceso a datos con cifrado real (AES-256-GCM) para persistir credenciales de API.
- [x] 1.3 Crear utilidades para la gestión de la sesión de scraping (guardar/recuperar estado).

## 2. Escaneo y Hashing de Archivos

- [x] 2.1 Implementar el escáner de directorios con soporte para rutas locales y de red.
- [x] 2.2 Desarrollar el motor de hashing para archivos directos (MD5, SHA1, CRC32).
- [x] 2.3 Implementar el hashing por flujo (streaming) para archivos dentro de contenedores ZIP, RAR y 7Z.

## 3. Integración de APIs de Metadatos

- [x] 3.1 Crear la interfaz genérica `Scraper` y el cliente para ScreenScraper.
- [x] 3.2 Implementar la lógica de búsqueda por Hash con soporte para múltiples algoritmos.
- [x] 3.3 Implementar la búsqueda por Nombre como respaldo y el sistema de reintentos con limitación de frecuencia.
- [x] 3.4 **(Extra)** Desarrollar el cliente Scraper completo para TheGamesDB como proveedor paralelo (Hash+Name).

## 4. Orquestador de Scraping

- [x] 4.1 Desarrollar el orquestador central con soporte para pausa, reanudación y parada.
- [x] 4.2 Implementar el pool de trabajadores concurrentes para el procesamiento de archivos.
- [x] 4.3 Añadir el sistema de eventos para notificar el progreso en tiempo real al frontend.
- [x] 4.4 **(Extra)** Lógica profunda de reanudación y conteo preciso sobre sesiones pausadas/detenidas.

## 5. Sistema de Renombrado y Plantillas

- [x] 5.1 Construir el motor de plantillas para soportar tokens como `{Name}`, `{Region}` y `{Hash}`.
- [x] 5.2 Implementar la lógica de generación de previsualización de cambios.
- [x] 5.3 Desarrollar la ejecución segura del renombrado masivo con verificación de errores.

## 6. Interfaz de Usuario (Wails + Svelte)

- [x] 6.1 Crear el dashboard principal con indicadores de progreso y logs en tiempo real.
- [x] 6.2 Implementar la gestión de configuración de APIs y credenciales.
- [x] 6.3 Añadir la tabla de previsualización para revisar cambios antes de aplicar.
- [x] 6.4 Diseñar una interfaz premium con modo oscuro y estética profesional.
- [x] 6.5 **(Extra)** Modal Interactivo (UI) para manejar amigablemente la recuperación / reinicio de sesiones previas en directorios.
