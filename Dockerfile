FROM golang:alpine

RUN mkdir -p /data 6& mkdir -p /app
WORKDIR /app

RUN apk add --update \
    build-base \
    git \
    curl \
    wget \
    zip \
    unzip

RUN go get github.com/mattn/go-sqlite3 \
    github.com/asaskevich/govalidator \
    github.com/gorilla/handlers \
    github.com/gorilla/mux

COPY . .

RUN go build main.go database.go

EXPOSE 1337

CMD ["./main"]
