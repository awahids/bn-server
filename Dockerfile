FROM golang:1.24-alpine AS builder

WORKDIR /src

RUN apk add --no-cache ca-certificates tzdata

ARG TARGETOS=linux
ARG TARGETARCH=amd64

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN set -eux; \
    fixdir() { from="$1"; to="$2"; if [ -d "$from" ] && [ ! -d "$to" ]; then mv "$from" "$to"; fi; }; \
    fixdir internal/delivery/data/request/authReq internal/delivery/data/request/authreq; \
    fixdir internal/delivery/data/response/authRes internal/delivery/data/response/authres; \
    fixdir internal/delivery/handlers/authHandler internal/delivery/handlers/authhandler; \
    fixdir internal/delivery/router/authRouter internal/delivery/router/authrouter; \
    fixdir internal/delivery/router/publicRouter internal/delivery/router/publicrouter; \
    fixdir internal/domain/repositories/appRepo internal/domain/repositories/apprepo; \
    fixdir internal/domain/repositories/authRepo internal/domain/repositories/authrepo; \
    fixdir internal/domain/repositories/repoInterface internal/domain/repositories/repointerface; \
    fixdir internal/domain/services/appService internal/domain/services/appservice; \
    fixdir internal/domain/services/authService internal/domain/services/authservice; \
    fixdir internal/domain/services/serviceInterface internal/domain/services/serviceinterface

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o /out/server ./cmd

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o /out/migrate ./cmd/migrate

FROM alpine:3.21

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S -G app app

COPY --from=builder --chown=app:app /out/server /app/bin/server
COPY --from=builder --chown=app:app /out/migrate /app/bin/migrate
COPY --from=builder --chown=app:app /src/internal/infrastructure/database/migrations /app/internal/infrastructure/database/migrations

USER app

EXPOSE 8080

CMD ["/app/bin/server"]
