# logdrift

A CLI tool that tails and diffs log streams across multiple services in real time.

---

## Installation

```bash
go install github.com/yourusername/logdrift@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logdrift.git
cd logdrift
go build -o logdrift .
```

---

## Usage

Tail and diff log streams from multiple services simultaneously:

```bash
logdrift --services auth-service,api-gateway,worker --follow
```

Compare log output between two services and highlight diverging lines:

```bash
logdrift diff --left auth-service --right api-gateway --since 10m
```

Watch all services defined in a config file:

```bash
logdrift --config ./logdrift.yaml
```

### Example Output

```
[auth-service]  2024-01-15T10:23:01Z INFO  user login successful uid=42
[api-gateway]   2024-01-15T10:23:01Z INFO  request routed path=/login
~ [worker]      2024-01-15T10:23:02Z WARN  queue depth high depth=150   ← drift detected
```

Lines prefixed with `~` indicate timing or content drift relative to other streams.

---

## Configuration

`logdrift.yaml` example:

```yaml
services:
  - name: auth-service
    source: journald
  - name: api-gateway
    source: file
    path: /var/log/api-gateway.log
  - name: worker
    source: docker
    container: worker-prod
```

---

## License

MIT © 2024 [yourusername](https://github.com/yourusername)