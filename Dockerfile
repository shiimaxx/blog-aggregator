FROM golang:1.11 as build

ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /go/src/app
COPY . .
RUN GOOS=linux GOARCH=amd64 go build


FROM alpine

COPY --from=build /go/src/app/blog-aggregator /blog-aggregator
EXPOSE 8080
ENTRYPOINT ["/blog-aggregator"]
