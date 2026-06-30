# Monitoring Stack — Operations Guide

The monitoring stack runs four containers on the VPS, all bound to `127.0.0.1` (not exposed to the internet):

| Container | Port | Purpose |
|---|---|---|
| `leviosa_cadvisor` | 8081 | Per-container CPU/memory/network metrics |
| `leviosa_node_exporter` | 9100 | Host-level OS metrics (CPU, RAM, disk, swap) |
| `leviosa_prometheus` | 9090 | Metrics storage and query engine |
| `leviosa_grafana` | 3001 | Dashboards and alerting |

Deployed by `infra/ansible/roles/monitoring/`, enabled via `monitoring_enabled: true` in `group_vars/all.yml`.

---

## Accessing the UIs

All ports are localhost-only. Use an SSH tunnel to reach them from your machine:

```bash
# Open Grafana (recommended entry point)
ssh -L 3001:127.0.0.1:3001 root@46.224.117.126 -N

# Open Prometheus (raw queries / debugging)
ssh -L 9090:127.0.0.1:9090 root@46.224.117.126 -N
```

Then visit `http://localhost:3001` (Grafana) or `http://localhost:9090` (Prometheus) in your browser.

For Grafana, the default login is `admin` / the password set in `grafana_admin_password` (see [Changing the Grafana password](#changing-the-grafana-password) below).

---

## First-Time Grafana Setup

After the first deploy:

1. Open Grafana via the SSH tunnel above
2. Log in as `admin`
3. Add Prometheus as a data source:
   - Go to **Connections → Data sources → Add new**
   - Type: **Prometheus**
   - URL: `http://prometheus:9090` (containers share the `monitoring` network)
   - Click **Save & test**
4. Import community dashboards:
   - **Node Exporter Full** — dashboard ID `1860` — host CPU, RAM, disk, swap
   - **cAdvisor** — dashboard ID `14282` — per-container resource usage
   - Go to **Dashboards → Import**, enter the ID, select your Prometheus data source

---

## What Each Metric Source Covers

**cAdvisor (container metrics)**
- Memory usage vs. limit per container — the key metric after the OOM incident
- CPU throttling — containers hitting their CPU cap
- Network I/O per container

**Node Exporter (host metrics)**
- Total RAM and swap usage across the host
- Disk usage per mount (catch disk-full before it kills the VPS)
- Load average and CPU saturation

**Prometheus (queries)**
- Stores 15 days of history by default (set via `--storage.tsdb.retention.time`)
- Access at `http://localhost:9090` — useful for ad-hoc PromQL queries

---

## Key Metrics to Watch

After setting up dashboards, keep an eye on:

```promql
# Container memory usage as % of its limit — alert if > 80%
container_memory_usage_bytes / container_spec_memory_limit_bytes * 100

# Host-level memory available (alert if < 300MB)
node_memory_MemAvailable_bytes

# Swap usage — non-zero means you're close to the edge
node_memory_SwapFree_bytes - node_memory_SwapTotal_bytes

# Disk usage % — alert if > 85%
(node_filesystem_size_bytes - node_filesystem_free_bytes) / node_filesystem_size_bytes * 100
```

---

## Setting Up Alerts

In Grafana: **Alerting → Alert rules → New alert rule**

Recommended rules:

| Alert | Condition | Severity |
|---|---|---|
| High memory pressure | `node_memory_MemAvailable_bytes < 300e6` | Critical |
| Container near limit | `container_memory_usage_bytes / container_spec_memory_limit_bytes > 0.85` | Warning |
| Disk filling up | filesystem used > 85% | Warning |
| Swap in use | `node_memory_SwapFree_bytes < node_memory_SwapTotal_bytes` | Warning |

For notifications, configure a **Contact point** (Grafana → Alerting → Contact points) — email via SMTP or a webhook to Slack/Telegram.

---

## Operational Tasks

### Check containers are running
```bash
ssh root@46.224.117.126
docker ps | grep leviosa_
# or via systemd
systemctl status leviosa-monitoring
```

### Restart the monitoring stack
```bash
ssh root@46.224.117.126
systemctl restart leviosa-monitoring
```

### View Prometheus logs
```bash
ssh root@46.224.117.126
docker logs leviosa_prometheus --tail 50
```

### Reload Prometheus config without restart
```bash
# After editing prometheus.yml on the server
curl -X POST http://127.0.0.1:9090/-/reload
```

### Changing the Grafana password

Option 1 — set it before deploy in `secrets.yml` (recommended):
```yaml
grafana_admin_password: "your-strong-password"
```
Then re-run `make ansible-deploy`.

Option 2 — change it in the UI after login:
- Profile → Change password

---

## Re-deploying / Updating

The monitoring stack is managed as a systemd service (`leviosa-monitoring.service`) wrapping a separate Docker Compose project at `~/leviosa/monitoring/docker-compose.yml` on the server.

Any change to `group_vars/all.yml`, the compose template, or `prometheus.yml.j2` is picked up by:

```bash
make ansible-deploy
```

Grafana dashboards and Prometheus data survive re-deploys because they use named Docker volumes (`prometheus_data`, `grafana_data`).
