FROM golang:1.21

WORKDIR /app
COPY . .

RUN go mod init client && go mod tidy && go build -o client .

CMD ["./client"]
