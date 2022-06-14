package router_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/leeyenter/deckhandler/internal/data"
	"github.com/leeyenter/deckhandler/internal/db"
	"github.com/leeyenter/deckhandler/internal/router"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDeckhandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Internal Routes Suite")
}

func parseDeck(rec *httptest.ResponseRecorder) router.DeckResponse {
	var resp router.DeckResponse
	Expect(json.Unmarshal(rec.Body.Bytes(), &resp)).To(Succeed())
	return resp
}

var _ = Describe("Create Deck", func() {
	var r *router.Router
	var e *echo.Echo
	var rec *httptest.ResponseRecorder
	var cards []data.Card

	BeforeEach(func() {
		r = router.New()
		e = echo.New()
		rec = httptest.NewRecorder()

		var err error
		cards, err = data.LoadCSVFile("../../assets/cards.csv")
		Expect(err).To(BeNil())
	})

	It("should create an unshuffled deck if shuffled is not passed in", func() {
		By("making the request")
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		c := e.NewContext(req, rec)
		Expect(r.CreateDeck(c)).To(Succeed())

		resp := parseDeck(rec)
		Expect(resp.DeckID).NotTo(BeEmpty())
		Expect(resp.Shuffled).To(Equal(false))
		Expect(resp.Remaining).To(Equal(len(cards)))
	})

	It("should create an unshuffled deck if shuffled is 'false'", func() {
		req := httptest.NewRequest(http.MethodPost, "/?shuffled=false", nil)
		c := e.NewContext(req, rec)
		Expect(r.CreateDeck(c)).To(Succeed())

		resp := parseDeck(rec)
		Expect(resp.DeckID).NotTo(BeEmpty())
		Expect(resp.Shuffled).To(Equal(false))
		Expect(resp.Remaining).To(Equal(len(cards)))
	})

	It("should create a shuffled deck if shuffled is 'true'", func() {
		req := httptest.NewRequest(http.MethodPost, "/?shuffled=true", nil)
		c := e.NewContext(req, rec)
		Expect(r.CreateDeck(c)).To(Succeed())

		resp := parseDeck(rec)
		Expect(resp.DeckID).NotTo(BeEmpty())
		Expect(resp.Shuffled).To(Equal(true))
		Expect(resp.Remaining).To(Equal(len(cards)))
	})

	It("should fail if shuffled is set to an unrecognised value", func() {
		req := httptest.NewRequest(http.MethodPost, "/?shuffled=unrecognised", nil)
		c := e.NewContext(req, rec)
		Expect(r.CreateDeck(c)).ToNot(Succeed())
	})

	It("should be able to create a subsetted deck", func() {
		url := "/?cards="
		for i, card := range cards[:10] {
			if i > 0 {
				url += ","
			}
			url += card.ID
		}

		req := httptest.NewRequest(http.MethodPost, url, nil)
		c := e.NewContext(req, rec)
		Expect(r.CreateDeck(c)).To(Succeed())
	})

	It("Should fail if cards are not recognised", func() {
		req := httptest.NewRequest(http.MethodPost, "/?cards=unknown,unknown2", nil)
		c := e.NewContext(req, rec)
		Expect(r.CreateDeck(c)).ToNot(Succeed())
	})
})

var _ = Describe("Open Deck", func() {
	var r *router.Router
	var e *echo.Echo
	var rec *httptest.ResponseRecorder

	BeforeEach(func() {
		r = router.New()
		e = echo.New()
		rec = httptest.NewRecorder()
	})

	Context("A deck has been created", func() {
		var err error
		cards, err := data.LoadCSVFile("../../assets/cards.csv")
		Expect(err).To(BeNil())
		dbObj, err := db.GetDB()
		Expect(err).To(BeNil())
		id, err := dbObj.CreateDeck(false, cards)
		Expect(err).To(BeNil())
		Expect(id).ToNot(BeEmpty())

		It("should be able to open the newly created deck", func() {
			req := httptest.NewRequest(http.MethodGet, "/"+id, nil)
			c := e.NewContext(req, rec)
			Expect(r.OpenDeck(c)).To(Succeed())
			resp := parseDeck(rec)
			Expect(resp.Shuffled).To(Equal(false))
			Expect(resp.Remaining).To(Equal(len(cards)))
			Expect(len(resp.Cards)).To(Equal(resp.Remaining))

			It("should have the correct cards", func() {
				for _, card := range resp.Cards {
					found := false
					for _, refCard := range cards {
						if card["code"] == refCard.ID {
							fail := false
							for k, v := range refCard.Values {
								if card[k] != v {
									fail = true
								}
							}

							if !fail {
								found = true
							}
						}
					}

					Expect(found).To(Equal(true), "no corresponding card found")
				}
			})
		})
	})

	It("should fail if not given an id", func() {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		c := e.NewContext(req, rec)
		Expect(r.OpenDeck(c)).ToNot(Succeed())
	})

	It("should fail if given an invalid id", func() {
		req := httptest.NewRequest(http.MethodGet, "/invalid", nil)
		c := e.NewContext(req, rec)
		Expect(r.OpenDeck(c)).ToNot(Succeed())
	})
})

var _ = Describe("Draw a Card", func() {
	var r *router.Router
	var e *echo.Echo
	var rec *httptest.ResponseRecorder
	var cards []data.Card
	var id string

	BeforeEach(func() {
		r = router.New()
		e = echo.New()
		rec = httptest.NewRecorder()

		var err error
		cards, err = data.LoadCSVFile("../../assets/cards.csv")
		Expect(err).To(BeNil())

		dbObj, err := db.GetDB()
		Expect(err).To(BeNil())
		id, err = dbObj.CreateDeck(false, cards)
		Expect(err).To(BeNil())
	})

	Context("A deck has been created", func() {
		It("should draw the correct number of cards", func() {
			By("Sending the draw request")
			req := httptest.NewRequest(http.MethodPost, "/"+id+"/draw?count=5", nil)
			c := e.NewContext(req, rec)
			Expect(r.DrawCards(c)).To(Succeed())
			Expect(len(parseDeck(rec).Cards)).To(Equal(5))

			By("opening the deck again")
			rec = httptest.NewRecorder()
			req = httptest.NewRequest(http.MethodGet, "/"+id, nil)
			c = e.NewContext(req, rec)
			Expect(r.OpenDeck(c)).To(Succeed())
			Expect(len(parseDeck(rec).Cards)).To(Equal(len(cards) - 5))
		})

		It("should fail if a count parameter is not provided", func() {
			By("sending the draw request")
			req := httptest.NewRequest(http.MethodPost, "/"+id+"/draw", nil)
			c := e.NewContext(req, rec)
			Expect(r.DrawCards(c)).ToNot(Succeed())

			By("opening the deck again")
			rec = httptest.NewRecorder()
			req = httptest.NewRequest(http.MethodGet, "/"+id, nil)
			c = e.NewContext(req, rec)
			Expect(r.OpenDeck(c)).To(Succeed())
			Expect(len(parseDeck(rec).Cards)).To(Equal(len(cards)))
		})

		It("should fail if too many cards are drawn", func() {
			By("sending the draw request")
			req := httptest.NewRequest(http.MethodPost, "/"+id+"/draw?count="+strconv.Itoa(len(cards)+1), nil)
			c := e.NewContext(req, rec)
			Expect(r.DrawCards(c)).ToNot(Succeed())

			By("opening the deck again")
			rec = httptest.NewRecorder()
			req = httptest.NewRequest(http.MethodGet, "/"+id, nil)
			c = e.NewContext(req, rec)
			Expect(r.OpenDeck(c)).To(Succeed())
			Expect(len(parseDeck(rec).Cards)).To(Equal(len(cards)))
		})
	})

	It("should fail if not given an id", func() {
		req := httptest.NewRequest(http.MethodPost, "/draw?count=5", nil)
		c := e.NewContext(req, rec)
		Expect(r.DrawCards(c)).ToNot(Succeed())
	})

	It("should fail if given an invalid id", func() {
		req := httptest.NewRequest(http.MethodPost, "/invalid-id/draw?count=5", nil)
		c := e.NewContext(req, rec)
		Expect(r.DrawCards(c)).ToNot(Succeed())
	})
})
