# 참고 : https://dev-racoon.tistory.com/23

FROM golang:alpine AS builder
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build
COPY . .
RUN go mod download
RUN go build -o main .

WORKDIR /dist
RUN cp /build/main .

# Copy .env files to the final stage
COPY .env /dist/.env
COPY .env.production /dist/.env.production

FROM scratch

ENV GO_ENV=production 
# GIN_MODE=release

COPY --from=builder /dist/main .
COPY --from=builder /dist/.env .env
COPY --from=builder /dist/.env.production .env.production

ENTRYPOINT ["/main"]



