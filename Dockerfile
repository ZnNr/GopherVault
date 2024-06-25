FROM golang:1.22.3-alpine3.18 as builder

COPY . /app
WORKDIR /app/
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/GopherVault github.com/ZnNr/GopherVault

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /app/build/GopherVault /usr/bin/GopherVault
EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/GopherVault", "run"]