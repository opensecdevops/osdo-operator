FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o manager .

FROM gcr.io/distroless/static:nonroot
LABEL org.opencontainers.image.title="OSDO Operator" \
      org.opencontainers.image.description="Kubernetes Operator para seguridad nativa de cluster OSDO" \
      org.opencontainers.image.source="https://github.com/opensecdevops/osdo-operator" \
      org.opencontainers.image.licenses="Apache-2.0"
COPY --from=builder /app/manager .
USER 65532:65532
ENTRYPOINT ["/manager"]
