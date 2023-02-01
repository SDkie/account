FROM golang:1.19

WORKDIR /usr/src/app

COPY *.go /usr/src/app/ 
COPY go.mod .
COPY go.sum .

CMD ["go", "test"]
