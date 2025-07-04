#!/bin/bash

# Script para probar el sistema de reintentos exponenciales
echo "=== Prueba del Sistema de Reintentos Exponenciales ==="

# 1. Verificar que todos los servicios están corriendo
echo "1. Verificando servicios..."
docker-compose -f infra/docker-compose.yml ps

# 2. Enviar mensaje de prueba a Kafka
echo "2. Enviando mensaje de prueba a Kafka..."
docker-compose -f infra/docker-compose.yml exec kafka kafka-console-producer.sh \
  --broker-list localhost:9092 --topic orders << EOF
{"orderId":"test-001","customerId":"customer-1","products":["product-1","product-2"]}
EOF

# 3. Verificar logs del order-worker
echo "3. Verificando logs del order-worker (últimas 50 líneas)..."
docker-compose -f infra/docker-compose.yml logs --tail=50 order-worker

# 4. Verificar datos en Redis (mensajes fallidos)
echo "4. Verificando datos en Redis..."
docker-compose -f infra/docker-compose.yml exec redis redis-cli KEYS "*"

# 5. Verificar datos en MongoDB
echo "5. Verificando datos en MongoDB..."
docker-compose -f infra/docker-compose.yml exec mongo mongosh --eval "
  use orders;
  db.orders.find().pretty();
"

# 6. Simular error parando las APIs
echo "6. Simulando errores - parando product-api..."
docker-compose -f infra/docker-compose.yml stop product-api

echo "7. Enviando otro mensaje para triggear retry..."
docker-compose -f infra/docker-compose.yml exec kafka kafka-console-producer.sh \
  --broker-list localhost:9092 --topic orders << EOF
{"orderId":"test-002","customerId":"customer-2","products":["product-3"]}
EOF

# 8. Verificar logs de retry
echo "8. Verificando logs de retry (últimas 30 líneas)..."
docker-compose -f infra/docker-compose.yml logs --tail=30 order-worker

# 9. Verificar cola de retry en Redis
echo "9. Verificando cola de retry en Redis..."
docker-compose -f infra/docker-compose.yml exec redis redis-cli ZRANGE retry_queue 0 -1 WITHSCORES

# 10. Restaurar product-api
echo "10. Restaurando product-api..."
docker-compose -f infra/docker-compose.yml start product-api

echo "=== Prueba completada ==="
echo "Revisa los logs para ver los reintentos exponenciales en acción"