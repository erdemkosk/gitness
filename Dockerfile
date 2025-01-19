FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.* .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o gitness

FROM alpine:latest
COPY --from=builder /app/gitness /gitness
RUN apk --no-cache add ca-certificates

ENTRYPOINT ["/gitness"] 