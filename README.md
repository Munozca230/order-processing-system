# Order Processing System

Sistema de procesamiento de pedidos con arquitectura de microservicios usando Java, Go, Kafka, MongoDB y Redis.

## 📋 Requerimientos completados
- ✅ Worker Java con Spring Boot WebFlux
  - ✅ Consumo de mensajes Kafka
  - ✅ Enriquecimiento con APIs Go externas
  - ✅ Validación de datos
  - ✅ Persistencia en MongoDB
  - ✅ **Reintentos exponenciales con Redis**
  - ✅ Distributed locking
  - ✅ Testing con Testcontainers

## 📁 Estructura del proyecto
  📦 order-processing-system/
  ├── 📂 docs/                    # Documentación
  │   ├── ARCHITECTURE.md         # Guía de arquitectura
  │   ├── CLAUDE.md              # Configuración para Claude
  │   └── architecture-diagram.md # Diagramas mermaid
  ├── 📂 scripts/                 # Scripts de testing
  │   ├── test-retry.ps1         # Script PowerShell
  │   └── test-retry.sh          # Script Bash
  ├── 📂 services/                # Microservicios
  │   ├── 📂 order-worker/       # Worker Java (Spring Boot)
  │   ├── 📂 product-api/        # API Go productos
  │   └── 📂 customer-api/       # API Go clientes
  ├── 📂 infra/                   # Infraestructura
  │   └── docker-compose.yml     # Orquestación Docker
  └── prueba.md                   # Especificaciones originales

## 🚀 Quick Start

### Desarrollo (Testing)
```bash
  cd services/order-worker
  mvn test

  Producción (Docker)

  cd infra
  docker-compose up -d

  Probar reintentos exponenciales

  cd scripts
  ./test-retry.ps1  # Windows
  ./test-retry.sh   # Linux/Mac
```

### 🏗️ Arquitectura

Sistema reactivo con reintentos exponenciales:
  - Kafka: Mensajería asíncrona
  - Redis: Cola de reintentos y distributed locking
  - MongoDB: Persistencia de pedidos procesados
  - Go APIs: Servicios de productos y clientes
  - Java Worker: Procesamiento reactivo con WebFlux

## 📚 Documentación adicional

- docs/ARCHITECTURE.md
- docs/CLAUDE.md
- docs/architecture-diagram.md