FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/api
RUN go build -o consumer ./cmd/consumer
RUN go build -o projection ./cmd/projections

EXPOSE 8080

CMD ["./main"]