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
