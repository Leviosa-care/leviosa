# HashiCorp Vault — Operations Guide

This project runs a [HashiCorp Vault](https://developer.hashicorp.com/vault) container alongside the app (one per environment) to hold encryption material used by the backend's `encx` field-encryption library — a transit key for encrypt/decrypt and a KV-v2 secret (the "pepper") used in key derivation. It is **not** related to Ansible Vault (the `ansible-vault encrypt`/`--ask-vault-pass` feature) — this project does not use that; see [Secrets files](#secrets-files-not-ansible-vault-encrypted) below for how deployment secrets are actually handled.

## How It's Deployed

Vault runs in production mode (file storage backend, no dev-mode flags) via the template at `roles/app/templates/vault.hcl.j2`, rendered into `docker-compose.yml.j2`. Each environment (production, staging) runs its own Vault container:

- Production: `leviosa_vault`, network alias `vault`
- Staging: `leviosa_staging_vault`, network alias `<app_name>-vault`, reachable from the shared production network

Both `playbooks/deploy.yml` and `playbooks/deploy-staging.yml` contain the full Vault lifecycle automation — there is nothing to do by hand on a normal deploy:

1. Create `{{ app_data_dir }}/vault/{config,data}` with correct ownership, render `vault.hcl`
2. Start (or recreate, if unhealthy) the Vault container
3. If not yet initialized: `vault operator init -key-shares=3 -key-threshold=2`, save unseal keys + root token to `{{ app_data_dir }}/vault/vault-keys.txt` (mode `0600`, owned by the deploy user) on the server
4. Unseal using 2 of the 3 saved keys
5. Write the live root token into the app's `.env` (`VAULT_TOKEN=...`) via `lineinfile`
6. Enable the `transit` and `secret` (KV-v2) secrets engines (idempotent, `|| true`)
7. Create the transit keys (`encx_kek_alias`, default `leviosa-kek`) and generate/store the pepper at `secret/encx/{{ encx_pepper_alias | default('leviosa') }}/pepper` if it doesn't already exist
8. Verify the pepper is reachable from the app's network before continuing

Re-running a deploy is safe — steps that have already succeeded are skipped based on `vault-keys.txt` existing and the container's reported seal status.

## Manual Operations

### Check status
```bash
ssh deploy@<server-ip>
docker exec leviosa_vault vault status        # or leviosa_staging_vault
```

### View the unseal keys / root token
The only copy lives on the server itself — there is no copy in git or in group_vars:
```bash
ssh deploy@<server-ip>
cat /opt/leviosa/data/vault/vault-keys.txt          # production
cat /opt/leviosa-staging/data/vault/vault-keys.txt  # staging
```

### Manually unseal (if a deploy left it sealed)
```bash
docker exec leviosa_vault vault operator unseal <UNSEAL_KEY_1>
docker exec leviosa_vault vault operator unseal <UNSEAL_KEY_2>
```

### Inspect the pepper / transit keys
```bash
docker exec -e VAULT_TOKEN=<root-token> leviosa_vault vault kv get secret/encx/leviosa/pepper
docker exec -e VAULT_TOKEN=<root-token> leviosa_vault vault list transit/keys
```

### Force a fresh Vault (re-initialize, losing existing keys)
Pass `vault_force_wipe=true` on a deploy — this wipes `{{ app_data_dir }}/vault/data` and lets the playbook re-initialize from scratch. **Only do this if you accept losing access to anything encrypted with the old transit keys/pepper** (this includes any encrypted columns in the database — there is no migration path back).
```bash
ansible-playbook playbooks/deploy.yml -e "ansible_host=<server-ip>" -e "vault_force_wipe=true"
```

## Secrets Files (not Ansible-Vault-encrypted)

Deployment secrets (DB password, AWS keys, Stripe keys, etc.) live in `group_vars/leviosa_staging.yml` and `group_vars/leviosa_production.yml` — plain YAML, gitignored, never encrypted with `ansible-vault`. Copy from the committed `.example` files and fill in real values:

```bash
cp group_vars/leviosa_staging.example.yml group_vars/leviosa_staging.yml
cp group_vars/leviosa_production.example.yml group_vars/leviosa_production.yml
```

`make update-staging-vault` / `make update-production-vault` (run from `infra/terraform`) refresh the AWS credential fields in these files from the latest Terraform output — despite the Makefile target name, this is unrelated to HashiCorp or Ansible Vault; see `infra/scripts/update-staging-vault.sh`.

## Troubleshooting

### CAP_SETFCAP error / restart loop
```
unable to set CAP_SETFCAP effective capability: Operation not permitted
```
Already handled — `vault.hcl.j2` sets `disable_mlock = true`, which is why the container doesn't need `privileged: true` or `IPC_LOCK` for mlock specifically (it still requests `IPC_LOCK` defensively in the compose file).

### Sealed after a host reboot
Vault's file storage persists, but a fresh container always starts sealed. Either redeploy (`make deploy` / `make deploy-staging`, which unseals automatically) or unseal manually as above.

### Backend can't reach Vault
```bash
docker ps | grep vault
docker exec leviosa_vault vault status
docker exec leviosa_backend env | grep VAULT
docker exec leviosa_backend ping leviosa_vault   # or the staging alias
```

### Lost `vault-keys.txt`
There is no recovery — the unseal keys and root token only ever existed in that file and in Vault's own internal state. You must wipe and re-initialize (`vault_force_wipe=true`), which destroys access to anything previously encrypted via `encx`.

## Security Notes

- Vault is not exposed outside the Docker network — no host port is published, only `expose:` for inter-container access
- TLS is disabled inside Vault's listener (`tls_disable = "true"`) because traffic never leaves the host's internal Docker network; Caddy terminates TLS for everything that does
- The root token is used directly by automation — this project does not yet use AppRole or scoped policies. Treat `vault-keys.txt` on the server with the same care as the database credentials

## Resources

- [Vault Production Hardening Guide](https://developer.hashicorp.com/vault/docs/operator/production)
- [Vault CLI Reference](https://developer.hashicorp.com/vault/docs/commands)
