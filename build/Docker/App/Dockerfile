FROM golang:1.17.3

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download

RUN apt-get update && apt-get install -y git vim
RUN go get -u github.com/cosmtrek/air

CMD ["air"]
