CREATE TABLE cards (
    code TEXT NOT NULL PRIMARY KEY,
    value JSONB NOT NULL UNIQUE
);

CREATE TABLE decks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shuffled BOOLEAN NOT NULL
);

CREATE TABLE deck_cards (
    deck_id UUID NOT NULL REFERENCES decks(id) ON DELETE CASCADE,
    card_code TEXT NOT NULL REFERENCES cards(code) ON DELETE CASCADE
);