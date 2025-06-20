INSERT INTO regions (id, name) VALUES
                                   ('550e8400-e29b-41d4-a716-446655440001', 'Sul'),
                                   ('550e8400-e29b-41d4-a716-446655440002', 'Sudeste'),
                                   ('550e8400-e29b-41d4-a716-446655440003', 'Centro-Oeste'),
                                   ('550e8400-e29b-41d4-a716-446655440004', 'Nordeste'),
                                   ('550e8400-e29b-41d4-a716-446655440005', 'Norte');

-- States
INSERT INTO states (sigla, name, region_id) VALUES
-- Sul
('RS', 'Rio Grande do Sul', '550e8400-e29b-41d4-a716-446655440001'),
('SC', 'Santa Catarina', '550e8400-e29b-41d4-a716-446655440001'),
('PR', 'Paraná', '550e8400-e29b-41d4-a716-446655440001'),
-- Sudeste
('SP', 'São Paulo', '550e8400-e29b-41d4-a716-446655440002'),
('RJ', 'Rio de Janeiro', '550e8400-e29b-41d4-a716-446655440002'),
('MG', 'Minas Gerais', '550e8400-e29b-41d4-a716-446655440002'),
('ES', 'Espírito Santo', '550e8400-e29b-41d4-a716-446655440002'),
-- Centro-Oeste
('GO', 'Goiás', '550e8400-e29b-41d4-a716-446655440003'),
('MT', 'Mato Grosso', '550e8400-e29b-41d4-a716-446655440003'),
('MS', 'Mato Grosso do Sul', '550e8400-e29b-41d4-a716-446655440003'),
('DF', 'Distrito Federal', '550e8400-e29b-41d4-a716-446655440003'),
-- Nordeste
('BA', 'Bahia', '550e8400-e29b-41d4-a716-446655440004'),
('SE', 'Sergipe', '550e8400-e29b-41d4-a716-446655440004'),
('AL', 'Alagoas', '550e8400-e29b-41d4-a716-446655440004'),
('PE', 'Pernambuco', '550e8400-e29b-41d4-a716-446655440004'),
('PB', 'Paraíba', '550e8400-e29b-41d4-a716-446655440004'),
('RN', 'Rio Grande do Norte', '550e8400-e29b-41d4-a716-446655440004'),
('CE', 'Ceará', '550e8400-e29b-41d4-a716-446655440004'),
('PI', 'Piauí', '550e8400-e29b-41d4-a716-446655440004'),
('MA', 'Maranhão', '550e8400-e29b-41d4-a716-446655440004'),
-- Norte
('AC', 'Acre', '550e8400-e29b-41d4-a716-446655440005'),
('RO', 'Rondônia', '550e8400-e29b-41d4-a716-446655440005'),
('AM', 'Amazonas', '550e8400-e29b-41d4-a716-446655440005'),
('RR', 'Roraima', '550e8400-e29b-41d4-a716-446655440005'),
('PA', 'Pará', '550e8400-e29b-41d4-a716-446655440005'),
('AP', 'Amapá', '550e8400-e29b-41d4-a716-446655440005'),
('TO', 'Tocantins', '550e8400-e29b-41d4-a716-446655440005');

INSERT INTO carriers (id, name) VALUES
                                    ('660e8400-e29b-41d4-a716-446655440001', 'Nebulix Logística'),
                                    ('660e8400-e29b-41d4-a716-446655440002', 'RotaFácil Transportes'),
                                    ('660e8400-e29b-41d4-a716-446655440003', 'Moventra Express');

-- Carrier Regions
INSERT INTO carrier_regions (carrier_id, region_id, prazo_estimado_dias, preco_por_kg) VALUES
-- Nebulix: Sul e Sudeste
('660e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440001', 4, 5.90), -- Sul
('660e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440002', 4, 5.90), -- Sudeste

-- RotaFácil: Sul, Sudeste, Centro-Oeste, Nordeste
('660e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440001', 7, 4.35),  -- Sul
('660e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440002', 7, 4.35),  -- Sudeste
('660e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440003', 9, 6.22),  -- Centro-Oeste
('660e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440004', 13, 8.00), -- Nordeste

-- Moventra: Centro-Oeste, Nordeste
('660e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440003', 7, 7.30),  -- Centro-Oeste
('660e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440004', 10, 9.50); -- Nordeste

-- Sample packages for testing
INSERT INTO packages (id, tracking_code, produto, peso_kg, estado_destino, status) VALUES
('770e8400-e29b-41d4-a716-446655440001', 'BR72619195', 'Camisa tamanho G', 0.6, 'PR', 'criado'),
('770e8400-e29b-41d4-a716-446655440002', 'BR38897894', 'Notebook Dell', 2.5, 'SP', 'enviado'),
('770e8400-e29b-41d4-a716-446655440003', 'BR14506220', 'Livro de programação', 0.8, 'RJ', 'coletado');