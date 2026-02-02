# Generate db stuff: jet -dsn=postgresql://program:Password123@localhost:5432/chatdb?sslmode=disable -path=./gen
# Copy sources with gen
# Regen swagger
# Run
FROM golang AS builder
WORKDIR /app
RUN go install github.com/swaggo/swag/cmd/swag@latest
#COPY docs ./docs
COPY domain ./domain
COPY gen ./gen
COPY handlers ./handlers
COPY migrations ./migrations
COPY ports ./ports
COPY usecases ./usecases
COPY go.mod .
COPY main.go .

RUN go mod tidy
RUN swag init
RUN CGO_ENABLED=0 go build .

FROM alpine
WORKDIR /app
COPY --from=builder /app/nextap-chat ./
COPY --from=builder /app/docs ./docs
CMD ["/app/nextap-chat"]
