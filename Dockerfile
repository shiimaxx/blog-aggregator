FROM golang:1.11 as build

ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /go/src/github.com/shiimaxx/blog-aggregator
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -a -tags netgo -installsuffix netgo -o blog-aggregator .


FROM alpine

COPY --from=build /go/src/github.com/shiimaxx/blog-aggregator/blog-aggregator /blog-aggregator
RUN apk add --no-cache --update ca-certificates
EXPOSE 8080
ENTRYPOINT ["/blog-aggregator"]
