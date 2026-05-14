# driftwatch

Lightweight daemon that detects config drift between running containers and their source manifests.

---

## Installation

```bash
go install github.com/yourusername/driftwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/driftwatch.git && cd driftwatch && make build
```

---

## Usage

Point driftwatch at your manifests directory and let it run alongside your containers:

```bash
driftwatch --manifests ./deploy/manifests --interval 30s
```

**Example output:**

```
[2024-01-15 08:42:11] INFO  Watching 12 manifests...
[2024-01-15 08:43:22] WARN  Drift detected: container "api-server"
                            expected image: myapp:v1.4.2
                            running image:  myapp:v1.4.1
[2024-01-15 08:43:22] INFO  No drift detected in remaining 11 containers.
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--manifests` | `./manifests` | Path to source manifest files |
| `--interval` | `60s` | How often to poll for drift |
| `--output` | `text` | Output format (`text`, `json`) |
| `--alert-webhook` | — | Webhook URL to notify on drift |

---

## Configuration

Optionally create a `driftwatch.yaml` in your working directory:

```yaml
manifests: ./deploy/manifests
interval: 30s
output: json
alert_webhook: https://hooks.example.com/notify
```

---

## License

MIT © 2024 [yourusername](https://github.com/yourusername)