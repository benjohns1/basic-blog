FROM golang:1.13.0 AS builder
COPY ./services /services/
WORKDIR /services/authentication/cmd
RUN env GOOS=linux GOARCH=386 go build

FROM scratch
COPY --from=builder /services/authentication/cmd/cmd /
ENTRYPOINT ["/cmd"]