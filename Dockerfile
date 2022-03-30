FROM golang:1.18.0-alpine3.15 as builder
COPY go.mod go.sum /app/
WORKDIR /app
RUN go mod download
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/service /app

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /app/build/service /usr/bin/service
EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/service"]
