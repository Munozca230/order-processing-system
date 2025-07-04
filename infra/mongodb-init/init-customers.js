// MongoDB initialization script for customers database
print('👥 Initializing customers database...');

// Connect to customers database  
db = db.getSiblingDB('catalog');

// Create customers collection with sample data
db.customers.deleteMany({}); // Clear existing data

const customers = [
  {
    customerId: 'customer-1',
    name: 'Juan Pérez García',
    email: 'juan.perez@email.com',
    phone: '+34 600 123 456',
    active: true,
    registrationDate: new Date('2023-01-15'),
    lastLogin: new Date(),
    preferences: {
      newsletter: true,
      notifications: true
    },
    address: {
      street: 'Calle Mayor 123',
      city: 'Madrid',
      postalCode: '28001',
      country: 'España'
    },
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    customerId: 'customer-2',
    name: 'María González López', 
    email: 'maria.gonzalez@email.com',
    phone: '+34 600 234 567',
    active: true,
    registrationDate: new Date('2023-02-20'),
    lastLogin: new Date(),
    preferences: {
      newsletter: true,
      notifications: false
    },
    address: {
      street: 'Avenida Libertad 456',
      city: 'Barcelona',
      postalCode: '08001',
      country: 'España'
    },
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    customerId: 'customer-3',
    name: 'Carlos Rodríguez Silva',
    email: 'carlos.rodriguez@email.com', 
    phone: '+34 600 345 678',
    active: false, // Inactive customer for testing
    registrationDate: new Date('2022-12-10'),
    lastLogin: new Date('2024-01-15'),
    preferences: {
      newsletter: false,
      notifications: false
    },
    address: {
      street: 'Plaza España 789',
      city: 'Valencia',
      postalCode: '46001',
      country: 'España'
    },
    deactivationReason: 'Account suspended due to inactivity',
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    customerId: 'customer-inactive',
    name: 'Cliente Inactivo',
    email: 'inactive@email.com',
    phone: '+34 600 456 789',
    active: false,
    registrationDate: new Date('2022-06-01'),
    lastLogin: new Date('2023-06-01'),
    preferences: {
      newsletter: false,
      notifications: false
    },
    address: {
      street: 'Calle Inactiva 000',
      city: 'Sevilla',
      postalCode: '41001',
      country: 'España'
    },
    deactivationReason: 'Customer requested account closure',
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    customerId: 'customer-premium',
    name: 'Ana Premium VIP',
    email: 'ana.premium@email.com',
    phone: '+34 600 567 890',
    active: true,
    registrationDate: new Date('2020-03-15'),
    lastLogin: new Date(),
    customerTier: 'VIP',
    preferences: {
      newsletter: true,
      notifications: true
    },
    address: {
      street: 'Paseo de la Castellana 100',
      city: 'Madrid',
      postalCode: '28046', 
      country: 'España'
    },
    loyaltyPoints: 15750,
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    customerId: 'customer-error',
    name: 'Cliente que causa error',
    email: 'error@test.com',
    phone: '+34 600 000 000',
    active: true,
    registrationDate: new Date(),
    lastLogin: new Date(),
    preferences: {
      newsletter: false,
      notifications: false
    },
    address: {
      street: 'Error Street 404',
      city: 'Test City',
      postalCode: '00000',
      country: 'Test Country'
    },
    testAccount: true, // Flag for error simulation
    createdAt: new Date(),
    updatedAt: new Date()
  }
];

// Insert customers
const result = db.customers.insertMany(customers);
print(`✅ Inserted ${result.insertedIds.length} customers`);

// Create indexes for performance
db.customers.createIndex({ customerId: 1 }, { unique: true });
db.customers.createIndex({ email: 1 }, { unique: true });
db.customers.createIndex({ active: 1 });
db.customers.createIndex({ customerTier: 1 });
db.customers.createIndex({ registrationDate: 1 });

print('📊 Created indexes: customerId (unique), email (unique), active, customerTier, registrationDate');

// Verify data
const activeCount = db.customers.countDocuments({ active: true });
const inactiveCount = db.customers.countDocuments({ active: false });
const totalCount = db.customers.countDocuments();

print(`🔢 Customer statistics:`);
print(`  - Total customers: ${totalCount}`);
print(`  - Active customers: ${activeCount}`);
print(`  - Inactive customers: ${inactiveCount}`);

print('🎉 Customers database initialization completed!');