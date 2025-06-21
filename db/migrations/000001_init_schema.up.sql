SET TIME ZONE 'America/Sao_Paulo';

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Table Regions
CREATE TABLE regions (
                         id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                         name VARCHAR(50) NOT NULL UNIQUE,
                         created_at TIMESTAMP DEFAULT NOW()
);

-- Table States
CREATE TABLE states (
                        code CHAR(2) PRIMARY KEY,
                        name VARCHAR(100) NOT NULL,
                        region_id UUID NOT NULL,
                        created_at TIMESTAMP DEFAULT NOW(),
                        CONSTRAINT fk_region FOREIGN KEY (region_id) REFERENCES regions(id)
);

-- Table Carriers
CREATE TABLE carriers (
                          id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                          name VARCHAR(255) NOT NULL,
                          created_at TIMESTAMP DEFAULT NOW()
);

-- Table Carrier Regions
CREATE TABLE carrier_regions (
                                 id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                 carrier_id UUID NOT NULL,
                                 region_id UUID NOT NULL,
                                 estimated_delivery_days INT NOT NULL,
                                 price_per_kg DECIMAL(10,2) NOT NULL,
                                 created_at TIMESTAMP DEFAULT NOW(),
                                 CONSTRAINT fk_carrier FOREIGN KEY (carrier_id) REFERENCES carriers(id) ON DELETE CASCADE,
                                 CONSTRAINT fk_region FOREIGN KEY (region_id) REFERENCES regions(id),
                                 UNIQUE(carrier_id, region_id)
);

-- Table Packages
CREATE TABLE packages (
                          id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                          tracking_code VARCHAR(11) UNIQUE NOT NULL,
                          product VARCHAR(255) NOT NULL,
                          weight_kg FLOAT NOT NULL CHECK (weight_kg > 0),
                          destination_state CHAR(2) NOT NULL,
                          status VARCHAR(50) NOT NULL DEFAULT 'criado',
                          hired_carrier_id UUID,
                          hired_price DECIMAL(10,2),
                          hired_delivery_days INT,
                          created_at TIMESTAMP DEFAULT NOW(),
                          updated_at TIMESTAMP DEFAULT NOW(),

                          CONSTRAINT fk_destination_state FOREIGN KEY (destination_state) REFERENCES states(code),
                          CONSTRAINT fk_hired_carrier FOREIGN KEY (hired_carrier_id) REFERENCES carriers(id),
                          CONSTRAINT check_status CHECK (status IN ('criado', 'esperando_coleta', 'coletado', 'enviado', 'entregue', 'extraviado'))
);

-- Indexes
CREATE INDEX idx_packages_status ON packages(status);
CREATE INDEX idx_packages_destination_state ON packages(destination_state);
CREATE INDEX idx_packages_tracking_code ON packages(tracking_code);
CREATE INDEX idx_carrier_regions_lookup ON carrier_regions(region_id, carrier_id);
CREATE INDEX idx_packages_hired_carrier ON packages(hired_carrier_id);
CREATE INDEX idx_states_region ON states(region_id);