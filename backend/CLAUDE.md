# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Build and Run
- `make run` - Start development server with dependencies (Redis, PostgreSQL, RabbitMQ)
- `make deps` - Start all required dependencies (Redis, PostgreSQL, RabbitMQ)
- `make up` - Start with docker-compose
- `make down` - Stop docker-compose services

### Testing
Individual microservices have their own test commands in makefiles:
- **authuser/**: `make test` (specific OTP tests), `make ti` (integration tests), `make all` (all tests)
- **catalog/**: `make test` (specific price tests), `make ti` (integration tests), `make test-all` (all internal tests)
- **settings/**: `make test` (specific S3 tests), `make ti` (integration tests), `make all` (all integration tests)

### Dependencies Management
- `make start-redis` - Start Redis container on localhost:6379
- `make start-postgres` - Start PostgreSQL container on localhost:5432  
- `make start-rabbit` - Start RabbitMQ container on localhost:5672/15672
- `make stop` - Stop all containers
- `make clean` - Remove all containers

## Project Architecture

### Multi-Module Monorepo Structure
This is a Go workspace with multiple microservices following hexagonal architecture:

```
├── core/           # Shared contracts, utilities, error definitions
├── authuser/       # Authentication and user management service
├── catalog/        # Product catalog and pricing service
├── settings/       # System configuration service
├── notification/   # Email and SMS notification service
└── cmd/leviosa/    # Main application entry point
```

### Hexagonal Architecture Pattern
Each microservice follows ports & adapters pattern:
- **internal/domain/** - Business entities and value objects (no external dependencies)
- **internal/ports/** - Interfaces for repositories and services
- **internal/application/** - Use cases and business workflows
- **internal/adapters/** - Infrastructure implementations:
  - `http/` - REST API handlers and routes
  - `postgres/` - Database persistence
  - `rabbitmq/` - Message queue integration
  - `redis/` - Caching layer
  - `s3/` - Object storage

### Core Package
The `core/` directory contains shared components:
- **core/errs/** - Centralized error definitions and constructors
- **core/contracts/** - RabbitMQ message contracts and routing keys
- **core/messaging/** - RabbitMQ utilities for exchanges, queues, and payloads
- **core/ctxutil/** - Context utilities for logger and role validation
- **core/httpx/** - HTTP utilities for CORS, JSON responses, error handling
- **core/logger/** - Structured logging configuration
- **core/middleware/** - Authentication middleware
- **core/validation/** - Email and phone validation utilities

### Database Migrations
- Located in `core/migrations/` and `internal/migrations/`
- Convention: `{timestamp}_{service}_{action}_{entity_or_scope}.sql`
- Example: `20250714103022_catalog_add_column_products_buffer_time.sql`

### Testing Structure
- **Unit tests**: Alongside source files (`*_test.go`)
- **Integration tests**: `test/integration/` directories per service
- **Test data**: `test/testdata/` with helpers for database, RabbitMQ, HTTP setup
- Uses real adapters for black-box testing

### Technology Stack
- **Language**: Go 1.24.2
- **Database**: PostgreSQL 17.5
- **Cache**: Redis (alpine)
- **Message Queue**: RabbitMQ 3 (management)
- **Object Storage**: AWS S3
- **Development**: Air for hot reload
- **Configuration**: Viper + environment variables

### Environment Configuration
- Development uses `development.env` file
- Production uses environment variables
- Required services: PostgreSQL, Redis, RabbitMQ, AWS S3, Vault (HashiCorp)
- External integrations: Stripe, Gmail SMTP, Twilio SMS

### Error Handling
- Sentinel errors defined in `core/errs/` package
- Domain-specific error constructors (e.g., `NewNotFoundErr`, `NewConflictErr`)
- PostgreSQL and Redis error classification utilities
- Consistent error wrapping with context

### Key Business Domains
- **User Management**: Authentication, OTP verification, session management
- **Product Catalog**: Categories, products, pricing with Stripe integration
- **Settings**: Company configuration, OTP policies, notification preferences
- **Events**: Event creation and management with capacity tracking
- **Messaging**: User-to-user conversations and notifications
- **Payments**: Stripe integration for pricing and checkout
- **Booking**: Room scheduling, availability management, utilization analytics

### Booking Service Features

The booking service provides advanced scheduling optimization for massage therapy practices:

#### 1. Gap Detection API
**Purpose**: Identifies unused time slots between bookings to maximize room utilization

**Key Features**:
- Analyzes room schedules for a specific date
- Finds gaps before first booking, between bookings, and after last booking
- Suggests products that fit within each gap (Duration + BufferTime <= gap)
- Sorts suggestions by duration (shortest first for flexibility)
- Respects room operating hours

**Endpoints**:
- `GET /availabilities/rooms/{room_id}/gaps?date=YYYY-MM-DD`

**Algorithm**: See `internal/booking/application/availability/get_room_gaps.go`

#### 2. Utilization Metrics
**Purpose**: Track room efficiency and identify scheduling improvements

**Key Features**:
- Daily metrics stored in materialized view (`booking.room_daily_metrics`)
- Tracks utilization percentage, fragmentation count, idle minutes
- Calculates efficiency score (utilization - fragmentation penalty - idle penalty)
- Supports room-level and partner-level aggregation
- Historical analysis with date ranges

**Endpoints**:
- `GET /rooms/{room_id}/metrics?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD`
- `GET /partners/{partner_id}/metrics?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD`

**Database**: Materialized view uses CTEs for complex calculations
**Migration**: `20251124110535_booking_add_metrics_materialized_view.sql`

#### 3. Availability Suggestions
**Purpose**: Recommend optimal availability block durations based on products

**Key Features**:
- Analyzes partner's product catalog (Duration + BufferTime)
- Generates single and multi-session block suggestions (1x, 2x, 3x, 4x)
- For shared rooms, suggests standard durations (60, 90, 120, 180, 240 min)
- Priority ranking:
  * Priority 1: Standard durations (60, 90, 120, 180 min)
  * Priority 2: 30-minute increments
  * Priority 3: Non-standard durations
- Consolidates multiple products suggesting same duration

**Endpoints**:
- `GET /partners/{partner_id}/rooms/{room_id}/suggest-blocks`

**Algorithm**: See `internal/booking/application/availability/suggestion_algorithm.go`

#### Error Handling Pattern
All booking handlers use `httpx.RespondWithServiceError()` for automatic error classification:
- Repository layer: `errs.ClassifyPgError()` converts DB errors to sentinel errors
- Service layer: Wraps errors with `fmt.Errorf` or returns sentinel errors
- Handler layer: `RespondWithServiceError()` maps to HTTP status codes (400/404/503/500)

### Development Workflow
1. Start dependencies: `make deps`
2. Run development server: `make run` (uses Air for hot reload)
3. Run tests per service using individual makefiles
4. Database migrations managed through Go embed and custom migration system
