# Setup Simple

## Para nuevos desarrolladores

1. **Clonar repo**
2. **Ejecutar**: `scripts/deploy-frontend.ps1`
3. **Listo** - abrir http://localhost:8080

## Para desarrolladores existentes

**Después de git pull**:
```powershell
scripts/fresh-start.ps1
```

Esto garantiza que todos tengan exactamente los mismos datos iniciales.

## Si algo se rompe

```powershell
scripts/dev-reset.ps1
```

## ¿Qué hace fresh-start?

1. Para servicios
2. Elimina datos viejos de MongoDB
3. Reinicia con datos frescos idénticos
4. Todos los desarrolladores tienen los mismos datos

## Datos iniciales siempre iguales

- **8 customers** fijos (mismos IDs, nombres, direcciones)
- **9 products** fijos (mismos IDs, precios, nombres)
- **Sin configuración** - funciona automáticamente
- **Sin problemas** de compatibilidad entre desarrolladores