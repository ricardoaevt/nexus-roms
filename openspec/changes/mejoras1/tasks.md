## 1. Filtrado de Archivos Comprimidos (`archive-filtering`)

- [x] 1.1 Implementar adaptador en Go para extraer índices de `.zip`, `.rar` y `.7z` sin descompresión total
- [x] 1.2 Crear la heurística de evaluación temporal para nombres paralelos frente a recopilaciones desconectadas
- [x] 1.3 Ajustar la fase inicial de escaneo para detectar e ignorar archivos bloqueados u omitidos por la heurística

## 2. Indicadores de Estado y UI Reactiva (`loading-indicators`)

- [x] 2.1 Diseñar componente overlay semitransparente de "proceso activo" en el frontend
- [x] 2.2 Desacoplar CSS de botones críticos de detención ("Pausar", "Detener") dotándoles de z-index elevado frente al overlay
- [x] 2.3 Enganchar señales de "inicio scraping" y "final scraping" a la visibilidad del overlay

## 3. Seguimiento Mensual de API (`api-rate-tracking`)

- [x] 3.1 Agregar sección temporal en los ajustes para hospedar el recordatorio persistente `{"Period": "YYYY-MM", "Count": 0}` para métricas de ScreenScraper
- [x] 3.2 Crear función middleware interceptora en Go (o wrap client HTTP) para incrementar el contador antes del *fetch*
- [x] 3.3 Implementar reinicio automático al detectar el cambio en el formato mes/año y enviar data actualizada a frontend

## 4. Renombramiento Inteligente y Manejo de Errores (`duplicate-handling`, `rename-error-reporting`)

- [x] 4.1 Refactorizar la función actual de `os.Rename` para atrapar *os errors* en lugar de paniquear o morir silenciosamente
- [x] 4.2 Desarrollar el cálculo de rutas y generación recursiva (`MkdirP`) bajo el nuevo directorio `[root]/duplicados/...` ante una colisión de nombre
- [x] 4.3 Empaquetar y retornar un listado explícito tipo `[]ErrorResult` tras concluir todo el ciclo batch del backend
- [x] 4.4 Mostrar mediante un modal al final el reporte claro de todo archivo saltado por bloqueo OS, o fallos en movimiento
- [x] 4.5 Omitir la carpeta `duplicados` de forma recursiva durante la fase de escaneo inicial para evitar re-procesar archivos aislados.
