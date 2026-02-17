FROM golang:alpine AS builder

RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o miniflux-mcp .

FROM scratch AS runtime

WORKDIR /app

COPY --from=builder /app/miniflux-mcp /app/miniflux-mcp
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/app/miniflux-mcp"]
