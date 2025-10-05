-- Calculated Estimates table to store estimate calculations
CREATE TABLE IF NOT EXISTS calculated_estimates (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id),
  total_price INTEGER NOT NULL,
  estimated_delivery_time_minutes INTEGER NOT NULL,
  estimate_data JSONB NOT NULL, -- Store the entire estimate request/response
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Orders table to store actual placed orders
CREATE TABLE IF NOT EXISTS orders (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id),
  calculated_estimate_id UUID NOT NULL REFERENCES calculated_estimates(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
