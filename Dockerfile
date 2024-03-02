# 참고 : https://dev-racoon.tistory.com/23

FROM golang:alpine AS builder
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build
ADD markdown ./markdown
COPY go.mod go.sum main.go ./
RUN go mod download
RUN go build -o main .
WORKDIR /dist
RUN cp /build/main .

FROM scratch
COPY --from=builder /dist/main .
ENTRYPOINT ["/main"]
