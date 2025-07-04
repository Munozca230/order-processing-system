# 📚 Guía de Arquitectura General – Proyecto de Procesamiento de Pedidos

Esta guía resume los principios, componentes y prácticas recomendadas para desarrollar y evolucio­­nar el proyecto completo. Sirve como referencia para nuevos integrantes y como checklist durante la implementación.

---

## 1. Objetivos Arquitectónicos
1. **Fiabilidad & Resiliencia**: Garantizar procesamiento exacto‐una‐vez incluso ante fallos de red o caídas parciales.
2. **Escalabilidad Horizontal**: Permitir añadir instancias del worker y de las APIs sin downtime.
3. **Baja Latencia**: Minimizar el tiempo desde la publicación del pedido hasta su persistencia enriquecida.
4. **Observabilidad**: Métricas, trazas y logs estructurados centralizados.
5. **Simplicidad Operacional**: Despliegues reproducibles vía contenedores y pipelines CI/CD.

---

## 2. Vista de Contexto
El sistema interactúa con:
• Operador/Cliente → envía pedidos a Kafka.
• Sistemas externos (Product API & Customer API) para enriquecimiento.
• Bases de datos (MongoDB y Redis) para persistencia y control de flujo.

---

## 3. Componentes Clave
| Categoría | Componente | Descripción Breve |
|-----------|-----------|-------------------|
| Aplicación | **Order Worker** (Java 21, Spring Boot + WebFlux) | Consume pedidos, llama APIs, valida y guarda resultados |
| Servicio | **Product API** (Go) | Obtiene información de catálogo |
| Servicio | **Customer API** (Go) | Obtiene detalles del cliente |
| Mensajería | **Kafka** | Bufferiza pedidos y desacopla emisores del worker |
| Base datos | **MongoDB** | Almacena documentos de pedido enriquecidos |
| Cache/Control | **Redis** | Retries (exponenciales) y lock distribuido |
| Observabilidad | Stack ELK / Grafana | Centralización de logs y métricas |

---

## 4. Flujo de Procesamiento
1. **Publicación**: El operador publica un mensaje con `orderId`, `customerId`, lista de productos.
2. **Consumo**: El worker (consumer group) extrae el mensaje.
3. **Lock**: Solicita un lock en Redis para evitar duplicidad.
4. **Enriquecimiento**: Obtiene datos de productos y cliente (APIs Go).
5. **Validación**: Verifica existencia y estado.
6. **Persistencia**: Guarda documento en MongoDB.
7. **Confirmación**: Confirma al offset de Kafka.
8. **Error Path**: Si falla un paso, registra intento en Redis y reintenta con backoff.

---

## 5. Aspectos No Funcionales
### 5.1 Resiliencia
• Retries exponenciales con número máximo configurable.
• Circuit-breakers ante fallos repetidos de APIs externas.
• Almacenamiento de mensajes fallidos en Redis para posterior análisis.

### 5.2 Escalabilidad
• Particiones de Kafka ≥ número previsto de instancias del worker.
• Sharding automático de datos en MongoDB si el volumen crece.

### 5.3 Seguridad
• Autenticación entre servicios con mTLS o JWT interno.
• Variables sensibles gestionadas en secretos de orquestador.

### 5.4 Observabilidad
• Logs JSON + trazas distribuidas (OpenTelemetry).
• Métricas de negocio (pedidos procesados, latencia) y técnicas (GC, heap).

---

## 6. Infraestructura & Despliegue
1. **Entorno Local**: Docker Compose con servicios: kafka, zookeeper, mongo, redis, worker, apis.
2. **Persistencia**: Volúmenes nombrados o bind mounts para bases de datos.
3. **Ambientes**: dev → staging → prod.
4. **CI/CD**: Pipeline que construye imágenes, ejecuta pruebas y despliega a Kubernetes o ECS.

---

## 7. Estructura de Carpetas Recomendada
```
/                     # raíz del repo mono o multi-repo
  ├─ order-worker/    # proyecto Java
  │   ├─ src/
  │   └─ Dockerfile
  ├─ product-api/     # servicio Go
  ├─ customer-api/    # servicio Go
  ├─ infra/
  │   ├─ docker-compose.yml
  │   └─ k8s/         # manifiestos opcionales
  ├─ docs/            # documentación y diagramas Mermaid/C4
  └─ ci/              # scripts y plantillas de pipeline
```
*(La estructura puede dividirse en repos separados si el equipo lo prefiere.)*

---

## 8. Roadmap de Desarrollo
1. **MVP End-to-End** (flujo feliz sin retries).
2. Validaciones de datos y manejo de errores.
3. Retries & locks distribuidos.
4. Observabilidad (logs, métricas, trazas).
5. CI/CD y pruebas de carga.
6. Endurecimiento de seguridad.

---

## 9. Riesgos y Mitigaciones
| Riesgo | Impacto | Mitigación |
|--------|---------|------------|
| Sobrecarga de APIs externas | Pedidos en cola, alta latencia | Circuit-breaker + cache | 
| Duplica­ción de procesamiento | Datos inconsistentes | Lock con expiración y verificación de idempotencia | 
| Falta de visibilidad | Difícil depurar | Observabilidad temprana | 
| Crecimiento explosivo de datos | Rendimiento degradado | Índices, TTL y sharding |

---

## 10. Glosario
• **Idempotencia**: Capacidad de procesar el mismo mensaje varias veces sin efectos adversos.
• **Back-pressure**: Regulación del flujo para no sobrecargar consumidores.

---

Esta guía complementa los diagramas Mermaid y debe mantenerse actualizada conforme el sistema evolucione.
