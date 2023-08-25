FROM golang:alpine

WORKDIR /go-app

COPY . .

ENV GOPROXY=https://goproxy.cn,direct

RUN go mod download
RUN go build -o main .

CMD ["./main"]