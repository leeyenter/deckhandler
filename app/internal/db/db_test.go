package db_test

import (
	"context"
	"testing"

	"github.com/leeyenter/deckhandler/internal/db"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDatabase(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Internal DB Suite")
}

var _ = Describe("Connect to database", func() {
	It("should be able to connect to the DB", func() {
		dbObj, err := db.GetDB()
		Expect(err).To(BeNil())
		Expect(dbObj.Conn.Ping(context.Background())).To(BeNil())
	})
})
