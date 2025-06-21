# Olist Shipping API

**API para gerenciamento de entregas da Olist** - Sistema completo para cotação de fretes, controle de pacotes e integração com transportadoras.

## 📋 Descrição

Aplicação desenvolvida para gerenciar envios de pacotes por diferentes transportadoras, permitindo registrar pacotes, consultar status de envio, simular custos com base em peso e região, e contratar a melhor transportadora para realizar a entrega.

## 🚀 Demo

- **API:** [https://olist-shipping-api.run.app](https://olist-shipping-api.run.app)
- **Frontend:** [https://olist-shipping-front.vercel.app](https://olist-shipping-front.vercel.app)
- **Swagger:** [https://olist-shipping-api.run.app/swagger/index.html](https://olist-shipping-api.run.app/swagger/index.html)

![Swagger Documentation](docs/swagger-preview.png)

## 📋 Índice

- [Requisitos](#requisitos)
- [Tecnologias Utilizadas](#tecnologias-utilizadas)
- [Instalação](#instalacao)
- [Como Rodar](#como-rodar)
- [Testes](#testes)
- [Endpoints](#endpoints)
- [Regras de Negócio](#regras-de-negocio)
- [Arquitetura](#arquitetura)
- [CI/CD](#cicd)
- [Autor](#autor)


## 🛠️ Tecnologias Utilizadas

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=flat&logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=flat&logo=docker&logoColor=white)
![GitHub Actions](https://img.shields.io/badge/GitHub%20Actions-2088FF?style=flat&logo=github-actions&logoColor=white)
![Google Cloud](https://img.shields.io/badge/Google%20Cloud-4285F4?style=flat&logo=google-cloud&logoColor=white)

### Backend
- **Golang 1.22** - Linguagem principal
- **Gin** - Framework web
- **PostgreSQL** - Banco de dados
- **SQLC** - Geração de código type-safe para SQL
- **Sqlx** - Driver PostgreSQL com melhor performance
- **Viper** - Gerenciamento de configurações
- **Zap** - Logging estruturado
- **Validator** - Validação de dados
- **UUID** - Identificadores únicos

### DevOps & Infraestrutura
- **Docker & Docker Compose** - Containerização
- **GitHub Actions** - CI/CD
- **Google Cloud Run** - Deploy e hospedagem
- **TestContainers** - Testes de integração
- **Golang-Migrate** - Migrations de banco
- **Swagger** - Documentação da API

### Testes
- **Testify** - Framework de testes
- **Mocks (Mockery)** - Testes unitários
- **TestContainers** - Testes de integração

## 📋 Requisitos

Para rodar o projeto localmente, você precisa ter instalado:

- **Go 1.22+**
- **Docker & Docker Compose**
- **Make** (opcional, mas recomendado)
- **Golang-Migrate** (para migrations)
- **SQLC** (para geração de código)

### Instalacao das Ferramentas

#### Golang-Migrate
```bash
# MacOS
brew install golang-migrate

# Ubuntu/Debian
apt install migrate

# Windows (chocolatey)
choco install migrate
```

#### SQLC
```bash
# MacOS
brew install sqlc

# Go install
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

## 🚀 Instalacao

### 1. Clone o projeto
```bash
git clone https://github.com/moura95/olist-shipping-api
cd olist-shipping-api
```

### 2. Instale as dependências
```bash
go mod tidy
```

### 3. Configure o ambiente
```bash
cp .envexample .env
# Edite o arquivo .env com suas configurações
```

## 💻 Como Rodar

### Opção 1: Com Docker (Recomendado)
```bash
# Inicia banco + migrations + aplicação
make start

# Ou manualmente
docker-compose up -d
make migrate-up
go run cmd/main.go
```

### Opção 2: Desenvolvimento Local
```bash
# 1. Suba apenas o banco
docker-compose up -d psql

# 2. Execute as migrations
make migrate-up

# 3. Rode a aplicação
make run
# ou
go run cmd/main.go
```

### 🧪 Testes

```bash
# Rodar todos os testes
make test

# Testes unitários apenas
make test-unit

# Testes de integração
make test-integration

# Testes por camada
make test-service
make test-repository
```

## 📚 Endpoints

### Pacotes
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `POST` | `/api/v1/packages` | Criar novo pacote |
| `GET` | `/api/v1/packages` | Listar todos os pacotes |
| `GET` | `/api/v1/packages/{id}` | Buscar pacote por ID |
| `GET` | `/api/v1/packages/tracking/{code}` | Buscar por código de rastreio |
| `PATCH` | `/api/v1/packages/{id}/status` | Atualizar status do pacote |
| `POST` | `/api/v1/packages/{id}/hire` | Contratar transportadora |
| `DELETE` | `/api/v1/packages/{id}` | Deletar pacote |

### Cotações
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/api/v1/quotes?estado_destino=SP&peso_kg=2.0` | Obter cotações de frete |

### Informações
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/api/v1/carriers` | Listar transportadoras |
| `GET` | `/api/v1/states` | Listar estados brasileiros |

### Exemplos de Uso

#### Criar um Pacote
```bash
curl -X POST http://localhost:8080/api/v1/packages \
  -H "Content-Type: application/json" \
  -d '{
    "produto": "Camisa tamanho G",
    "peso_kg": 0.6,
    "estado_destino": "PR"
  }'
```

#### Cotação de Frete
```bash
curl "http://localhost:8080/api/v1/quotes?estado_destino=SP&peso_kg=2.0"
```

#### Contratar Transportadora
```bash
curl -X POST http://localhost:8080/api/v1/packages/{id}/hire \
  -H "Content-Type: application/json" \
  -d '{
    "transportadora_id": "660e8400-e29b-41d4-a716-446655440001",
    "preco": "25.90",
    "prazo_dias": 5
  }'
```

## 🏗️ Regras de Negocio

### Transportadoras Disponíveis

#### Nebulix Logística
- **Regiões:** Sul e Sudeste
- **Prazo:** 4 dias
- **Preço:** R$ 5,90/kg

#### RotaFácil Transportes
- **Sul/Sudeste:** 7 dias - R$ 4,35/kg
- **Centro-Oeste:** 9 dias - R$ 6,22/kg
- **Nordeste:** 13 dias - R$ 8,00/kg

#### Moventra Express
- **Centro-Oeste:** 7 dias - R$ 7,30/kg
- **Nordeste:** 10 dias - R$ 9,50/kg

### Status dos Pacotes
- `criado` → `esperando_coleta` → `coletado` → `enviado` → `entregue`
- `extraviado` (status especial)

### Cálculo de Preços
```
Preço Final = Peso (kg) × Preço por KG da Transportadora
```

## 🏛️ Arquitetura

O projeto segue uma arquitetura com separação clara de responsabilidades:

```
cmd/                    # Entry points
├── main.go

internal/               # Código da aplicação
├── handler/           # HTTP handlers (controllers)
├── service/           # Regras de negócio
├── repository/        # Acesso a dados (gerado pelo SQLC)
├── middleware/        # Middlewares (CORS, rate limit, logging)
└── server.go         # Setup do servidor

api/v1/               # Tipos da API
├── packages.go       # Request/Response types
└── response.go       # Helpers de resposta

config/               # Configurações
pkg/                  # Utilities
├── validator/        # Validações customizadas
├── ginx/            # Helpers do Gin
└── tracking/        # Geração de códigos

db/                   # Database
├── migrations/       # SQL migrations
└── queries/         # SQL queries (SQLC)

tests/               # Testes organizados
├── repository/      # Testes de repository
└── service/        # Testes unitários e integração
```

## 🔄 CI/CD

O projeto possui pipeline completa com GitHub Actions:

### Pipeline Stages
1. **Test** - Execução de testes
2. **Build** - Compilação da aplicação
3. **Docker Build** - Criação da imagem
4. **Deploy** - Deploy automático no Google Cloud Run

### Deploy
- **Ambiente:** Google Cloud Run
- **Trigger:** Push na branch `main`
- **Database:** PostgreSQL dedicado
- **Monitoring:** Health checks automáticos

## 📁 Collection Postman

Importe a collection para testar a API:
```
docs/olist_shipping_collection.json
```

## 📊 Métricas Disponíveis

- **Total de Pacotes por Status**
- **Pacotes por Transportadora**
- **Tempo Médio de Entrega**
- **Cotações por Região**
- **Volume de Entregas por Estado**

## 🔧 Comandos Úteis

```bash
# Desenvolvimento
make run              # Roda a aplicação
make test            # Executa todos os testes
make migrate-up      # Aplica migrations
make migrate-down    # Reverte migrations

# Docker
make up              # Sobe ambiente completo
make down            # Para ambiente
make restart         # Reinicia ambiente

# Database
make sqlc            # Gera código SQLC
make migrate-create  # Cria nova migration

# Documentação
make swag            # Gera Swagger docs
```

## 👨‍💻 Autor

**Guilherme Moura** - *Engenheiro de Software*
- GitHub: [@moura95](https://github.com/moura95)
- LinkedIn: [Guilherme Moura](https://linkedin.com/in/guilherme-moura95)
- Email: dev@guilhermemoura.dev

---
