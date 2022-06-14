package db_test

import (
	"github.com/leeyenter/deckhandler/internal/data"
	"github.com/leeyenter/deckhandler/internal/db"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DB deck queries", func() {
	var dbObj *db.Database
	var cards []data.Card

	BeforeEach(func() {
		var err error
		dbObj, err = db.GetDB()
		Expect(err).To(BeNil())
		cards, err = data.LoadCSVFile("../../assets/cards.csv")
		Expect(err).To(BeNil())
	})

	Describe("Deck management", func() {
		It("should be able to create a shuffled deck", func() {
			By("creating a new deck")

			// Make a copy of input, as shuffle will be done in-place,
			// and we don't want it to affect our 'reference'
			cardsInput := make([]data.Card, len(cards))
			copy(cardsInput, cards)

			id, err := dbObj.CreateDeck(true, cardsInput)
			Expect(err).To(BeNil())
			Expect(id).ToNot(BeEmpty())

			By("retrieving information on the deck")
			deck, err := dbObj.GetDeck(id)
			Expect(err).To(BeNil())
			Expect(deck.ID).To(Equal(id))
			Expect(deck.Shuffled).To(Equal(true))
			Expect(len(deck.Cards)).To(Equal(len(cards)))
			Expect(deck.Cards).NotTo(Equal(cards)) // should be shuffled
		})

		It("should be able to create a smaller deck", func() {
			By("creating a new deck")
			subset := cards[:int(len(cards)/2)]
			id, err := dbObj.CreateDeck(false, subset)
			Expect(id).ToNot(BeEmpty())
			Expect(err).To(BeNil())

			By("retrieving information on the deck")
			Expect(dbObj.GetDeck(id)).To(Equal(data.Deck{
				ID:       id,
				Shuffled: false,
				Cards:    subset,
			}))
		})

		It("should be able to create an unshuffled deck", func() {
			By("creating a new deck")
			id, err := dbObj.CreateDeck(false, cards)
			Expect(id).ToNot(BeEmpty())
			Expect(err).To(BeNil())

			By("retrieving information on the deck")
			Expect(dbObj.GetDeck(id)).To(Equal(data.Deck{
				ID:       id,
				Shuffled: false,
				Cards:    cards,
			}))

			By("fetching all the cards in the deck")
			Expect(dbObj.FetchCardsFromDeck(id, -1)).To(Equal(cards))

			By("fetching 5 cards in the deck")
			Expect(dbObj.FetchCardsFromDeck(id, 5)).To(Equal(cards[:5]))

			By("removing the first 5 cards in the deck")
			cardCodes := make([]string, 0)
			for _, card := range cards[:5] {
				cardCodes = append(cardCodes, card.ID)
			}
			Expect(dbObj.RemoveCardsFromDeck(id, cardCodes)).To(Succeed())

			By("fetching 5 more cards in the deck")
			Expect(dbObj.FetchCardsFromDeck(id, 5)).To(Equal(cards[5:10]))

			By("fetching all of the remaining cards in the deck")
			Expect(dbObj.FetchCardsFromDeck(id, -1)).To(Equal(cards[5:]))
		})

	})
})
