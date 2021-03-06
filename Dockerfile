FROM golang:1.15

WORKDIR /go/src/github.com/Aiscom-LLC/meals-api
COPY . .

RUN go get -v

CMD ["meals-api"]

