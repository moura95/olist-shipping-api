-- name: ListCarriers :many
SELECT id, name, created_at
FROM carriers
ORDER BY name;

-- name: GetCarrierRegions :many
SELECT
    cr.id,
    cr.carrier_id,
    cr.region_id,
    cr.estimated_delivery_days,
    cr.price_per_kg,
    r.name as region_name
FROM carrier_regions cr
         JOIN regions r ON cr.region_id = r.id
WHERE cr.carrier_id = $1;

-- name: GetRegionByState :one
SELECT r.id, r.name
FROM regions r
         JOIN states s ON s.region_id = r.id
WHERE s.code = $1;