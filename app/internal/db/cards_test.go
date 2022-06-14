package db_test

import (
	"github.com/leeyenter/deckhandler/internal/data"
	"github.com/leeyenter/deckhandler/internal/db"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DB card queries", func() {
	var dbObj *db.Database
	var cards []data.Card

	BeforeEach(func() {
		var err error
		dbObj, err = db.GetDB()
		Expect(err).To(BeNil())
		cards, err = data.LoadCSVFile("../../assets/cards.csv")
		Expect(err).To(BeNil())
	})

	Describe("Card management", func() {
		It("should correctly add and retrieve the cards", func() {
			Expect(dbObj.ClearCards()).To(Succeed())

			By("adding the cards")
			Expect(dbObj.CreateCards(cards)).To(Succeed())

			By("retrieving the cards")
			Expect(dbObj.FetchCards()).To(Equal(cards))
		})
	})
})
