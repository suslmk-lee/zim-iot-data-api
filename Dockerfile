FROM golang:1.22-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod go.sum main.go ./

COPY config ./config
COPY database ./database
COPY handlers ./handlers
COPY models ./models
COPY utils ./utils
RUN go mod download
RUN go mod tidy
RUN go build -o main .

WORKDIR /dist

RUN cp /build/main .

FROM scratch

COPY --from=builder /dist/main .

COPY --from=builder /dist/*.properties .

ENV PROFILE=prod \
    DATABASE_HOST=${DATABASE_HOST} \
    DATABASE_NAME=${DATABASE_NAME} \
    DATABASE_USER=${DATABASE_USER} \
    DATABASE_PASSWORD=${DATABASE_PASSWORD} \
    DATABASE_PORT=${DATABASE_PORT}

ENTRYPOINT ["/main"]
