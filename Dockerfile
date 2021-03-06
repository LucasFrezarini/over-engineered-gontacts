# build stage
FROM golang:1.14-alpine as builder

ENV GO111MODULE on

RUN apk update; \
    apk add --no-cache git; 

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

# final stage
FROM scratch

COPY --from=builder /app/go-contacts /app/

EXPOSE 8080

ENTRYPOINT [ "/app/go-contacts" ]
