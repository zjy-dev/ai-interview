FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git make

WORKDIR /src

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .

ARG VERSION=dev
RUN CGO_ENABLED=0 go build -ldflags "-s -w -X main.Version=${VERSION}" -o /app/server ./cmd/server/

# ---
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/server .
COPY backend/configs/ ./configs/

EXPOSE 8000

ENTRYPOINT ["./server"]
CMD ["-conf", "configs/"]
