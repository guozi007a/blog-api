FROM golang:alpine

WORKDIR /go-app

COPY . .

RUN go build -o main .

CMD ["./main"]