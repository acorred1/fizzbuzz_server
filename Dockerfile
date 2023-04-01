# syntax=docker/dockerfile:1

FROM golang:1.20.2-alpine3.17 AS build

WORKDIR /fizzbuzz_server

# Install dependencies first so this layer can be cached before our files
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o fizzbuzz_server


FROM scratch
WORKDIR /
COPY --from=build fizzbuzz_server/fizzbuzz_server fizzbuzz_server
EXPOSE 8080

CMD [ "./fizzbuzz_server" ]
