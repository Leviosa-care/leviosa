# Vault Production Setup — Current State

This document used to be a forward-looking migration plan ("move Vault from dev mode to production mode"). That migration is **done** — both `playbooks/deploy.yml` and `playbooks/deploy-staging.yml` already deploy and manage Vault in production mode. This document now describes what's actually in place and how to operate it; for day-to-day commands see `VAULT.md`.

## Current State

- Vault runs with the `file` storage backend (`roles/app/templates/vault.hcl.j2`), not dev mode — no `VAULT_DEV_ROOT_TOKEN_ID`, no in-memory storage
- Each environment runs its own container: `leviosa_vault` (production), `leviosa_staging_vault` (staging)
- TLS is disabled inside the listener (`tls_disable = "true"`) because Vault is never exposed outside the Docker network — no host port is published; Caddy terminates TLS for everything that does leave the host
- `disable_mlock = true` is set in `vault.hcl.j2`, so the container does not need `privileged: true` (it still requests the `IPC_LOCK` capability defensively)
- Initialization, unsealing, secrets-engine setup, and key/pepper provisioning are fully automated by the deploy playbooks — there is no manual setup step on a normal deploy
- The root token and unseal keys are generated on first deploy and saved only on the server, at `{{ app_data_dir }}/vault/vault-keys.txt` (mode `0600`) — they are never stored in git, group_vars, or Ansible Vault

## What the Deploy Playbooks Do (Automated)

On `make deploy` / `make deploy-staging`:

1. Render `vault.hcl` and start the Vault container if it isn't already healthy
2. If not yet initialized: `vault operator init -key-shares=3 -key-threshold=2`, write the keys/token to `vault-keys.txt` on the server
3. Unseal with 2 of the 3 saved keys
4. Write the live root token into the app's `.env` (`VAULT_TOKEN=...`)
5. Enable the `transit` and `secret` (KV-v2) engines if not already enabled
6. Create the transit key (`encx_kek_alias`, default `leviosa-kek`) and the `encx` pepper at `secret/encx/{{ encx_pepper_alias | default('leviosa') }}/pepper` if missing
7. Verify the pepper is reachable from the app's Docker network before continuing the rest of the deploy

Re-running deploys is idempotent — already-completed steps are skipped based on `vault-keys.txt` and the container's reported seal/init status.

## Operating Vault

See `VAULT.md` for: checking status, viewing the unseal keys/root token, manually unsealing, inspecting the pepper/transit keys, and forcing a fresh re-initialization with `vault_force_wipe=true`.

## Gaps / Not Yet Done

These were aspirational items in the original migration plan that are **not** implemented — worth knowing if you're hardening this further:

- **AppRole authentication**: the app still authenticates to Vault with the root token (written into `.env`), not a scoped AppRole/policy. There is no `leviosa-production`/`leviosa-staging` policy file or path separation — both environments use the default `secret/` and `transit/` mounts on their own Vault instance (already separate since each environment runs its own container)
- **Auto-unseal**: unsealing after a restart is manual (or via re-running a deploy) — no AWS KMS/HSM auto-unseal is configured
- **Vault-specific backups**: there's no scheduled `vault operator raft snapshot` or equivalent for Vault's own data; only application-level backup (`rclone`/`gpg` roles, `playbooks/backup.yml`) exists, and that isn't wired to a bucket yet either (see `README.md`)
- **Audit logging**: not enabled

## Troubleshooting

### CAP_SETFCAP Error (Vault Restart Loop)
```
unable to set CAP_SETFCAP effective capability: Operation not permitted
```
Already handled by `disable_mlock = true` in `vault.hcl.j2`. If you see this, check that the rendered config on the server actually has that line — a stale `vault.hcl` from before this was added would need a redeploy to pick it up.

### Vault Sealed After Restart
```bash
ssh deploy@<server-ip>
cat /opt/leviosa/data/vault/vault-keys.txt   # or /opt/leviosa-staging/...
docker exec leviosa_vault vault operator unseal <UNSEAL_KEY_1>
docker exec leviosa_vault vault operator unseal <UNSEAL_KEY_2>
```
Or just re-run `make deploy` / `make deploy-staging` — the playbook unseals automatically.

### Backend Cannot Connect to Vault
1. Check Vault is running: `docker ps | grep vault`
2. Check Vault status: `docker exec leviosa_vault vault status`
3. Check environment variables in the backend container: `docker exec leviosa_backend env | grep VAULT`
4. Verify network connectivity: `docker exec leviosa_backend ping leviosa_vault` (or the staging network alias)

### Lost `vault-keys.txt`
There is no recovery path — re-initialize with `vault_force_wipe=true` and accept that anything encrypted via `encx` under the old transit key/pepper is unrecoverable.

## Additional Resources

- [Vault Production Hardening Guide](https://developer.hashicorp.com/vault/docs/operator/production)
- [Vault AppRole Auth Method](https://developer.hashicorp.com/vault/docs/auth/approle) — relevant if closing the "Gaps" above
- [Vault Best Practices](https://developer.hashicorp.com/vault/docs/operations/best-practices)
