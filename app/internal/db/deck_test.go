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
			id, err := dbObj.CreateDeck(true, cards)
			Expect(id).ToNot(BeEmpty())
			Expect(err).To(BeNil())

			By("retrieving information on the deck")
			Expect(dbObj.GetDeck(id)).To(Equal(data.Deck{
				ID:       id,
				Shuffled: true,
			}))

			By("fetching the cards in the deck")
			fetchedCards, err := dbObj.FetchCardsFromDeck(id, -1)
			Expect(err).To(BeNil())
			Expect(len(fetchedCards)).To(Equal(len(cards)))
			Expect(fetchedCards).NotTo(Equal(cards)) // should be shuffled
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
			}))

			By("fetching the cards in the deck")
			Expect(dbObj.FetchCardsFromDeck(id, -1)).To(Equal(subset))
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
			}))

			By("fetching all the cards in the deck")
			Expect(dbObj.FetchCardsFromDeck(id, -1)).To(Equal(cards))

			By("fetching 5 cards in the deck")
			Expect(dbObj.FetchCardsFromDeck(id, 5)).To(Equal(cards[:5]))

			By("removing the first 5 cards in the deck")
			Expect(dbObj.RemoveCardsFromDeck(id, 5)).To(Succeed())

			By("fetching 5 more cards in the deck")
			Expect(dbObj.FetchCardsFromDeck(id, 5)).To(Equal(Equal(cards[5:10])))

			By("fetching all of the remaining cards in the deck")
			Expect(dbObj.FetchCardsFromDeck(id, -1)).To(Equal(cards[5:]))
		})

	})
})