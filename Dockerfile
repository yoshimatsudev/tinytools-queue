From golang:alpine

WORKDIR /go/src/app

COPY . .

RUN go build -o main main.go



EXPOSE 9090

ENV HOST=0.0.0.0
ENV PORT=9090

ENTRYPOINT ["./main"]