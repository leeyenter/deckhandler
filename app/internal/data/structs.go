package data

// Card contains information on a single card.
// ID is the card code, which is extracted from the first column
// of the CSV. The remaining values are stored as a map/JSON dict
// in Values.
type Card struct {
	ID     string
	Values map[string]string
}

// Deck contains information of a current hand of cards,
// including whether they are shuffled.
type Deck struct {
	ID       string
	Shuffled bool
	Cards    []Card
}

// ToMap converts a card into a map, that makes it easier
// to return as a flat JSON object.
func (c *Card) ToMap() map[string]interface{} {
	mappedCard := make(map[string]interface{})
	mappedCard["code"] = c.ID

	for k, v := range c.Values {
		mappedCard[k] = v
	}

	return mappedCard
}
