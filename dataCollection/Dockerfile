FROM golang:latest

RUN go get github.com/Shopify/sarama

WORKDIR /app

COPY . .

CMD ["go", "run", "main.go"]
