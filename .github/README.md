# GitHub Actions Workflows

This directory contains the CI/CD workflows for the Leviosa project, implementing GitHub Flow with environment promotion.

## Workflow Overview

```
feature/new-feature
    ↓ git push
    ↓ PR to main
    ↓ ci.yaml runs (validation only)
main branch
    ↓ merge triggers staging.yaml
staging environment (staging.leviosa.com)
    ↓ git tag v*.*.*
    ↓ production.yaml runs
production environment (leviosa.com)
```

## Workflows

### 1. CI Workflow (`ci.yaml`)

**Trigger:** Pull requests to `main` branch
**Purpose:** Validate code quality without deployment
**Status Required:** Yes (blocks merge if failed)

#### Jobs:
- **Frontend Checks**: Install deps, build, unit tests
- **Backend Checks**: Install deps, build, unit tests
- **Backend Integration Tests**: Full integration tests with testcontainers (PostgreSQL, Redis, RabbitMQ, S3)
- **Security Scan**: Go vulnerability scanning with govulncheck
- **Summary**: Aggregate all job results

**Failure Behavior:** PR cannot merge until all checks pass

---

### 2. Staging Deployment (`staging.yaml`)

**Trigger:** Push to `main` branch
**Purpose:** Deploy to password-protected staging environment
**Environment:** `staging.leviosa.com`

#### Jobs:
1. **Frontend Build**: Build and push frontend Docker image
2. **Backend Build**: Build, test, and push backend Docker image
3. **Security Scan**: Trivy scan for CRITICAL/HIGH vulnerabilities
4. **Deploy**: Deploy to staging environment (self-hosted runner)

**Deployment Flow:**
1. Build frontend and backend images in parallel
2. Scan images with Trivy
3. If no critical vulnerabilities, deploy to staging
4. Verify health checks pass

---

### 3. Production Deployment (`production.yaml`)

**Trigger:** Git tags matching `v*.*.*` (e.g., `v1.0.0`, `v2.1.3`)
**Purpose:** Deploy to production environment
**Environment:** `leviosa.com`

#### Jobs:
1. **Frontend Build**: Build and push frontend Docker image
2. **Backend Build**: Build, test, and push backend Docker image
3. **Security Scan**: Trivy scan for CRITICAL/HIGH vulnerabilities
4. **Deploy**: Deploy to production environment (self-hosted runner)

**Manual Trigger:** Includes `workflow_dispatch` for emergency manual deployments

**Deployment Flow:**
1. Build frontend and backend images in parallel
2. Scan images with Trivy
3. If no critical vulnerabilities, deploy to production
4. Verify health checks pass

---

## Custom Actions

### Workflow Actions

Located in `.github/actions/workflow-actions/`:

#### `back-build`
Builds, tests, and pushes backend Docker images.
- Setup Go environment
- Build Go binary
- Run unit tests
- Push to Docker Hub

#### `front-build`
Builds, tests, and pushes frontend Docker images.
- Setup pnpm and Node.js
- Install dependencies
- Build SvelteKit application
- Run unit tests
- Push to Docker Hub

#### `deploy-app`
Deploys application using docker-compose.
- Login to Docker Hub
- Pull latest images
- Remove existing containers
- Start services with docker-compose
- Verify health checks
- Cleanup old resources

### Utility Actions

Located in `.github/actions/utility-actions/`:

#### `push-docker-image`
Builds and pushes Docker images to registry.
- Login to Docker Hub
- Setup BuildKit cache
- Build multi-stage Docker image
- Push to registry with proper tagging

---

## Security Features

### Vulnerability Scanning

**Go Dependencies:** `govulncheck` in CI workflow
**Docker Images:** Trivy scanning before deployment

**Severity Threshold:** CRITICAL and HIGH vulnerabilities block deployment

### Image Signing

Images are tagged with:
- Environment: `staging-frontend`, `production-backend`
- Git SHA: Included as build argument for traceability

---

## Environment Variables & Secrets

### Required Secrets (GitHub Repository Settings):

**Docker Hub:**
- `DOCKERHUB_USERNAME`
- `DOCKERHUB_TOKEN`

**AWS:**
- `AWS_REGION`
- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`
- `BUCKETNAME`

**External Services:**
- `STRIPE_SECRET_KEY`
- `GMAIL_EMAIL`
- `GMAIL_PASSWORD`
- `GOOGLE_CLIENT_ID`
- `GOOGLE_CLIENT_SECRET`

**Infrastructure:**
- `REDIS_ADDR`
- `REDIS_PASSWORD`
- `REDIS_DB`
- `USER_ENCRYPTION_KEY`
- `LOGGING_SALT`

---

## Development Workflow

### Making Changes

1. **Create Feature Branch:**
   ```bash
   git checkout -b feature/my-feature
   ```

2. **Make Changes & Commit:**
   ```bash
   git add .
   git commit -m "feat: add new feature"
   ```

3. **Push & Create PR:**
   ```bash
   git push origin feature/my-feature
   # Create PR on GitHub to main branch
   ```

4. **CI Validation:**
   - CI workflow runs automatically
   - All checks must pass (tests, security, build)

5. **Merge to Main:**
   - Use "Rebase and merge" (preserves commit history)
   - Staging deployment triggers automatically

6. **Test on Staging:**
   - Access staging.leviosa.com
   - Verify functionality
   - Test with realistic data

7. **Promote to Production:**
   ```bash
   git tag -a v1.2.3 -m "Release v1.2.3: Description"
   git push origin v1.2.3
   ```

---

## Repository Settings

### Branch Protection (Recommended)

Protect `main` branch:
- ✅ Require status checks before merging
- ✅ Require branches to be up to date
- ✅ Include CI workflow as required check
- ❌ Automatically delete head branches (keep history)

### Merge Strategy

**Configured:** Rebase and merge only
**Disabled:** Squash merging, merge commits

**Why:** Preserves full commit history for portfolio showcase

---

## Troubleshooting

### CI Workflow Failures

**Integration Tests Failing:**
- Testcontainers require Docker
- Check Docker daemon is running in CI
- Verify sufficient resources for containers

**Security Scan Failures:**
- Update vulnerable dependencies
- Check Go modules: `go get -u ./...`
- Review Trivy scan output for specific CVEs

### Deployment Failures

**Health Check Timeouts:**
- Check application logs: `docker compose logs`
- Verify all services started: `docker compose ps`
- Ensure dependencies (Redis, RabbitMQ) are healthy

**Image Pull Failures:**
- Verify Docker Hub credentials
- Check image tags match environment
- Ensure images were pushed successfully

---

## Monitoring & Logs

**CI Workflow Logs:** Available in GitHub Actions tab
**Deployment Logs:** Self-hosted runner output
**Application Logs:** Grafana + Loki (configured in compose.yaml)

---

## Future Enhancements

Potential improvements not yet implemented:

- Code coverage reporting (Codecov)
- Performance benchmarking
- Automated dependency updates (Dependabot)
- Slack/Discord deployment notifications
- Blue-green deployment strategy
- Deployment rollback automation
- Infrastructure validation (Terraform)

---

## Contact & Support

For issues with CI/CD workflows, check:
1. GitHub Actions logs for specific error messages
2. Docker build logs for image issues
3. Application health check endpoints
4. Grafana dashboards for runtime issues
