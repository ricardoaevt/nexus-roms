## Por qué

Actualmente, el renombrador de ROMs tiene algunas deficiencias en el manejo de archivos comprimidos complejos, transparencia en los errores, manejo de duplicados, feedback visual durante procesos largos (lo que genera incertidumbre o "angustia" en el usuario), y falta de visibilidad en el consumo de las APIs utilizadas. Estos cambios son necesarios para mejorar drásticamente la experiencia de usuario, evitar pérdida de datos o sobreescrituras accidentales, y mantener un control del límite de llamadas a los proveedores de metadatos.

## Qué Cambia

- **Filtrado de archivos comprimidos**: Se descartarán los archivos .zip, .rar y .7z que en su interior contengan múltiples ROMs o archivos que no estén relacionados (como los típicos "romsets" o grandes recopilaciones). Solamente se mantendrán aquellos archivos comprimidos que, aun conteniendo varios archivos, sean correspondientes a DLCs, parches o expansiones de una misma ROM válida.
- **Indicadores de carga (Mata-angustias)**: Se añadirán indicadores de estado de carga en los procesos prolongados. Los controles para pausar o detener el scraping seguirán funcionando y no serán bloqueados por estos indicadores.
- **Reporte de errores de renombramiento**: Se mejorará el sistema de registro y visualización para mostrar claramente por qué falló un renombramiento y proporcionar una lista exacta de qué ROMs no pudieron ser procesadas.
- **Manejo inteligente de duplicados**: Si una ROM renombrada entra en conflicto de nombre con otra existente, se moverá a una carpeta especial llamada `duplicados`. Dentro de esta, se recreará la estructura del árbol de directorios original para saber exactamente de dónde provino la ROM repetida.
- **Seguimiento de cuotas de API**: Se implementará un contador mensual de llamadas por cada proveedor de API configurado, reiniciándose automáticamente cada mes.

## Capacidades

### Nuevas Capacidades
- `archive-filtering`: Lógica para inspeccionar el contenido de archivos .zip, .rar y .7z, con el objetivo de descartar aquellos que funcionen como recopilaciones o directorios de ROMs múltiples no relacionadas, reteniendo aquellos que contengan únicamente archivos vinculados a un mismo título (como DLCs o actualizaciones).
- `loading-indicators`: Sistema de feedback visual no bloqueante para operaciones de larga duración.
- `rename-error-reporting`: Recopilación detallada y presentación de errores durante el proceso de renombramiento en lote.
- `duplicate-handling`: Lógica para detección de colisiones de nombres y reubicación de archivos preservando su ruta relativa en un directorio de contención.
- `api-rate-tracking`: Almacenamiento persistente y visualización de estadísticas de uso de APIs con reinicio mensual.

### Capacidades Modificadas


## Impacto

- **Código afectado**: Lógica de escaneo de archivos (para el interior de zips/rars), lógica principal de renombramiento y movimiento de archivos, controladores de frontend para la UI, y la capa de persistencia de configuraciones para llevar la cuenta de las APIs.
- **Sistemas**: Impacto mínimo en dependencias, posiblemente se requiera mejorar la integración con librerías de descompresión o de lectura de archivos comprimidos para inspeccionar el índice sin extraer todo.
