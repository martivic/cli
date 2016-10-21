package integration

import (
	. "code.cloudfoundry.org/cli/integration/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = FDescribe("auth command", func() {
	var (
		validUserName   string
		validPassword   string
		invalidUserName string
		invalidPassword string
		session         *Session
	)

	BeforeEach(func() {
		// Skip("until #126256907")
		validUserName = "admin"
		validPassword = "admin"
		invalidUserName = "invalid-username"
		invalidPassword = "invalid-password"
	})

	JustBeforeEach(func() {
		session = CF("auth", validUserName, validPassword)
	})

	FContext("when no API endpoint is set", func() {
		BeforeEach(func() {
			unsetAPI()
		})

		It("fails with no API endpoint set message", func() {
			Eventually(session).Should(Exit(1))
			Expect(session.Out).To(Say("FAILED"))
			Expect(session.Err).To(Say("No API endpoint set. Use 'cf login' or 'cf api' to target an endpoint."))
		})
	})

	Context("when an API endpoint is set", func() {
		BeforeEach(func() {
			setAPI()
		})

		It("displays authenticating message", func() {
			Eventually(session).Should(Exit())
			Expect(session.Out).To(Say("API endpoint:"))
			Expect(session.Out).To(Say("Authenticating..."))
		})

		Context("when invalid credentials are provided", func() {
			JustBeforeEach(func() {
				session = CF("auth", invalidUserName, invalidPassword)
			})

			It("displays an error message", func() {
				Eventually(session).Should(Exit(1))
				Expect(session.Out).To(Say("FAILED"))
				Expect(session.Err).To(Say("Credentials were rejected, please try again."))
			})
		})

		Context("when valid credentials are provided", func() {
			It("logs the user in and displays a success message", func() {
				Eventually(session).Should(Exit(0))
				Expect(session.Out).To(Say("OK"))
				Expect(session.Out).To(Say("Use 'cf target' to view or set your target org and space"))
			})
		})
	})
})
