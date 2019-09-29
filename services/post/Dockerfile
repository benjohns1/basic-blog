FROM golang:1.13.0 AS builder
COPY ./services /services/
WORKDIR /services/post/cmd
RUN env GOOS=linux GOARCH=386 go build

FROM scratch
COPY --from=builder /services/post/cmd/cmd /
ENTRYPOINT ["/cmd"]