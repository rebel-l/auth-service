# Builder
FROM golang:latest
WORKDIR /usr/src/app
COPY . .
RUN go build

# Service
FROM ubuntu:latest
WORKDIR /usr/bin
COPY --from=0 /usr/src/app/auth-service .
EXPOSE 3000/tcp
ENTRYPOINT ["auth-service"]
