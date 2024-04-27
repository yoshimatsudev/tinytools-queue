From golang:alpine

WORKDIR /go/src/app

COPY . .

RUN go build -o main main.go

CMD ["./main"]