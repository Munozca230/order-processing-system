# 🏗️ Arquitectura C4 – Worker de Procesamiento de Pedidos (Mermaid)

> Documento de arquitectura usando el modelo C4 representado en **Mermaid**.
> Incluye diagramas de Contexto (Nivel 1), Contenedores (Nivel 2) y Componentes (Nivel 3).

---

## 1. Diagrama de Contexto (Nivel 1)

```mermaid
%%{init: { 'theme': 'default' } }%%
flowchart TD
    user([Operador/Cliente])
    kafka[[Kafka]]
    worker["Order Worker (Java)"]
    productApi[("Product API (Go)")]
    clientApi[("Customer API (Go)")]
    mongo[(MongoDB)]
    redis[(Redis)]

    user -->|"Publica pedidos"| kafka
    kafka -->|"Consume"| worker
    worker -->|"Consulta"| productApi
    worker -->|"Consulta"| clientApi
    worker -->|"Persiste"| mongo
    worker -->|"Retries / locks"| redis
```

---

## 2. Diagrama de Contenedores (Nivel 2)

```mermaid
%%{init: { 'theme': 'default' } }%%
flowchart TD
    kafka[[Kafka]]

    subgraph "Sistema de Procesamiento de Pedidos"
        worker["Order Worker (Spring Boot WebFlux)"]
        productApi["Product API (Go)"]
        clientApi["Customer API (Go)"]
        mongo[(MongoDB)]
        redis[(Redis)]
    end

    kafka --> worker
    worker --> productApi
    worker --> clientApi
    worker --> mongo
    worker --> redis
```

---

## 3. Diagrama de Componentes (Nivel 3) – Worker

```mermaid
%%{init: { 'theme': 'default' } }%%
flowchart TD
    consumer[KafkaConsumer]
    lockMgr[LockManager]
    enrichment[EnrichmentService]
    validation[ValidationService]
    retryMgr[RetryManager]
    repo[OrderRepository]

    consumer --> lockMgr
    consumer --> enrichment
    enrichment --> validation
    validation --> repo
    enrichment --> retryMgr
```

---

### Flujo resumido
1. `KafkaConsumer` recibe el mensaje con `orderId`.
2. Se solicita un lock en `Redis` para evitar doble procesamiento.
3. `EnrichmentService` consulta las APIs Go (productos y clientes).
4. `ValidationService` verifica consistencia y reglas de negocio.
5. Si hay errores, `RetryManager` guarda el intento fallido en `Redis` y reprograma.
6. Pedido válido se almacena a través de `OrderRepository` en `MongoDB`.

---

### Notas adicionales
- **Reactivo:** Cadena completa con back-pressure (`Project Reactor`).
- **Resiliencia:** Retries exponenciales y locks distribuidos en Redis.
- **Escalabilidad:** Varias instancias del worker en el mismo consumer group.
- **Observabilidad:** Logs y métricas centralizadas (no dibujadas para simplificar).
