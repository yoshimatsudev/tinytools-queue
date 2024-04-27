From golang:alpine

WORKDIR /go/src/app

COPY . .

RUN go build -o main main.go



EXPOSE 9090

ENV PORT=9090

ENTRYPOINT ["./main"]