FROM golang:1.19

WORKDIR /usr/src/app

COPY *.go ./
COPY go.mod .
COPY go.sum .

CMD ["go", "test"]
