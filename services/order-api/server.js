const express = require('express');
const cors = require('cors');
const { Kafka } = require('kafkajs');
const fetch = require('node-fetch');

const app = express();
const port = process.env.PORT || 3000;

// Middleware
app.use(cors());
app.use(express.json());

// Kafka configuration
const kafka = new Kafka({
  clientId: 'order-api',
  brokers: [process.env.KAFKA_BROKERS || 'localhost:9092'],
});

const producer = kafka.producer();

// Initialize Kafka producer
async function initKafka() {
  try {
    await producer.connect();
    console.log('✅ Connected to Kafka');
  } catch (error) {
    console.error('❌ Failed to connect to Kafka:', error);
    process.exit(1);
  }
}

// Health check endpoint
app.get('/health', (req, res) => {
  res.json({ 
    status: 'healthy', 
    service: 'order-api',
    timestamp: new Date().toISOString()
  });
});

// Submit order endpoint
app.post('/api/orders', async (req, res) => {
  try {
    const { orderId, customerId, products } = req.body;

    // Validate request
    if (!orderId || !customerId || !products || !Array.isArray(products)) {
      return res.status(400).json({
        error: 'Invalid request',
        message: 'orderId, customerId, and products array are required'
      });
    }

    // Validate products structure
    for (const product of products) {
      if (!product.productId) {
        return res.status(400).json({
          error: 'Invalid product',
          message: 'Each product must have a productId'
        });
      }
    }

    const orderMessage = {
      orderId,
      customerId,
      products
    };

    console.log(`📨 Sending order to Kafka: ${orderId}`);

    // Send to Kafka
    await producer.send({
      topic: 'orders',
      messages: [
        {
          key: orderId,
          value: JSON.stringify(orderMessage),
        },
      ],
    });

    console.log(`✅ Order sent successfully: ${orderId}`);

    res.json({
      success: true,
      message: 'Order submitted successfully',
      orderId: orderId,
      timestamp: new Date().toISOString()
    });

  } catch (error) {
    console.error('❌ Error sending order:', error);
    res.status(500).json({
      success: false,
      error: 'Failed to submit order',
      message: error.message
    });
  }
});

// Get orders endpoint (placeholder - would query MongoDB in real implementation)
app.get('/api/orders', (req, res) => {
  res.json({
    message: 'Order retrieval not implemented - check MongoDB directly',
    mongoCommand: 'docker-compose exec mongo mongosh orders --eval "db.orders.find().forEach(printjson)"'
  });
});

// Get specific order endpoint
app.get('/api/orders/:orderId', (req, res) => {
  res.json({
    message: 'Order retrieval not implemented - check MongoDB directly',
    orderId: req.params.orderId,
    mongoCommand: `docker-compose exec mongo mongosh orders --eval "db.orders.find({orderId: '${req.params.orderId}'}).forEach(printjson)"`
  });
});

// Get order status endpoint (proxy to Order Worker)
app.get('/api/orders/:orderId/status', async (req, res) => {
  try {
    const response = await fetch(`http://order-worker:8080/api/orders/${req.params.orderId}/status`, {
      method: 'GET',
      headers: {
        'Cache-Control': 'no-cache, no-store, must-revalidate',
        'Pragma': 'no-cache',
        'Expires': '0'
      }
    });
    
    if (!response.ok) {
      if (response.status === 404) {
        return res.status(404).json({
          error: 'Order not found',
          orderId: req.params.orderId
        });
      }
      return res.status(response.status).json({
        error: 'Unable to retrieve order status'
      });
    }
    
    const statusData = await response.json();
    
    // Add no-cache headers to response
    res.set({
      'Cache-Control': 'no-cache, no-store, must-revalidate',
      'Pragma': 'no-cache',
      'Expires': '0'
    });
    
    res.json(statusData);
  } catch (error) {
    console.error('❌ Error fetching order status:', error);
    res.status(503).json({ 
      error: 'Order status service unavailable',
      message: 'Please try again later'
    });
  }
});

// Start server
app.listen(port, async () => {
  console.log(`🚀 Order API running on port ${port}`);
  await initKafka();
});

// Graceful shutdown
process.on('SIGTERM', async () => {
  console.log('📋 Shutting down gracefully...');
  await producer.disconnect();
  process.exit(0);
});

process.on('SIGINT', async () => {
  console.log('📋 Shutting down gracefully...');
  await producer.disconnect();
  process.exit(0);
});