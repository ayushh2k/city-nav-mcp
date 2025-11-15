FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o ./bin/mcp-server ./cmd/mcp-server/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/mcp-server .

RUN apk --no-cache add ca-certificates

EXPOSE 8000

CMD ["./mcp-server"]