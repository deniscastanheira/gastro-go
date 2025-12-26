CREATE TABLE restaurants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    category VARCHAR(100),
    rating INTEGER NOT NULL DEFAULT 0,
    total_reviews INTEGER NOT NULL DEFAULT 0,
    delivery_fee BIGINT NOT NULL DEFAULT 0,
    min_order_value BIGINT NOT NULL DEFAULT 0,
    preparation_time_min INTEGER NOT NULL DEFAULT 0,
    supports_pickup BOOLEAN NOT NULL DEFAULT false,
    supports_delivery BOOLEAN NOT NULL DEFAULT false,
    logo_url TEXT,
    banner_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

