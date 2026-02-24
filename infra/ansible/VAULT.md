# Ansible Vault - Secure Variable Management

Ansible Vault encrypts sensitive variables so they can be safely stored in git.

## Quick Start

### 1. Create Your Vault File

```bash
# Copy the example file
cp group_vars/secrets.vault.example.yml group_vars/secrets.yml

# Edit with your actual values
nano group_vars/secrets.yml
```

### 2. Encrypt the Vault File

```bash
# Encrypt with a password (you'll be prompted)
ansible-vault encrypt group_vars/secrets.yml

# Or use a password file
echo "your_vault_password" > .vault_pass
chmod 600 .vault_pass
ansible-vault encrypt --vault-password-file .vault_pass group_vars/secrets.yml
```

### 3. Run Playbooks with Vault

```bash
# You'll be prompted for the vault password
ansible-playbook playbook.yml

# Or use a password file
ansible-playbook playbook.yml --vault-password-file .vault_pass
```

## Common Vault Commands

```bash
# View encrypted file (decrypts to stdout)
ansible-vault view group_vars/secrets.yml

# Edit encrypted file (decrypts, opens editor, re-encrypts on save)
ansible-vault edit group_vars/secrets.yml

# Change vault password
ansible-vault rekey group_vars/secrets.yml

# Decrypt file (removes encryption - be careful!)
ansible-vault decrypt group_vars/secrets.yml

# Encrypt file
ansible-vault encrypt group_vars/secrets.yml

# Check if file is encrypted
ansible-vault view group_vars/secrets.yml --vault-password-file /dev/null
```

## Password Management

### Option 1: Prompt (Default)

```bash
ansible-playbook playbook.yml --ask-vault-pass
```

### Option 2: Password File

Create `.vault_pass` (add to `.gitignore`!):

```bash
echo "your_vault_password" > .vault_pass
chmod 600 .vault_pass
```

Then:

```bash
ansible-playbook playbook.yml --vault-password-file .vault_pass
```

### Option 3: Environment Variable

```bash
export ANSIBLE_VAULT_PASSWORD_FILE=.vault_pass
ansible-playbook playbook.yml
```

### Option 4: Script (Dynamic Password)

Create `vault_pass.sh`:

```bash
#!/bin/bash
# Use a password manager, TPM, or other secure source
pass show ansible/vault_password
```

Make executable:

```bash
chmod +x vault_pass.sh
```

Then:

```bash
ansible-playbook playbook.yml --vault-password-file vault_pass.sh
```

## Security Best Practices

1. **Never commit `.vault_pass` to git** (it's in `.gitignore`)
2. **Use strong, unique passwords** (use a password manager)
3. **Rotate vault passwords periodically** (`ansible-vault rekey`)
4. **Limit vault access** (file permissions: `600`)
5. **Backup vault password securely** (password manager, not plain text)
6. **Use different passwords** for different environments

## Project Setup

### 1. Add to `.gitignore`

The `.gitignore` already includes:
```
.vault_pass
vault_pass.sh
```

### 2. Update `ansible.cfg`

The `ansible.cfg` file is configured with:

```ini
[defaults]
vault_password_file = .vault_pass  # Optional: uncomment to use
```

### 3. Store Password Securely

Recommended: Use a password manager like `pass`, `1password`, or `keepassxc`:

```bash
# Using pass (Linux/Mac)
pass insert -m ansible/leviosa/vault_password
chmod +x vault_pass.sh  # make script executable
```

## Troubleshooting

### "Decryption failed"

Wrong password or file corrupted. Verify:
```bash
ansible-vault view group_vars/secrets.yml
```

### "Vault password file not found"

Check file path and permissions:
```bash
ls -la .vault_pass
cat .vault_pass
```

### "Vault password mismatch"

You encrypted with one password, trying to decrypt with another.

## Migration Guide

### Moving from `.env` to Vault

1. Create `group_vars/secrets.yml` with your `.env` values
2. Encrypt the vault file
3. Remove sensitive values from `.env` (keep non-sensitive ones)
4. Update playbooks to reference vault variables

### Example Vault Variables

```yaml
# In group_vars/secrets.yml (encrypted)
vault_db_password: "secure_password"
vault_aws_secret_access_key: "AKIA..."

# In playbooks, reference as:
db_password: "{{ vault_db_password }}"
```

## Multiple Vaults (Advanced)

For multiple environments:

```bash
group_vars/
  secrets.yml              # Shared secrets
  manager_secrets.yml      # Manager-specific secrets
  worker_secrets.yml       # Worker-specific secrets
```

Use with inventory:

```yaml
# playbook.yml
vars_files:
  - group_vars/secrets.yml

# Different playbooks can use different vault files
```

## Resources

- [Ansible Vault Documentation](https://docs.ansible.com/ansible/latest/vault_guide/index.html)
- [Best Practices for Secrets Management](https://docs.ansible.com/ansible/latest/user_guide/vault.html#best-practices)
