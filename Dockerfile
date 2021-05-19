FROM golang:1.16-alpine as builder

# ENV GO111MODULE=on \
#    GOPROXY=https://goproxy.cn,direct

RUN apk add --no-cache gcc musl-dev linux-headers

WORKDIR /go/src/github.com/chainflag/eth-faucet

COPY . .

RUN go build .

FROM alpine

RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=builder /go/src/github.com/chainflag/eth-faucet/eth-faucet .

EXPOSE 8080

ENTRYPOINT ["./eth-faucet"]
