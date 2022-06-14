package router_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDeckhandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Internal Routes Suite")
}

// var _ = Describe("Create Deck", func() {
// 	var r *router.Router
// 	var e *echo.Echo
// 	var rec *httptest.ResponseRecorder

// 	BeforeEach(func() {
// 		r = router.New()
// 		e = echo.New()
// 		rec = httptest.NewRecorder()
// 	})

// 	It("should be able to create a complete, unshuffled deck", func() {
// 		req := httptest.NewRequest(http.MethodPost, "/", nil)
// 		c := e.NewContext(req, rec)
// 		Expect(r.CreateDeck(c)).To(Succeed())

// 		// Expect response to have full set of cards, with no shuffling
// 	})

// 	It("should be able to create a complete, shuffled deck", func() {
// 		req := httptest.NewRequest(http.MethodPost, "/", nil)
// 		c := e.NewContext(req, rec)
// 		Expect(r.CreateDeck(c)).To(Succeed())

// 		// Expect response to have full set of cards, with shuffling
// 	})

// 	It("should be able to create a subsetted deck", func() {
// 		req := httptest.NewRequest(http.MethodPost, "/", nil)
// 		c := e.NewContext(req, rec)
// 		c.SetPath("/")
// 		Expect(r.CreateDeck(c)).To(Succeed())

// 		// Expect response to have subset set of cards, with no shuffling
// 	})
// })
