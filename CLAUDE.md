# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Root Level Commands
- `make dev` - Run both frontend and backend in development mode
- `make up` - Start application with docker-compose (local development)
- `make down` - Stop docker-compose services
- `make help` - Show available commands
- `make help-front` - Show frontend-specific commands
- `make help-back` - Show backend-specific commands

### Frontend Commands (from /frontend)
- `npm run dev` - Start SvelteKit development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run check` - Type checking
- `npm run check:watch` - Type checking in watch mode

### Backend Commands (from /backend)
- `make run` - Start development server with all dependencies (uses Air for hot reload)
- `make deps` - Start required dependencies (Redis, PostgreSQL, RabbitMQ)
- `make start-redis` - Start Redis container
- `make start-postgres` - Start PostgreSQL container
- `make start-rabbit` - Start RabbitMQ container
- `make stop` - Stop all dependency containers
- `make clean` - Remove all dependency containers

### Testing Commands
Each microservice has specific test commands in their individual Makefiles:
- **authuser**: `make test`, `make ti` (integration), `make all`
- **catalog**: `make test`, `make ti` (integration), `make test-all`
- **settings**: `make test`, `make ti` (integration), `make all`
- **notification**: `make test`, `make ti` (integration)

## Project Architecture

### Full-Stack Application Structure
```
├── frontend/           # SvelteKit 5 application
├── backend/           # Go microservices monorepo
├── config/            # Infrastructure configuration (Caddy, Loki, Vault, etc.)
├── infra/             # Terraform and Ansible deployment scripts
└── compose.yaml       # Production docker-compose configuration
```

### Frontend Architecture (SvelteKit 5)
- **Framework**: SvelteKit 5 with TypeScript and Tailwind CSS v4
- **Authentication**: Session-based with cookies, environment-dependent (mock in dev/staging)
- **Route Structure**:
  - `(app)/` - Protected application routes
  - `auth/` - Authentication routes
  - `legal/` - Public legal pages
- **API Integration**: Automatic Bearer token injection via hooks.server.ts
- **Key Libraries**: bits-ui, sveltekit-superforms, arktype, @internationalized/date

### Backend Architecture (Go Microservices)
Go workspace with multiple microservices following hexagonal architecture:

```
backend/
├── core/              # Shared contracts, utilities, error definitions
├── authuser/          # Authentication and user management
├── catalog/           # Product catalog and pricing (Stripe integration)
├── settings/          # System configuration service
├── notification/      # Email/SMS notifications
└── cmd/leviosa/       # Main application entry point
```

**Hexagonal Architecture Pattern** for each microservice:
- `internal/domain/` - Business entities (no external dependencies)
- `internal/ports/` - Repository and service interfaces
- `internal/application/` - Use cases and workflows
- `internal/adapters/` - Infrastructure implementations:
  - `http/` - REST API handlers
  - `postgres/` - Database persistence
  - `rabbitmq/` - Message queue integration
  - `redis/` - Caching layer
  - `s3/` - Object storage

### Technology Stack

**Frontend**:
- SvelteKit 5, TypeScript, Tailwind CSS v4
- Node.js adapter for SSR

**Backend**:
- Go 1.24.2 with Air for hot reload
- PostgreSQL 17.5, Redis (alpine), RabbitMQ 3 (management)
- AWS S3, HashiCorp Vault
- External integrations: Stripe, Gmail SMTP

**Infrastructure**:
- Docker with docker-compose
- Caddy reverse proxy
- Grafana + Loki for logging
- Prometheus for monitoring

### Key Business Domains
- **User Management**: Authentication, OTP verification, sessions
- **Product Catalog**: Categories, products, pricing with Stripe
- **Settings**: Company configuration, OTP policies, notifications
- **Events**: Event creation with capacity tracking
- **Messaging**: User-to-user conversations
- **Payments**: Stripe integration for checkout

### Database Migrations
- Located in `core/migrations/` and `internal/migrations/`
- Convention: `{timestamp}_{service}_{action}_{entity_or_scope}.sql`
- Managed through Go embed system

### Error Handling
- Centralized error definitions in `core/errs/` package
- Domain-specific error constructors
- PostgreSQL and Redis error classification utilities
- Consistent error wrapping with context

### Development Workflow
1. Start dependencies: `make deps` (from backend directory)
2. Start frontend: `npm run dev` (from frontend directory)
3. Start backend: `make run` (from backend directory)
4. Access application: Frontend typically on http://localhost:5173
5. Run tests per service using individual makefiles

### Environment Configuration
- **Development**: Uses `.env` files and mock data
- **Production**: Environment variables with real services
- **Required Services**: PostgreSQL, Redis, RabbitMQ, AWS S3, Vault
- **Docker**: Full containerized setup available via compose.yaml

### Testing Approach
- **Backend**: Unit tests alongside source files, integration tests in `test/integration/`
- **Frontend**: No testing framework currently configured
- Uses real adapters for black-box integration testing
- Test data helpers available in `test/testdata/`

### CI/CD Pipeline (GitHub Actions)

**Workflow Strategy:** GitHub Flow with environment promotion

```
feature/new-feature
    ↓ Push to remote
    ↓ Create PR to main
    ↓ CI validation (tests, security scan)
main branch
    ↓ Merge triggers staging deployment
staging.leviosa.com (password-protected testing)
    ↓ Create git tag v*.*.*
production.leviosa.com (public users)
```

**Key Workflows:**

1. **CI Workflow** (`ci.yaml`) - Runs on all PRs to main:
   - Frontend: Build, unit tests
   - Backend: Build, unit tests, integration tests (with testcontainers)
   - Security: Go vulnerability scanning (govulncheck)
   - Blocks merge if any checks fail

2. **Staging Deployment** (`staging.yaml`) - Triggers on push to main:
   - Build and test frontend/backend
   - Scan Docker images with Trivy
   - Deploy to staging environment
   - Run health checks

3. **Production Deployment** (`production.yaml`) - Triggers on git tags (v*.*.*):
   - Build and test frontend/backend
   - Scan Docker images with Trivy
   - Deploy to production environment
   - Run health checks

**Security Features:**
- Docker image vulnerability scanning with Trivy (CRITICAL/HIGH threshold)
- Go dependency scanning with govulncheck
- Integration tests run with real dependencies (PostgreSQL, Redis, RabbitMQ, S3)
- Deployment blocked if security vulnerabilities found

**Merge Strategy:** Rebase and merge only (preserves full commit history)

See `.github/README.md` for detailed workflow documentation.