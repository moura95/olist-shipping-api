-- name: ListStates :many
SELECT s.code, s.name, r.name as region_name
FROM states s
         JOIN regions r ON s.region_id = r.id
ORDER BY s.name;

-- name: GetStateByCode :one
SELECT s.code, s.name, r.name as region_name
FROM states s
         JOIN regions r ON s.region_id = r.id
WHERE s.code = $1;

-- name: ListRegions :many
SELECT id, name, created_at
FROM regions
ORDER BY name;