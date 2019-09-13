FROM golang:1.13.0 AS builder
COPY . /api/
WORKDIR /api/cmd
RUN env GOOS=linux GOARCH=386 go build -o api

FROM scratch
COPY --from=builder /api/cmd/api /
ENTRYPOINT ["/api"]