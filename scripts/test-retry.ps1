# Script para probar el sistema de reintentos exponenciales
Write-Host "=== Prueba del Sistema de Reintentos Exponenciales ===" -ForegroundColor Green

# Cambiar al directorio de infra
Set-Location "../infra"

# 1. Verificar que todos los servicios están corriendo
Write-Host "1. Verificando servicios..." -ForegroundColor Yellow
docker-compose ps

# 2. Enviar mensaje de prueba a Kafka
Write-Host "2. Enviando mensaje de prueba a Kafka..." -ForegroundColor Yellow
$testMessage = '{"orderId":"test-001","customerId":"customer-1","products":["product-1","product-2"]}'
echo $testMessage | docker-compose exec -T kafka kafka-console-producer.sh --bootstrap-server localhost:9092 --topic orders

# 3. Verificar logs del order-worker
Write-Host "3. Verificando logs del order-worker..." -ForegroundColor Yellow
docker-compose logs --tail=50 order-worker

# 4. Verificar datos en Redis
Write-Host "4. Verificando datos en Redis..." -ForegroundColor Yellow
docker-compose exec redis redis-cli KEYS "*"

# 5. Simular error parando las APIs
Write-Host "5. Simulando errores - parando product-api..." -ForegroundColor Yellow
docker-compose stop product-api

# 6. Enviar otro mensaje para triggear retry
Write-Host "6. Enviando mensaje para triggear retry..." -ForegroundColor Yellow
$retryMessage = '{"orderId":"test-002","customerId":"customer-2","products":["product-3"]}'
echo $retryMessage | docker-compose exec -T kafka kafka-console-producer.sh --bootstrap-server localhost:9092 --topic orders

# 7. Verificar logs de retry
Write-Host "7. Verificando logs de retry..." -ForegroundColor Yellow
Start-Sleep -Seconds 5
docker-compose logs --tail=30 order-worker

# 8. Verificar cola de retry en Redis
Write-Host "8. Verificando cola de retry en Redis..." -ForegroundColor Yellow
docker-compose exec redis redis-cli ZRANGE retry_queue 0 -1 WITHSCORES

# 9. Restaurar product-api
Write-Host "9. Restaurando product-api..." -ForegroundColor Yellow
docker-compose start product-api

Write-Host "=== Prueba completada ===" -ForegroundColor Green
Write-Host "Revisa los logs para ver los reintentos exponenciales en acción" -ForegroundColor Cyan