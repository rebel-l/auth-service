# Builder
FROM golang:latest
WORKDIR /usr/src/app
COPY . .
RUN go build

# Service
FROM ubuntu:latest
RUN apt-get update \
    && apt-get install -y ca-certificates
WORKDIR /usr/bin
COPY --from=0 /usr/src/app/auth-service .
EXPOSE 3000/tcp
ENTRYPOINT ["auth-service"]
