CREATE TABLE restaurant_payment_methods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    method VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(restaurant_id, method)
);

CREATE INDEX idx_restaurant_payment_methods_restaurant_id ON restaurant_payment_methods(restaurant_id);

