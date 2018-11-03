# blog-aggregator

[![wercker status](https://app.wercker.com/status/80c41247bbc1cf9592f13bba6216b6ba/s/master "wercker status")](https://app.wercker.com/project/byKey/80c41247bbc1cf9592f13bba6216b6ba)

## Build

```
docker build -t docker-blog-aggregator .
```

## Run

```
docker run -it --rm --env USER_ID=<user id> -p 8080:8080 docker-blog-aggregator
```

## License

[MIT](https://github.com/shiimaxx/blog-aggregator/blob/master/LICENSE)

## Author

[shiimaxx](https://github.com/shiimaxx)
