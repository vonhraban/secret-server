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

	Context("when a secret 123abc is added", func() {
		vault := persistence.NewInMemoryVault()
		cmd := &secret.AddSecret{}
		_, err := cmd.Execute(vault, "123abc")
		if err != nil {
			panic(err)
		}
		Context("the the secret 123abc should be stored", func() {
			storedSecret, err := vault.Retrieve("123abc")
			if err != nil {
				panic(err)
			}
			It("has a number of current uses set to 0", func() {
				Expect(storedSecret.Uses).To(Equal(0))

			})
			It("has a maximum allowed number of uses set to 5", func() {
				Expect(storedSecret.MaxUses).To(Equal(5))
			})
			It("has the time created set to 2019-06-15 11:14:00", func() {})

			It("has the time expires set to 2019-06-15 12:14:00", func() {})
		})
	})

})
