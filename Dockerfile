FROM golang:1.24.1-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o joeburgess

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/joeburgess .

COPY --from=builder /app/content/ /app/content/
COPY static/ /app/static/
COPY templates/ /app/templates/

RUN mkdir -p /app/data

EXPOSE 8081

CMD ["./joeburgess"]
