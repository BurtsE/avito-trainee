FROM golang:1.22-alpine AS builder

WORKDIR /build

ADD /src/go.mod .
ADD /src/go.sum .
RUN go mod download
RUN pwd


COPY /src/cmd cmd
COPY /src/internal internal

RUN GOOS=linux go build -o app ./cmd

FROM golang:1.22-alpine
WORKDIR /root/

COPY /src/configs configs

ENV POSTGRES_USERNAME=${POSTGRES_USERNAME}
ENV POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
ENV POSTGRES_HOST=${POSTGRES_HOST}
ENV POSTGRES_PORT=${POSTGRES_PORT}
ENV POSTGRES_DATABASE=${POSTGRES_DATABASE}
ENV SERVER_ADDRESS=${SERVER_ADDRESS}

COPY --from=builder /build/app .


CMD ["./app"]