FROM golang:1.12.7 as builder

WORKDIR /build/middleman

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

COPY go.mod .
COPY go.sum .
RUN go get

COPY main.go .
RUN go build -ldflags="-s -w" -o middleman main.go
RUN chmod +x middleman

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/middleman/middleman /usr/bin/middleman

ENTRYPOINT ["middleman"]
