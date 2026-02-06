FROM golang:alpine

WORKDIR /subscriptions-api

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/run

CMD ["./main"]