package data

type Card struct {
	ID     string
	Values map[string]string
}

type Deck struct {
	ID       string
	Shuffled bool
	Cards    []Card
}

// ToMap converts a card into a JSON-suitable object
func (c *Card) ToMap() map[string]interface{} {
	mappedCard := make(map[string]interface{})
	mappedCard["code"] = c.ID

	for k, v := range c.Values {
		mappedCard[k] = v
	}

	return mappedCard
}
