# Verify MongoDB data
Write-Host "=== MongoDB Verification ====" -ForegroundColor Green

Set-Location "../infra"

# Start only infrastructure for verification
Write-Host "1. Starting services..." -ForegroundColor Yellow
docker-compose up -d

Start-Sleep -Seconds 15

# Check MongoDB data
Write-Host "2. Checking MongoDB data..." -ForegroundColor Yellow
docker-compose exec mongo mongosh --eval "use orders; db.orders.find().pretty(); print('Total orders: ' + db.orders.countDocuments());"

Write-Host "=== Verification Complete ====" -ForegroundColor Green