## Context

El proyecto actual cuenta con un sistema básico de escaneo, scraping y renombramiento de ROMs. Sin embargo, carece de mecanismos robustos para lidiar con el lado oscuro de las colecciones (archivos comprimidos con múltiples juegos no relacionados, colisiones de nombres), presenta flaquezas en el feedback visual (falta de "mata-angustias" que no bloqueen los botones de interrumpir proceso), carece de contabilidad de la API e invisibiliza errores de renombramiento. Estos problemas degradan la experiencia de usuario y aumentan la posibilidad de errores y baneos de API.

## Goals / Non-Goals

**Goals:**
- Identificar y omitir paquetes comprimidos de recopilaciones de ROMs sin necesidad de extraerlos completamente.
- Proporcionar feedback visual constante (Mata-angustias) en el Frontend durante tareas largas sin bloquear interacciones críticas (Parar, Pausar).
- Manejar conflictos de nombres moviendo las ROMs duplicadas a un subdirectorio especial `duplicados` manteniendo la estructura de directorios relativa original.
- Mostrar un resumen final detallado de los errores de renombramiento.
- Contabilizar llamadas a la API y llevar un conteo mensual persistente en el backend para no exceder los límites de las APIs (como ScreenScraper u otras).

**Non-Goals:**
- Implementar extracción automática de archivos (solo lectura de índices/cabeceras).
- Solucionar los contenidos anómalos dentro del zip (si el archivo no cumple el criterio de ser "único título o relacionados", se ignora completamente el zip).
- Interfaz compleja para el resumen de errores; un panel claro y conciso es suficiente.

## Decisions

- **Filtrado de Archivos Comprimidos**: Utilizaremos librerías nativas o bien soportadas en Go (como `archive/zip` y adaptadores para `.rar` o `.7z`) para leer **exclusivamente el índice** del archivo. 
  - *Heurística a aplicar*: Si el archivo contiene múltiples archivos con extensiones válidas de ROM, se aplicará el heurístico de similitud de nombres. Si los bases names (sin extensión ni tokens tipo "Disc 1") son completamente diferentes o la cantidad de archivos excede un umbral lógico de DLCs/múltiples discos, el archivo se descartará.

- **Mata-Angustias (Loading Indicators)**: Se implementará en el Frontend mediante overlays dinámicos o *spinners*. El estado de "procesando" aislará el resto de la interfaz (inputs, settings) pero mantendrá explícitamente habilitados y con un `z-index` mayor a los botones de control de proceso (Pausar/Detener scraping o escaneo).

- **Gestión de Errores de Renombramiento**: La función que interacciona con los archivos (renombrado y movimiento) en el Backend no fallará silenciosamente ante un archivo bloqueado o faltante. Recopilará una estructura `[]RenameError{Path, Reason}` que se retornará al Frontend y mostrará un modal o tabla en la UI centralizando el resultado final del lote.

- **Duplicados y Recreación de Rutas**: Al detectar que un nombre destino ya existe en el directorio final, se cambiará el path destino por `<directorio_destino>/duplicados/.../nuevo_nombre.ext`.
  - Se requiere conocer el directorio raíz de la carpeta seleccionada en el escaneo para calcular y reconstruir la ruta relativa mediante funciones combinadas de `path/filepath`. `os.MkdirAll` se utilizará para crear las subcarpetas necesarias antes de efectuar el `os.Rename`.

- **Métrica de Uso de API**: Se agregará una estructura persistente en el fichero de configuración de la app o backend state (`config.json` o similar), alojando un par `{"MonthYear": "2026-03", "Count": X}` por proveedor.
  - En cada petición HTTP, en el middleware o capa del proveedor del cliente HTTP de Go, se incrementará el valor. Si el "MonthYear" difiere de `time.Now().Format("2006-01")`, el contador se reinicia a 1.

## Risks / Trade-offs

- **Risk**: Leer índices de archivos .7z o .rar pesados o solid-archives podría ser lento o traer problemas de compatibilidad en Go.
  - *Mitigation*: Emplear librerías específicas bien mantenidas o advertir al usuario/hacer time-out si el índice tarda mucho en obtenerse. Priorizar `archive/zip` inicialmente de ser imperativo.
- **Risk**: La recursión al aislar duplicados puede complicarse si involucra unidades separadas por permisos.
  - *Mitigation*: Las rutas dentro de la carpeta "duplicados" manejarán posibles conflictos resolviendo a la carpeta destino base de las otras ROMs elegida por el usuario.
- **Risk**: El contador de APIs puede perder desincronización si la app colapsa antes de escribirse a disco.
  - *Mitigation*: Guardar el conteo en batches regulares (p.ej. cada 10 peticiones) o utilizar un fichero temporal persistente rápido, o simplemente un `defer flushConfig()` en el apagado ordenado.
