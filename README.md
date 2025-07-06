# 🚀 **Order Processing System**

**Sistema enterprise de procesamiento de órdenes en tiempo real** que consume mensajes de Kafka, enriquece datos consultando APIs externas, valida reglas de negocio y almacena resultados en MongoDB. Implementa patrones de resiliencia con reintentos exponenciales, distributed locking y manejo de errores.

**Arquitectura**: Order Worker (Java 21 + WebFlux) consume de Kafka → Enriquece datos via Product/Customer APIs (Go) → Valida cliente activo → Persiste en MongoDB con estructura específica.

**Compatibilidad**: Windows (PowerShell), macOS (Terminal), Linux (Bash). Scripts PowerShell incluidos para máxima compatibilidad Windows.

![Architecture](https://img.shields.io/badge/Architecture-Clean%20Architecture-blue) ![Java](https://img.shields.io/badge/Java-21-orange) ![Go](https://img.shields.io/badge/Go-1.22-blue) ![Spring](https://img.shields.io/badge/Spring-WebFlux-green) ![MongoDB](https://img.shields.io/badge/MongoDB-7.0-green) ![Kafka](https://img.shields.io/badge/Kafka-3.6-red)

---

## ⚡ **Quick Start por Sistema Operativo**

### **🖥️ Windows (PowerShell)**
```powershell
# Opción 1: Scripts PowerShell (sin instalaciones)
.\scripts\deploy-frontend.ps1
.\scripts\test-final-system.ps1

# Opción 2: Git Bash para usar make
# Abrir Git Bash, luego:
make frontend
make test
```

### **🍎 macOS**
```bash
# 1. Instalar make (una vez)
brew install make

# 2. Deploy y validación
make frontend && make test
```

### **🐧 Linux**
```bash
# 1. Instalar make (una vez)
sudo apt install make  # Ubuntu/Debian
# sudo yum install make  # CentOS/RHEL

# 2. Deploy y validación
make frontend && make test
```

**Resultado**: Sistema funcionando en http://localhost:8080 + testing ejecutado

### **🚀 Comandos Principales**

```bash
make frontend    # Deploy completo con web UI
make backend     # Solo APIs + Worker  
make test        # Testing nativo (Go + Java)
make status      # Ver estado de servicios
make clean       # Restart con datos frescos
make help        # Ver todos los comandos
```

---

## 🎯 **Para Revisores - Validación Rápida**

### **📋 Opción 1: Scripts PowerShell (Windows - Sin instalaciones)**
```powershell
# Deploy completo (2-3 minutos)
.\scripts\deploy-frontend.ps1

# Testing automático (2-3 minutos) 
.\scripts\test-final-system.ps1

# Validar en navegador
start http://localhost:8080
```

### **📋 Opción 2: Make (Linux/macOS/Git Bash)**
```bash
# Instalar make (una vez)
brew install make           # macOS
sudo apt install make       # Linux

# Deploy y testing
make frontend && make test

# Validar
open http://localhost:8080   # macOS
# xdg-open http://localhost:8080  # Linux
```

### **📊 URLs del Sistema**
- **Frontend Web**: http://localhost:8080 (Interfaz completa)
- **Product API**: http://localhost:8081/health
- **Customer API**: http://localhost:8082/health  

### **🔧 Mínimo (Solo Docker)**
```bash
cd infra
docker-compose --profile frontend up -d
# Esperar ~30 segundos, luego abrir http://localhost:8080
```

---

## 📋 **Arquitectura del Sistema**

```
🌐 Frontend → 📨 Order API → 📨 Kafka → ⚙️ Order Worker (Java 21)
                                            ↓
🛍️ Product API ← 🔍 Enrichment ← 👥 Customer API ← ✅ Validation  
     ↓                                ↓
💾 MongoDB (Catalog) ← 📊 Storage ← 💾 MongoDB (Orders)
```

### **🔧 Componentes**
- **Order Worker** (Java 21 + WebFlux): Consume Kafka, enriquece datos, valida y persiste
- **Product/Customer APIs** (Go + Clean Architecture): Proveen datos del catálogo  
- **Frontend** (HTML/JS): Interfaz web que consume las APIs directamente
- **Infrastructure**: Kafka, MongoDB, Redis para locking y retries

### **📊 Testing Options**

**Visual**: http://localhost:8080 (crear órdenes en interfaz web)  
**Postman**: Importar collection desde `/postman/`  
**CLI**: `make test` (unit + integration + e2e automático)

---

## 🛠️ **Troubleshooting**

### **Error "&&" no válido (Windows PowerShell)**
```powershell
# En lugar de: make frontend && make test
# Usar comandos separados:
make frontend
make test

# O usar scripts PowerShell:
.\scripts\deploy-frontend.ps1
.\scripts\test-final-system.ps1
```

### **make: command not found**
```bash
# Windows: Usar Git Bash o PowerShell con scripts
# macOS: brew install make
# Linux: sudo apt install make
```

### **Error de red Docker: "network not found"**
```bash
# Limpiar redes Docker huérfanas
make docker-clean

# Luego intentar de nuevo
make frontend
```

### **Servicios no responden**
```powershell
# Windows PowerShell
.\scripts\dev-reset.ps1

# Git Bash/Linux/macOS
make clean && make status
```

### **Tests fallan**
```bash
# Verificar dependencias
java -version     # Debe ser Java 21
docker --version  # Docker Desktop debe estar corriendo
```

---

## 📚 **Documentación Técnica**

- **[📋 Arquitectura Completa](docs/COMPLETE_ARCHITECTURE_DIAGRAMS.md)** - Diagramas detallados, secuencias, componentes
- **[📄 Especificación Original](prueba.md)** - Requerimientos técnicos cumplidos
- **[🔧 Cross-Platform Setup](docs/CROSS_PLATFORM_SETUP.md)** - Instalación avanzada

---

## ✅ **Requerimientos Cumplidos**

✅ **Worker Java 21** con Spring Boot WebFlux reactivo  
✅ **Consumo Kafka** con consumer groups + rebalancing automático  
✅ **APIs Go** con Clean Architecture + MongoDB persistence  
✅ **Enriquecimiento de datos** via WebClient reactivo  
✅ **Validación** business rules + customer active  
✅ **MongoDB storage** con estructura exacta según especificación  
✅ **Reintentos exponenciales** con backoff + dead letter queue  
✅ **Distributed locking** con Redis locks + TTL automático  
✅ **Testing completo** Testcontainers + Integration + E2E + Postman

**🚀 Sistema enterprise-ready con 100% cumplimiento de requerimientos + frontend interactivo.**