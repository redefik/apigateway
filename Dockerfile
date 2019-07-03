#build stage
FROM golang:1.11 AS build-env

WORKDIR /go/src/github.com/redefik/sdccproject/apigateway

COPY . .

RUN go get -d -v ./...

RUN cd cmd && CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o /go/bin/apigateway 

#production stage
FROM alpine:latest

WORKDIR /root/

COPY --from=build-env /go/bin/apigateway .

EXPOSE 80

CMD ["./apigateway"]
