FROM golang:1.21.1

WORKDIR /go/src/app

COPY . .

RUN go build -o main main.go
EXPOSE 9000

CMD ["./main"]