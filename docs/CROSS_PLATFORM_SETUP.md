# 🌍 Cross-Platform Setup Guide

Este proyecto funciona en **Linux**, **macOS** y **Windows** usando diferentes métodos.

## 🎯 **Métodos Disponibles (Ordenados por Recomendación)**

### **1. Make (Universal) - ⭐ RECOMENDADO**

Funciona en todos los sistemas operativos:

```bash
# Deploy completo
make frontend

# Solo backend
make backend

# Testing
make test

# Ver estado
make status
```

**Instalación de Make:**
- **macOS**: `brew install make` (o usar Xcode tools)
- **Linux**: `sudo apt install make` (Ubuntu) o `sudo yum install make` (RHEL)
- **Windows**: 
  - Git Bash (incluido con Git)
  - Chocolatey: `choco install make`
  - WSL: `sudo apt install make`

### **2. PowerShell Core - 🚀 MANTIENE FUNCIONALIDAD COMPLETA**

PowerShell ahora es cross-platform:

```bash
# Instalar PowerShell Core
# macOS
brew install powershell

# Linux (Ubuntu)
wget -q https://packages.microsoft.com/config/ubuntu/20.04/packages-microsoft-prod.deb
sudo dpkg -i packages-microsoft-prod.deb
sudo apt-get update
sudo apt-get install -y powershell

# Windows: Ya incluido
```

**Uso:**
```bash
pwsh ./scripts/deploy-frontend.ps1
pwsh ./scripts/deploy-backend.ps1
```

### **3. Scripts Bash - 🐧 LINUX/MACOS NATIVO**

Scripts bash universales:

```bash
# Hacer ejecutables (una vez)
chmod +x scripts/*.sh

# Usar
./scripts/deploy-frontend.sh
./scripts/deploy-backend.sh
```

### **4. Docker Compose Directo - 🐳 MÍNIMO DEPENDENCIAS**

Solo requiere Docker:

```bash
cd infra

# Frontend completo
docker-compose --profile frontend up -d

# Solo backend
docker-compose up -d

# Ver estado
docker-compose ps

# Logs
docker-compose logs -f order-worker
```

## 📋 **Comparación de Métodos**

| Método | Pros | Contras | Mejor Para |
|--------|------|---------|------------|
| **Make** | Universal, simple, profesional | Requiere instalación | **Producción, CI/CD** |
| **PowerShell Core** | Mantiene toda la funcionalidad | Requiere instalación en macOS/Linux | **Desarrollo avanzado** |
| **Bash** | Nativo en Unix, sin dependencias | No funciona en Windows cmd | **Desarrolladores Unix** |
| **Docker Compose** | Mínimas dependencias | Menos validaciones | **Setup rápido** |

## 🛠️ **Instalación por Sistema Operativo**

### **macOS (Desarrolladores revisores)**

```bash
# Opción 1: Homebrew (recomendado)
brew install make powershell

# Verificar
make --version
pwsh --version

# Usar cualquier método
make frontend                    # Método 1
pwsh ./scripts/deploy-frontend.ps1  # Método 2
./scripts/deploy-frontend.sh        # Método 3
```

### **Linux**

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install make

# CentOS/RHEL
sudo yum install make

# PowerShell (opcional)
wget -q https://packages.microsoft.com/config/ubuntu/20.04/packages-microsoft-prod.deb
sudo dpkg -i packages-microsoft-prod.deb
sudo apt-get update
sudo apt-get install -y powershell
```

### **Windows**

```powershell
# PowerShell (nativo)
.\scripts\deploy-frontend.ps1

# Git Bash
make frontend

# WSL
sudo apt install make
make frontend
```

## 🎯 **Recomendación para Revisores**

**Para revisores en macOS:**

1. **Opción más simple:**
   ```bash
   cd infra
   docker-compose --profile frontend up -d
   ```

2. **Opción con validaciones completas:**
   ```bash
   brew install make
   make frontend
   ```

3. **Opción manteniendo PowerShell:**
   ```bash
   brew install powershell
   pwsh ./scripts/deploy-frontend.ps1
   ```

## 📖 **Funcionalidades por Método**

| Funcionalidad | Make | PowerShell | Bash | Docker Direct |
|---------------|------|------------|------|---------------|
| Health checks dinámicos | ✅ | ✅ | ✅ | ❌ |
| Validaciones Docker | ✅ | ✅ | ✅ | ❌ |
| Output coloreado | ✅ | ✅ | ✅ | ❌ |
| Progress feedback | ✅ | ✅ | ✅ | ❌ |
| Error handling | ✅ | ✅ | ✅ | ❌ |
| Cross-platform | ✅ | ✅ | ⚠️ | ✅ |

**Leyenda**: ✅ Completo | ⚠️ Limitado | ❌ No disponible

## 🚨 **Troubleshooting Cross-Platform**

### **Error: "make: command not found"**
```bash
# macOS
brew install make

# Ubuntu/Debian  
sudo apt install make

# Windows
# Use Git Bash or install via Chocolatey
```

### **Error: "pwsh: command not found"**
```bash
# macOS
brew install powershell

# Linux
# Follow PowerShell installation guide
```

### **Permisos en scripts bash**
```bash
chmod +x scripts/*.sh
```

### **Windows: Scripts no ejecutan**
```powershell
# Cambiar política de ejecución
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```