FROM golang:latest AS builder

WORKDIR /app

COPY go.sum go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/web


FROM alpine:3.21.3

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup 

COPY --from=builder /app/server .
COPY --from=builder /app/tls ./tls
COPY --from=builder /app/ui ./ui

RUN chown -R appuser:appgroup /app 

USER appuser 

EXPOSE 4000

CMD ["./server"]