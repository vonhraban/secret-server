package secret_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vonhraban/secret-server/persistence"
	"github.com/vonhraban/secret-server/secret"
)

var _ = Describe("Secret", func() {
	Context("initially", func() {
		It("does not have 123abc", func() {
			// nothing to do here, although perhaps ensure?
		})
		It("is 2019-06-15 11:14:00 ", func() {
			// somehow set the time
		})
	})

	Context("when a secret 123abc is added with allowed max views of 5", func() {
		vault := persistence.NewInMemoryVault()
		cmd := &secret.AddSecret{}
		id, err := cmd.Execute(vault, "123abc", 5)
		if err != nil {
			panic(err)
		}
		Context("then this secret should be stored", func() {
			storedSecret, err := vault.Retrieve(id)
			if err != nil {
				panic(err)
			}
			It("has a 4 remaining views since it has been just retrieved", func() { // TODO Don't like this
				Expect(storedSecret.RemainingViews).To(Equal(4))
			})
			It("has the time created set to 2019-06-15 11:14:00", func() {})

			It("has the time expires set to 2019-06-15 12:14:00", func() {})
		})
	})

})
