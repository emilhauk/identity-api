FROM golang:alpine

WORKDIR /app
COPY . .
RUN go get ./...
RUN go build

ENTRYPOINT ["./identity-api"]

EXPOSE 9002
