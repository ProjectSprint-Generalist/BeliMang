-- Install PostGIS and pgcrypto
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- User Role enum
CREATE TYPE user_role AS ENUM ('user', 'admin');

-- User table
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  username TEXT NOT NULL,
  password TEXT NOT NULL,
  email TEXT NOT NULL,
  role user_role NOT NULL DEFAULT 'user'
);
-- Unique constraint on username for all roles
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users (username);
-- Unique constraint on email per role (allows same email for different roles)
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_role_email ON users (role, email);

-- User Images table
CREATE TABLE IF NOT EXISTS images (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  filename TEXT NOT NULL,
  url TEXT NOT NULL,
  size_bytes BIGINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Merchant Enum
CREATE TYPE merchant_category AS ENUM (
  'SmallRestaurant',
  'MediumRestaurant',
  'LargeRestaurant',
  'MerchandiseRestaurant',
  'BoothKiosk',
  'ConvenienceStore'
);

-- Merchant table
CREATE TABLE IF NOT EXISTS merchants (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  merchant_category merchant_category NOT NULL,
  location GEOMETRY(POINT, 4326) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Product Enum
CREATE TYPE product_category AS ENUM (
  'Beverage',
  'Food',
  'Snack',
  'Condiments',
  'Additions'
);

-- Merchant Items table
CREATE TABLE IF NOT EXISTS merchant_items (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  merchant_id UUID NOT NULL REFERENCES merchants(id),
  name TEXT NOT NULL,
  product_category product_category NOT NULL,
  price INTEGER NOT NULL CHECK (price >= 1),
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

