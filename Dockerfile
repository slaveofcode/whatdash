FROM golang:alpine as builder
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

# Create appuser
RUN adduser -D -g '' appuser

COPY . $GOPATH/src/whatdash/
WORKDIR $GOPATH/src/whatdash/

RUN go get -d -v

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/whatdash

FROM scratch as app-whatdash
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/whatdash /go/bin/whatdash
COPY --from=builder /go/src/whatdash/static /static

USER root

EXPOSE 8081

ENTRYPOINT [ "/go/bin/whatdash" ]
