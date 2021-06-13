FROM golang:1.16

WORKDIR /go/src/github.com/VolticFroogo/htmlrender-api
COPY . .
RUN go build -o main .

CMD ["./main"]