DELETE FROM carrier_regions;
DELETE FROM carriers;
DELETE FROM states;
DELETE FROM regions;

DROP INDEX IF EXISTS idx_states_region;
DROP INDEX IF EXISTS idx_packages_hired_carrier;
DROP INDEX IF EXISTS idx_carrier_regions_lookup;
DROP INDEX IF EXISTS idx_packages_estado_destino;
DROP INDEX IF EXISTS idx_packages_status;

DROP TABLE IF EXISTS packages;
DROP TABLE IF EXISTS carrier_regions;
DROP TABLE IF EXISTS carriers;
DROP TABLE IF EXISTS states;
DROP TABLE IF EXISTS regions;

DROP EXTENSION IF EXISTS "uuid-ossp";