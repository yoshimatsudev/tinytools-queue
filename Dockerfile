From golang:alpine

WORKDIR /go/src/app

COPY . .

RUN go build -o main main.go

EXPOSE 8080

ENV PORT="8080"

CMD ["./main"]