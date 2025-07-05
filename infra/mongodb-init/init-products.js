// MongoDB initialization script for products catalog
print('🛍️ Initializing products catalog...');

// Connect to products database
db = db.getSiblingDB('catalog');

// Create products collection with sample data
db.products.deleteMany({}); // Clear existing data

const products = [
  {
    productId: 'product-1',
    name: 'Laptop Gaming MSI',
    description: 'High-performance gaming laptop with RTX graphics',
    price: 1299.99,
    category: 'laptops',
    stock: 15,
    active: true,
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    productId: 'product-2', 
    name: 'Mouse Gamer Logitech',
    description: 'Wireless gaming mouse with RGB lighting',
    price: 59.99,
    category: 'peripherals',
    stock: 50,
    active: true,
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    productId: 'product-3',
    name: 'Teclado Mecánico RGB',
    description: 'Mechanical keyboard with customizable RGB lighting',
    price: 129.99,
    category: 'peripherals', 
    stock: 30,
    active: true,
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    productId: 'product-4',
    name: 'Monitor 4K 27 pulgadas',
    description: 'Ultra HD 4K monitor for gaming and productivity',
    price: 399.99,
    category: 'monitors',
    stock: 20,
    active: true,
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    productId: 'product-5',
    name: 'Auriculares Gaming',
    description: 'Professional gaming headset with noise cancellation',
    price: 89.99,
    category: 'peripherals',
    stock: 40,
    active: true,
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    productId: 'product-6',
    name: 'SSD NVMe 1TB Samsung',
    description: 'High-speed NVMe SSD for ultra-fast data transfer',
    price: 149.99,
    category: 'storage',
    stock: 25,
    active: true,
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    productId: 'product-7',
    name: 'Webcam 4K Logitech',
    description: 'Ultra HD webcam for streaming and video calls',
    price: 199.99,
    category: 'peripherals',
    stock: 35,
    active: true,
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    productId: 'product-8',
    name: 'Tarjeta Gráfica RTX 4060',
    description: 'High-performance graphics card for gaming',
    price: 899.99,
    category: 'components',
    stock: 10,
    active: true,
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    productId: 'product-9',
    name: 'Silla Gaming Ergonómica',
    description: 'Ergonomic gaming chair with lumbar support',
    price: 299.99,
    category: 'furniture',
    stock: 15,
    active: true,
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    productId: 'product-error',
    name: 'Producto que causa error',
    description: 'Test product that simulates errors',
    price: 999.99,
    category: 'testing',
    stock: 0,
    active: false,
    createdAt: new Date(),
    updatedAt: new Date()
  }
];

// Insert products
const result = db.products.insertMany(products);
print(`✅ Inserted ${result.insertedIds.length} products`);

// Create indexes for performance
db.products.createIndex({ productId: 1 }, { unique: true });
db.products.createIndex({ category: 1 });
db.products.createIndex({ active: 1 });
db.products.createIndex({ price: 1 });

print('📊 Created indexes: productId (unique), category, active, price');

// Verify data
const count = db.products.countDocuments();
print(`🔢 Total products in catalog: ${count}`);

print('🎉 Products catalog initialization completed!');