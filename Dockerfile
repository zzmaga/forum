FROM golang:latest

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o forum ./cmd/main.go

EXPOSE 8080
CMD ["./forum"]
