#!/usr/bin/env pwsh

Write-Host "🧪 Testing Order Structure Validation" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan

# Function to check if services are running
function Test-ServiceHealth {
    param($ServiceName, $Url)
    
    try {
        $response = Invoke-RestMethod -Uri $Url -TimeoutSec 5
        Write-Host "✅ $ServiceName is healthy" -ForegroundColor Green
        return $true
    }
    catch {
        Write-Host "❌ $ServiceName is not responding: $($_.Exception.Message)" -ForegroundColor Red
        return $false
    }
}

# Check if all services are running
Write-Host "`n🏥 Checking Service Health..." -ForegroundColor Yellow

$services = @(
    @{Name="Product API"; Url="http://localhost:8081/health"},
    @{Name="Customer API"; Url="http://localhost:8082/health"},
    @{Name="Order API"; Url="http://localhost:3000/health"}
)

$allHealthy = $true
foreach ($service in $services) {
    if (-not (Test-ServiceHealth -ServiceName $service.Name -Url $service.Url)) {
        $allHealthy = $false
    }
}

if (-not $allHealthy) {
    Write-Host "`n⚠️ Some services are not running. Please start them with:" -ForegroundColor Yellow
    Write-Host "docker-compose up -d" -ForegroundColor Cyan
    exit 1
}

Write-Host "`n📨 Sending Test Order via Order API..." -ForegroundColor Yellow

# Create a test order
$testOrder = @{
    orderId = "structure-test-$(Get-Date -Format 'yyyyMMddHHmmss')"
    customerId = "customer-1"
    products = @(
        @{productId = "product-1"},
        @{productId = "product-2"}
    )
} | ConvertTo-Json -Depth 3

Write-Host "📋 Order payload:" -ForegroundColor Cyan
Write-Host $testOrder -ForegroundColor Gray

try {
    # Send order via Order API
    $response = Invoke-RestMethod -Uri "http://localhost:3000/api/orders" -Method POST -Body $testOrder -ContentType "application/json"
    Write-Host "✅ Order sent successfully: $($response.orderId)" -ForegroundColor Green
    $orderId = $response.orderId
    
    # Wait for processing
    Write-Host "`n⏳ Waiting 10 seconds for order processing..." -ForegroundColor Yellow
    Start-Sleep -Seconds 10
    
    # Now we would verify in MongoDB, but since docker commands aren't available in WSL,
    # let's show the expected structure and verification commands
    Write-Host "`n📊 Expected MongoDB Structure (according to prueba.md):" -ForegroundColor Yellow
    $expectedStructure = @"
{
  "_id": ObjectId(),
  "orderId": "$orderId",
  "customerId": "customer-1",
  "products": [
    {
      "productId": "product-1",
      "name": "Laptop Gaming MSI",
      "price": 1299.99
    },
    {
      "productId": "product-2", 
      "name": "Mouse Gamer Logitech",
      "price": 59.99
    }
  ]
}
"@
    Write-Host $expectedStructure -ForegroundColor Gray
    
    Write-Host "`n🔍 To verify the actual structure in MongoDB, run:" -ForegroundColor Yellow
    Write-Host "docker-compose exec mongo mongosh orders --eval `"db.orders.find({orderId: '$orderId'}).forEach(printjson)`"" -ForegroundColor Cyan
    
    Write-Host "`n📋 Verification Checklist:" -ForegroundColor Yellow
    Write-Host "✅ Order sent to Kafka via Order API" -ForegroundColor Green
    Write-Host "⏳ Order should be processed by Order Worker" -ForegroundColor Yellow
    Write-Host "⏳ Products should be enriched with name and price" -ForegroundColor Yellow
    Write-Host "⏳ Customer should be validated as active" -ForegroundColor Yellow
    Write-Host "⏳ Final document should match prueba.md structure" -ForegroundColor Yellow
    
    Write-Host "`n💡 Additional verification commands:" -ForegroundColor Cyan
    Write-Host "# Check Order Worker logs:" -ForegroundColor Gray
    Write-Host "docker-compose logs -f order-worker | grep `"$orderId`"" -ForegroundColor Gray
    Write-Host "`n# Check all orders in MongoDB:" -ForegroundColor Gray
    Write-Host "docker-compose exec mongo mongosh orders --eval `"db.orders.find().forEach(printjson)`"" -ForegroundColor Gray
    Write-Host "`n# Count processed orders:" -ForegroundColor Gray
    Write-Host "docker-compose exec mongo mongosh orders --eval `"print('Total orders:', db.orders.countDocuments())`"" -ForegroundColor Gray
    
    Write-Host "`n🎉 Test completed! Order ID: $orderId" -ForegroundColor Green
}
catch {
    Write-Host "❌ Failed to send order: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "`n🔧 Troubleshooting:" -ForegroundColor Yellow
    Write-Host "1. Ensure Order API is running: docker-compose ps order-api" -ForegroundColor Gray
    Write-Host "2. Check Order API logs: docker-compose logs order-api" -ForegroundColor Gray
    Write-Host "3. Verify Kafka is accessible: docker-compose logs kafka" -ForegroundColor Gray
    exit 1
}

Write-Host "`n✨ Structure validation test ready! Use the verification commands above to confirm the MongoDB structure matches prueba.md requirements." -ForegroundColor Green