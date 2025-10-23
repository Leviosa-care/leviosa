# 🧭 Expose a Minimal Catalog API in the Core Microservice

Implement a **minimal read-only Catalog API** inside the **Core microservice** to provide shared access to product and category data for other microservices (e.g. `authuser`).

---

## 🎯 Goal

Create a lightweight interface in the Core microservice that exposes:

- `ListCategories()`
- `ListProducts()`

and defines the corresponding shared data models:

- `Category`
- `Product`

This API will allow other microservices to retrieve catalog data without directly depending on the Catalog service.

---

## 🧠 Motivation

- **Decoupling:** Other services, like `authuser`, can access catalog data without creating a tight dependency on the Catalog service.  
- **Single source of truth:** Catalog remains the owner of the data, Core only exposes it for reading.  
- **Consistency:** Shared domain models (`Category`, `Product`) ensure all services interpret the data the same way.  
- **Scalability:** In the future, Core can implement its own repository, caching, or replication strategies without affecting the Catalog service.  
- **Minimal surface area:** Only expose what is needed (`ListCategories` and `ListProducts`), keeping the service simple and maintainable.

---

## 🧩 Implementation Notes

- The Core microservice will **reuse the same database tables** currently managed by the Catalog microservice for now.
- In the future, the Core service may evolve to use its **own repository implementation** (e.g. via event-driven sync or read replicas).
- Core must remain **read-only** — it should never modify catalog data.
- Expose a **clean, minimal repository interface**, for example:

```go
  type CatalogRepository interface {
      ListCategories(ctx context.Context) ([]domain.Category, error)
      ListProducts(ctx context.Context) ([]domain.Product, error)
  }
```
- Implement this interface using the existing Catalog database until decoupling is complete.
- Core’s domain models (Category, Product) should be defined in a domain package and reusable across services.

## ⚙️ Design Constraints

- Keep Core’s surface minimal — only expose what’s needed for cross-service data access.
- No business logic, filtering, or write operations.
- Preserve clear ownership boundaries: Catalog owns the data, Core only exposes it.
- Future changes may move this functionality to a dedicated read-only store or a replication pipeline.


