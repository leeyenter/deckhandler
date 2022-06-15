# Deck Handler

## Setup

The list of cards that are to be used, and will be loaded into the database,
should be placed in `app/assets` folder. By default, the application reads `app/assets/cards.csv`.

The first row should be the header, and the remaining data should be ordered sequentially.

The first column should be the card codes, which will be used to identify the cards.
The remaining columns are flexible and automatically parsed.

## Running the Code

1. Start the database

```shell
cd database
docker-compose up
```

2. To run tests

```shell
cd app
ginkgo -r
```

3. To run the backend

```shell
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

## APIs

| Method | Path        | Name              | Params                                                                                                                                                                                                    |
| ------ | ----------- | ----------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| POST   | `/`         | Create a new Deck | `cards`: optional; comma-separated list of card codes to include. If not included, will use the full deck.<br />`shuffle`: optional; `true` or `false`, whether to shuffle the deck. Defaults to `false`. |
| GET    | `/:id`      | Open a Deck       | None                                                                                                                                                                                                      |
| POST   | `/:id/draw` | Draw a Card       | `count`: required; number of cards to draw.                                                                                                                                                               |
