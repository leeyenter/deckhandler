package data_test

import (
	"testing"

	"github.com/leeyenter/deckhandler/internal/data"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDeckhandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Internal Data Suite")
}

var _ = Describe("Load CSV", func() {
	It("should be able to load cards from disk", func() {
		cards, err := data.LoadCSVFile("../../assets/cards.csv")
		// Tests are based on challenge of having a 52-card deck
		Expect(err).To(BeNil())
		Expect(len(cards)).To(Equal(52))          // for a full deck
		Expect(len(cards[0].Values)).To(Equal(2)) // value and suit
	})

	It("should fail if file is not existent", func() {
		cards, err := data.LoadCSVFile("fakefile.csv")
		Expect(err).ToNot(BeNil())
		Expect(cards).To(BeNil())
	})
})

var _ = Describe("ToMap", func() {
	It("should correctly format a card to map", func() {
		card := data.Card{
			ID: "hello",
			Values: map[string]string{
				"value1": "a",
				"value2": "b",
				"value3": "c",
			},
		}

		cardMap := card.ToMap()

		Expect(cardMap).To(Equal(map[string]interface{}{
			"code":   "hello",
			"value1": "a",
			"value2": "b",
			"value3": "c",
		}))
	})
})
