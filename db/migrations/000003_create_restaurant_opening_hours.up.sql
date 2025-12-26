CREATE TABLE restaurant_opening_hours (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    weekday INTEGER NOT NULL CHECK (weekday >= 0 AND weekday <= 6),
    opens_at INTEGER NOT NULL CHECK (opens_at >= 0 AND opens_at < 1440),
    closes_at INTEGER NOT NULL CHECK (closes_at >= 0 AND closes_at < 1440),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_restaurant_opening_hours_restaurant_id ON restaurant_opening_hours(restaurant_id);
CREATE INDEX idx_restaurant_opening_hours_weekday ON restaurant_opening_hours(restaurant_id, weekday);

