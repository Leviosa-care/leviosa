# ADR-0001: Modular Monolith Over Microservices

**Status**: Accepted  
**Date**: 2025-07

---

## Context

Leviosa needs to serve multiple distinct domains (auth, catalog, booking, notifications, settings) with clean separation of concerns. The team considered two primary deployment models: a distributed microservices architecture versus a modular monolith.

At this stage of the product, deployment complexity, operational overhead, and network latency between services are costs that outweigh the benefits of independent deployability.

## Decision

Deploy a single Go binary (`cmd/app`) containing five logically isolated service modules: `authuser`, `catalog`, `booking`, `settings`, and `notification`. Each module has its own domain layer, ports, application layer, and infrastructure adapters with no direct imports across module boundaries.

Inter-service calls happen through:
- **Synchronous**: typed port interfaces (e.g., `bookingPorts.AuthUserClient`) resolved at startup via DI container
- **Asynchronous**: RabbitMQ with typed message contracts for side effects (notifications, events)

## Consequences

**Positive:**
- Single deployment artifact simplifies CI/CD and infrastructure
- In-process calls eliminate network latency and serialization overhead between services
- Shared database with clear schema ownership per service
- Easier local development: one process to start

**Negative:**
- A crash in one module affects all modules
- All services must scale together (no independent scaling)
- Requires discipline to maintain module boundaries without import cycles

**Future path**: The port/adapter structure means individual services can be extracted to separate processes or containers if independent scaling or deployment becomes necessary. The interfaces already define the service contracts.
