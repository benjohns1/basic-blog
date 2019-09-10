FROM golang:1.13.0 AS builder
COPY . /app/
WORKDIR /app
RUN env GOOS=linux GOARCH=386 go build

FROM scratch
COPY --from=builder /app/basic-blog /
COPY --from=builder /app/index.html /
ENTRYPOINT ["/basic-blog"]