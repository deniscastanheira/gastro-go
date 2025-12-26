-- name: CreateRestaurant :one
INSERT INTO restaurants (
    name, slug, description, status, category, rating, total_reviews,
    delivery_fee, min_order_value, preparation_time_min,
    supports_pickup, supports_delivery, logo_url, banner_url
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
) RETURNING *;

-- name: GetRestaurantByID :one
SELECT * FROM restaurants WHERE id = $1 LIMIT 1;

-- name: GetRestaurantBySlug :one
SELECT * FROM restaurants WHERE slug = $1 LIMIT 1;

-- name: ListRestaurants :many
SELECT * FROM restaurants
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateRestaurantStatus :one
UPDATE restaurants
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: CreateRestaurantAddress :one
INSERT INTO restaurant_addresses (
    restaurant_id, street, number, complement, city, state, zip_code, lat, lng
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: UpdateRestaurantAddress :one
UPDATE restaurant_addresses
SET street = $2, number = $3, complement = $4, city = $5, state = $6, zip_code = $7, lat = $8, lng = $9, updated_at = NOW()
WHERE restaurant_id = $1
RETURNING *;

-- name: GetRestaurantAddress :one
SELECT * FROM restaurant_addresses WHERE restaurant_id = $1 LIMIT 1;

-- name: CreateOpeningHour :one
INSERT INTO restaurant_opening_hours (
    restaurant_id, weekday, opens_at, closes_at
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: DeleteOpeningHoursByRestaurant :exec
DELETE FROM restaurant_opening_hours WHERE restaurant_id = $1;

-- name: GetOpeningHoursByRestaurant :many
SELECT * FROM restaurant_opening_hours
WHERE restaurant_id = $1
ORDER BY weekday, opens_at;

-- name: CreatePaymentMethod :one
INSERT INTO restaurant_payment_methods (
    restaurant_id, method
) VALUES (
    $1, $2
) RETURNING *;

-- name: DeletePaymentMethodsByRestaurant :exec
DELETE FROM restaurant_payment_methods WHERE restaurant_id = $1;

-- name: GetPaymentMethodsByRestaurant :many
SELECT * FROM restaurant_payment_methods
WHERE restaurant_id = $1
ORDER BY method;

