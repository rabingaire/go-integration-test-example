FROM golang:alpine

RUN apk add --no-cache make curl gcc libc-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD make integration-tests
