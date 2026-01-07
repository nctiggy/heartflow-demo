# HeartFlow Demo Service

A demo web service for Spectro Cloud Edge workshops that displays HeartFlow branding and shows whether it's running on systemd or Kubernetes.

## Features

- HeartFlow-branded UI with gradient effects
- Automatic runtime detection (systemd vs Kubernetes)
- Health check endpoints (`/health`, `/healthz`)
- Lightweight Go binary

## Deployment Options

### 1. Kubernetes (Helm)

```bash
helm install heartflow-demo ./helm/heartflow-demo
```

Access via NodePort: `http://<node-ip>:30080`

### 2. Systemd (CanvOS)

Add the binary and service file to your CanvOS Dockerfile. See `canvos-dockerfile-snippet.txt` for integration instructions.

### 3. Docker

```bash
docker run -p 8080:8080 nctiggy/heartflow-demo:main
```

## Building

```bash
# Build binary
go build -o heartflow-demo .

# Build container
docker build -t heartflow-demo .
```

## Environment Variables

- `PORT` - Server port (default: 8080)
- `POD_NAME` - Pod name (auto-detected in K8s)
- `POD_NAMESPACE` - Pod namespace (auto-detected in K8s)
