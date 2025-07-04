# Script para probar el flujo end-to-end completo
Write-Host "=== Test End-to-End: Order Processing System ===" -ForegroundColor Green

# Cambiar al directorio de infra
Set-Location "../infra"

# 1. Verificar servicios
Write-Host "1. Verificando que todos los servicios estén corriendo..." -ForegroundColor Yellow
docker-compose ps

# 2. Test del flujo exitoso
Write-Host "`n2. Test Caso Exitoso - Cliente activo + productos existentes" -ForegroundColor Yellow
$successMessage = '{"orderId":"e2e-success-001","customerId":"customer-1","products":["product-1","product-2"]}'
echo $successMessage | docker-compose exec -T kafka kafka-console-producer.sh --bootstrap-server localhost:9092 --topic orders

Write-Host "   Esperando procesamiento..." -ForegroundColor Gray
Start-Sleep -Seconds 5

Write-Host "   Logs del procesamiento:" -ForegroundColor Gray
docker-compose logs --tail=20 order-worker

# 3. Verificar datos en MongoDB
Write-Host "`n3. Verificando datos guardados en MongoDB..." -ForegroundColor Yellow
docker-compose exec mongo mongosh --eval "use orders; db.orders.find({orderId: 'e2e-success-001'}).pretty();"

# 4. Test con cliente inactivo
Write-Host "`n4. Test Caso Error - Cliente inactivo" -ForegroundColor Yellow
$inactiveCustomerMessage = '{"orderId":"e2e-inactive-001","customerId":"customer-inactive","products":["product-1"]}'
echo $inactiveCustomerMessage | docker-compose exec -T kafka kafka-console-producer.sh --bootstrap-server localhost:9092 --topic orders

Write-Host "   Esperando procesamiento..." -ForegroundColor Gray
Start-Sleep -Seconds 3

# 5. Test con producto inexistente
Write-Host "`n5. Test Caso Error - Producto inexistente" -ForegroundColor Yellow
$missingProductMessage = '{"orderId":"e2e-missing-001","customerId":"customer-1","products":["product-999"]}'
echo $missingProductMessage | docker-compose exec -T kafka kafka-console-producer.sh --bootstrap-server localhost:9092 --topic orders

Write-Host "   Esperando procesamiento..." -ForegroundColor Gray
Start-Sleep -Seconds 3

# 6. Ver todos los logs recientes
Write-Host "`n6. Logs de procesamiento de todos los casos:" -ForegroundColor Yellow
docker-compose logs --tail=40 order-worker

# 7. Verificar todos los datos en MongoDB
Write-Host "`n7. Resumen de datos en MongoDB:" -ForegroundColor Yellow
docker-compose exec mongo mongosh --eval "use orders; db.orders.find().pretty(); print('Total orders: ' + db.orders.countDocuments());"

# 8. Verificar datos en Redis
Write-Host "`n8. Estado de Redis (locks y retry queues):" -ForegroundColor Yellow
docker-compose exec redis redis-cli KEYS "*"

Write-Host "`n=== Test End-to-End Completado ===" -ForegroundColor Green
Write-Host "Resultados esperados:" -ForegroundColor Cyan
Write-Host "- e2e-success-001: Debería estar guardado en MongoDB con datos enriquecidos" -ForegroundColor White
Write-Host "- e2e-inactive-001: Debería fallar en validación (cliente inactivo)" -ForegroundColor White  
Write-Host "- e2e-missing-001: Debería fallar en enriquecimiento (producto no existe)" -ForegroundColor White