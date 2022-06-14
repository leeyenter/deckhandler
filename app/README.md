1. Start the database

```shell
docker run --rm --name db \
    -e POSTGRES_USER=deckhandler \
    -e POSTGRES_PASSWORD=wmWLWyoqsKJtXwisAqwaPkA9yT8MvrzRj \
    -p 127.0.0.1:5432:5432 \
    -d postgres
```

To run tests

```
ginkgo -r
```
