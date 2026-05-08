# syntax=docker/dockerfile:1.7

ARG GO_VERSION=1.26.3
ARG TRIVY_VERSION=0.69.3
ARG SYFT_VERSION=1.42.4
ARG GRYPE_VERSION=0.111.0
ARG DUCKDB_VERSION=1.4.2

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS tools
ARG TARGETARCH
ARG TRIVY_VERSION
ARG SYFT_VERSION
ARG GRYPE_VERSION
ARG DUCKDB_VERSION
RUN apk add --no-cache ca-certificates curl tar unzip
WORKDIR /tools
RUN test "$TARGETARCH" = "amd64"
RUN curl -fsSL "https://github.com/aquasecurity/trivy/releases/download/v${TRIVY_VERSION}/trivy_${TRIVY_VERSION}_Linux-64bit.tar.gz" \
  | tar -xz trivy \
  && curl -fsSL "https://github.com/anchore/syft/releases/download/v${SYFT_VERSION}/syft_${SYFT_VERSION}_linux_amd64.tar.gz" \
  | tar -xz syft \
  && curl -fsSL "https://github.com/anchore/grype/releases/download/v${GRYPE_VERSION}/grype_${GRYPE_VERSION}_linux_amd64.tar.gz" \
  | tar -xz grype \
  && curl -fsSLo duckdb.zip "https://github.com/duckdb/duckdb/releases/download/v${DUCKDB_VERSION}/duckdb_cli-linux-amd64.zip" \
  && unzip duckdb.zip duckdb \
  && chmod +x trivy syft grype duckdb
RUN cp /bin/busybox /tools/busybox

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS builder
ARG VERSION=dev
ARG COMMIT=none
ARG DATE=unknown
WORKDIR /src
RUN apk add --no-cache ca-certificates git
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath \
  -ldflags="-s -w -X github.com/baditaflorin/audit-in-a-box/pkg/version.Version=${VERSION} -X github.com/baditaflorin/audit-in-a-box/pkg/version.Commit=${COMMIT} -X github.com/baditaflorin/audit-in-a-box/pkg/version.Date=${DATE}" \
  -o /out/audit-server ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
LABEL org.opencontainers.image.source="https://github.com/baditaflorin/audit-in-a-box"
LABEL org.opencontainers.image.licenses="MIT"
WORKDIR /app
COPY --from=builder /out/audit-server /app/audit-server
COPY --from=tools /tools/trivy /usr/local/bin/trivy
COPY --from=tools /tools/syft /usr/local/bin/syft
COPY --from=tools /tools/grype /usr/local/bin/grype
COPY --from=tools /tools/duckdb /usr/local/bin/duckdb
COPY --from=tools /tools/busybox /busybox
ENV SERVER_ADDR=:8080
ENV WORK_DIR=/tmp/audit-work
ENV TRIVY_CACHE_DIR=/tmp/trivy
ENV GRYPE_DB_CACHE_DIR=/tmp/grype
ENV HOME=/tmp
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 CMD ["/busybox", "wget", "-qO-", "http://127.0.0.1:8080/healthz"]
USER nonroot:nonroot
ENTRYPOINT ["/app/audit-server"]
