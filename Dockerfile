#Production Dockerfile

FROM golang:1.25.1-alpine

WORKDIR /home/portfolio


# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/server/main.go

CMD ["./main"]

EXPOSE 3306


