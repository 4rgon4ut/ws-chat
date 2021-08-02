FROM golang:1.16-buster as builder

# go deps
ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

COPY go.mod go.sum build/

WORKDIR build

RUN go mod download

COPY . .

WORKDIR cmd/server/

RUN CGO_ENABLED=0 go build -o /bin/server



FROM alpine:latest

# time zone
RUN apk add --no-cache tzdata
ENV TZ=Europe/Moscow

WORKDIR /root/

COPY --from=builder /bin/server .

ARG APP_VERSION
ARG COMMIT_SHORT

ENV APP_VERSION $APP_VERSION

LABEL APP_VERSION=$APP_VERSION \
    COMMIT=$COMMIT_SHORT


CMD ["./server"]