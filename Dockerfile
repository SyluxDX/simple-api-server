FROM golang:alpine as builder
WORKDIR /src
COPY . .
RUN go build -o api-server server.go
RUN rm -rf .dockerignore server.govim

## can't use scratch due to net dynamic libraries
FROM alpine
WORKDIR /src
COPY --from=builder /src .
ENTRYPOINT ["/src/api-server"]