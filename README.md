# Deck Handler

## Setup

The list of cards that are to be used, and will be loaded into the database,
should be placed in `assets` folder. By default, the application reads `assets/cards.csv`.

The data should be ordered sequentially. The first row should be the header,
and the first column should be the card codes, which will be used to identify
the cards. The remaining columns are flexible and automatically parsed.

## Running the Code

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

The app also uses the following environment variables:

| Variable     | Description                                 | Default                             |
| ------------ | ------------------------------------------- | ----------------------------------- |
| `CARDS_PATH` | Path of file to load into the database      | `assets/cards.csv`                  |
| `SKIP_SEED`  | Whether to skip the seeding of the database | `false`                             |
| `PORT`       | Port to run the web app on                  | `3000`                              |
| `DB_USER`    | Database account username                   | `deckhandler`                       |
| `DB_PASS`    | Database account password                   | `wmWLWyoqsKJtXwisAqwaPkA9yT8MvrzRj` |
| `DB_HOST`    | Database host                               | `127.0.0.1`                         |
| `DB_PORT`    | Database port                               | `5432`                              |
