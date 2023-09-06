FROM golang:alpine

WORKDIR /go-app/blog-api

COPY . .

ENV GOPROXY=https://goproxy.cn,direct
ENV GIN_MODE=release

RUN go mod download
RUN go build -o main .

CMD ["./main"]