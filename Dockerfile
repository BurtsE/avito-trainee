FROM golang:1.22-alpine AS builder

WORKDIR /src

ADD go.mod .
ADD go.sum .
RUN go mod download
RUN pwd


COPY cmd cmd
COPY internal internal

RUN GOOS=linux go build -o app ./cmd

FROM golang:1.22-alpine
WORKDIR /root/

COPY configs configs

ENV HOUSE_DB_USER=${HOUSE_DB_USER}
ENV HOUSE_DB_PASSWORD=${HOUSE_DB_PASSWORD}
ENV USER_DB_USER=${HOUSE_DB_USER}
ENV USER_DB_PASSWORD=${HOUSE_DB_PASSWORD}
COPY --from=builder /build/app .


CMD ["./app"]