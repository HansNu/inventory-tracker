# ---- build stage ----
FROM golang:1.26.2-alpine AS builder

WORKDIR /src

# Dependencies first. This layer stays cached until go.mod/go.sum change,
# so editing a handler doesn't re-download every module.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO off => a fully static binary, which is what lets the final stage be alpine.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/server .

# ---- run stage ----
FROM alpine:3.20

RUN apk add --no-cache ca-certificates && adduser -D -u 10001 app

COPY --from=builder /out/server /app/server

USER app
EXPOSE 8080
ENTRYPOINT ["/app/server"]