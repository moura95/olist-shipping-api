-- name: CreatePackage :one
INSERT INTO packages (tracking_code, product, weight_kg, destination_state, status)
VALUES ($1, $2, $3, $4, 'criado')
RETURNING id, tracking_code, product, weight_kg, destination_state, status, hired_carrier_id, hired_price, hired_delivery_days, created_at, updated_at;

-- name: GetPackageById :one
SELECT id, tracking_code, product, weight_kg, destination_state, status, hired_carrier_id, hired_price, hired_delivery_days, created_at, updated_at
FROM packages
WHERE id = $1;

-- name: GetPackageByTrackingCode :one
SELECT id, tracking_code, product, weight_kg, destination_state, status, hired_carrier_id, hired_price, hired_delivery_days, created_at, updated_at
FROM packages
WHERE tracking_code = $1;

-- name: ListPackages :many
SELECT id, tracking_code, product, weight_kg, destination_state, status, hired_carrier_id, hired_price, hired_delivery_days, created_at, updated_at
FROM packages
ORDER BY created_at DESC;

-- name: UpdatePackageStatus :exec
UPDATE packages
SET status = $2, updated_at = NOW()
WHERE id = $1;

-- name: HireCarrier :exec
UPDATE packages
SET hired_carrier_id = $2,
    hired_price = $3,
    hired_delivery_days = $4,
    status = 'esperando_coleta',
    updated_at = NOW()
WHERE id = $1;

-- name: DeletePackage :exec
DELETE FROM packages
WHERE id = $1;

-- name: TrackingCodeExists :one
SELECT EXISTS(
    SELECT 1 FROM packages
    WHERE tracking_code = $1
);

-- name: GetQuotesForPackage :many
SELECT
    c.name as carier,
    (cr.price_per_kg * $2) as estimated_price,
    cr.estimated_delivery_days
FROM carriers c
         JOIN carrier_regions cr ON c.id = cr.carrier_id
         JOIN states s ON s.region_id = cr.region_id
WHERE s.code = $1;

-- name: GetCarrierById :one
SELECT id, name, created_at
FROM carriers
WHERE id = $1;