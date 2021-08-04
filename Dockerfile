FROM node:lts-alpine as frontend

WORKDIR /frontend-build

COPY ./web/package*.json ./
RUN npm install

COPY ./web .
RUN npm run build

FROM golang:1.16-alpine as backend

RUN apk add --no-cache gcc musl-dev linux-headers

WORKDIR /backend-build

COPY . .
COPY --from=frontend /frontend-build/public ./web/public

RUN go build -o eth-faucet -ldflags "-w -s"

FROM alpine

RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=backend /backend-build/eth-faucet .

EXPOSE 8080

ENTRYPOINT ["./eth-faucet"]
