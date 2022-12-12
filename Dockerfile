FROM golang:1.19.2 as builder
WORKDIR /src/
COPY cmd cmd
COPY service service
COPY public public
COPY vendor vendor
COPY go.mod .
COPY go.sum .
RUN go build -o /tmp/webapi ./cmd/webapi/

FROM debian:stable
COPY --from=builder /tmp/webui /bin/webapi





CMD ["/bin/webapi"]
