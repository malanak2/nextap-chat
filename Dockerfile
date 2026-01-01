# jet -dsn=postgresql://program:Password123@localhost:5432/chatdb?sslmode=disable -path=./gen
# Copy sources with gen
# Regen swagger
# Run
FROM golang AS builder
WORKDIR /app
COPY data ./data
COPY docs ./docs
COPY domain ./domain
COPY gen ./gen
COPY handlers ./handlers
COPY migrations ./
COPY ports ./
COPY usecases ./
COPY go.mod .
COPY main.go .

RUN go mod tidy
RUN CGO_ENABLED=0 go build .
RUN ls

FROM alpine
WORKDIR /app
COPY --from=builder /app/nextap-chat ./
RUN ls
CMD ["/app/nextap-chat"]