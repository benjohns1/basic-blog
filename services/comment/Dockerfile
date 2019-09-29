FROM golang:1.13.0 AS builder
COPY ./services /services/
WORKDIR /services/comment/cmd
RUN env GOOS=linux GOARCH=386 go build

FROM scratch
COPY --from=builder /services/comment/cmd/cmd /
ENTRYPOINT ["/cmd"]