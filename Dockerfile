# Builder
FROM golang:latest as builder
WORKDIR /usr/src/app
COPY . .
RUN go build

# Service
FROM ubuntu:latest
RUN apt-get update \
    && apt-get install -y ca-certificates
WORKDIR /usr/bin/app
ENV PATH="/usr/bin/app:${PATH}"
COPY --from=builder /usr/src/app/auth-service .
COPY --from=builder /usr/src/app/scripts ./scripts
EXPOSE 3000/tcp
ENTRYPOINT ["auth-service"]
