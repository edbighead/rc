FROM golang:1.15-alpine AS build

RUN apk add --no-cache git

WORKDIR /tmp/app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./out/rc .

FROM alpine:3.9 

COPY --from=build /tmp/app/out/rc /usr/local/bin/rc
ENTRYPOINT ["rc"]