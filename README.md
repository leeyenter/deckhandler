1. Start the database

```shell
cd database
docker-compose up
```

2. To run tests

```
cd app
ginkgo -r
```

3. To run the backend

```
cd app
go build
./deckhandler
```
