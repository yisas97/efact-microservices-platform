FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

RUN swag init -g cmd/api/main.go --output ./docs

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ms1-documents ./cmd/api

FROM alpine:3.19

RUN apk --no-cache add ca-certificates

RUN adduser -D appuser
USER appuser

WORKDIR /home/appuser

COPY --from=builder /app/ms1-documents .

EXPOSE 5000

CMD ["./ms1-documents"]