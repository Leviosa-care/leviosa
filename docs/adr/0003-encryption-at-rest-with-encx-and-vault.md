# ADR-0003: Field-Level Encryption at Rest with encx and HashiCorp Vault

**Status**: Accepted  
**Date**: 2025-07

---

## Context

The platform stores personally identifiable information (PII): client and partner notes, building addresses, and contact details. GDPR requires that this data be protected at rest. A breach of the database should not expose readable PII.

Encrypting the entire database or using transparent disk encryption protects against physical media theft but not against a compromised application user with direct DB access.

## Decision

Use field-level encryption via the internal `encx` library. Fields tagged `encx:"encrypt"` in domain structs are encrypted before being persisted and decrypted after retrieval — transparently to the application layer.

Encryption keys are **per-service** secrets stored in HashiCorp Vault at the path `secret/data/encx/{service}/pepper`. This means a compromise of one service's key does not expose another service's data.

Key management:
- Vault is the single source of truth for keys
- Playbooks use `vault kv put/get secret/encx/{service}/pepper`
- Integration tests use `testutils.SetupServiceVault()` to provision per-service test keys
- Staging Vault uses a separate alias from production (`VAULT_ADDR` must be set explicitly per environment)

## Consequences

**Positive:**
- PII is unreadable from raw DB access or a DB backup leak
- Per-service key isolation limits blast radius of a key compromise
- Encryption is transparent to application and domain code

**Negative:**
- Encrypted fields cannot be indexed or searched by the database
- Key rotation requires re-encrypting all affected rows
- Vault becomes a critical dependency: if Vault is unavailable, the application cannot start

**Constraint**: Do not store Vault credentials or pepper values in `.env` files committed to the repository. Production keys are injected via environment at deploy time.
