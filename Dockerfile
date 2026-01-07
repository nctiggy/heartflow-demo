# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY main.go ./
COPY templates/ ./templates/

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o heartflow-demo .

# Runtime stage - Ubuntu with systemd
FROM ubuntu:22.04

# Install systemd and clean up
RUN apt-get update && apt-get install -y --no-install-recommends \
    systemd \
    systemd-sysv \
    ca-certificates \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    # Remove unnecessary systemd services
    && rm -f /lib/systemd/system/multi-user.target.wants/* \
    && rm -f /etc/systemd/system/*.wants/* \
    && rm -f /lib/systemd/system/local-fs.target.wants/* \
    && rm -f /lib/systemd/system/sockets.target.wants/*udev* \
    && rm -f /lib/systemd/system/sockets.target.wants/*initctl* \
    && rm -f /lib/systemd/system/basic.target.wants/* \
    && rm -f /lib/systemd/system/anaconda.target.wants/*

# Install the heartflow-demo binary
COPY --from=builder /app/heartflow-demo /usr/local/bin/heartflow-demo
RUN chmod +x /usr/local/bin/heartflow-demo

# Install the systemd service file
COPY heartflow-demo.service /etc/systemd/system/heartflow-demo.service

# Enable the service
RUN systemctl enable heartflow-demo.service

EXPOSE 8080

# Use systemd as init
STOPSIGNAL SIGRTMIN+3
CMD ["/lib/systemd/systemd"]
