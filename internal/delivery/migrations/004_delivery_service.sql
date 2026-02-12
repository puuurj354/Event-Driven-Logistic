-- Migration untuk Delivery Service Database
-- Database: db_delivery (untuk future implementation)
-- Note: Saat ini belum ada delivery-service di cmd/, tapi schema sudah disiapkan

-- Create ENUM type for shipment status
CREATE TYPE shipment_status AS ENUM ('PICKING_UP', 'ON_THE_WAY', 'DELIVERED', 'CANCELLED');

-- Create shipments table
CREATE TABLE IF NOT EXISTS shipments (
    id SERIAL PRIMARY KEY,
    order_id UUID NOT NULL UNIQUE,
    courier_name VARCHAR(255),
    current_lat FLOAT,
    current_long FLOAT,
    destination_lat FLOAT,
    destination_long FLOAT,
    status shipment_status NOT NULL DEFAULT 'PICKING_UP',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_shipments_order_id ON shipments(order_id);
CREATE INDEX idx_shipments_status ON shipments(status);
CREATE INDEX idx_shipments_created_at ON shipments(created_at DESC);

-- Create trigger to update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_shipments_updated_at
    BEFORE UPDATE ON shipments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
