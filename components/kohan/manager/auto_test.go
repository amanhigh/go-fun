package manager

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// These specs register on Ginkgo's global suite, which the package's
// TestManager entrypoint runs. They live in the internal `manager` package so
// they can reach the unexported recoverConnection and recovery state directly.
var _ = Describe("Auto", func() {
	Context("OS network manager", func() {
		Describe("recoverConnection", func() {
			var (
				mgr *OSManagerImpl

				reconnectCalls []string
				restartCalls   int

				reconnectErr error
			)

			// fakeReconnect records the connection it was asked to reconnect and
			// returns the configured error.
			reconnect := func(conn string) error {
				reconnectCalls = append(reconnectCalls, conn)
				return reconnectErr
			}

			// fakeRestart counts invocations and always succeeds.
			restart := func() error {
				restartCalls++
				return nil
			}

			BeforeEach(func() {
				// Minimal manager: recovery state lives on the struct and is exercised
				// directly via recoverConnection, so no scheduler is required.
				mgr = &OSManagerImpl{}
				reconnectCalls = nil
				restartCalls = 0
				reconnectErr = nil
			})

			// probe drives n consecutive unreachable gateway checks through the policy.
			probe := func(n int) {
				for range n {
					mgr.recoverConnection("wlan0", false, reconnect, restart)
				}
			}

			Context("failure counting", func() {
				It("does not recover after a single failure", func() {
					probe(1)

					Expect(reconnectCalls).To(BeEmpty(), "reconnect must not run on the first failure")
					Expect(restartCalls).To(Equal(0), "restart must never run without a reconnect attempt")
					Expect(mgr.consecutiveFailures).To(Equal(1), "failure count should be tracked")
				})

				It("invokes targeted reconnect (and no restart) after two failures", func() {
					probe(2)

					Expect(reconnectCalls).To(Equal([]string{"wlan0"}), "reconnect should target the resolved connection")
					Expect(restartCalls).To(Equal(0), "restart must not run when reconnect succeeds")
					Expect(mgr.consecutiveFailures).To(Equal(0), "failure count resets after a recovery attempt")
				})
			})

			Context("reconnect failure fallback", func() {
				BeforeEach(func() {
					reconnectErr = errors.New("reconnect boom")
				})

				It("invokes restart exactly once when reconnect fails", func() {
					probe(2)

					Expect(reconnectCalls).To(Equal([]string{"wlan0"}), "reconnect attempted once")
					Expect(restartCalls).To(Equal(1), "restart must run exactly once as fallback")
				})
			})

			Context("cooldown", func() {
				It("suppresses recovery while still within the cooldown window", func() {
					// First recovery attempt records lastRecovery = now.
					probe(2)
					Expect(reconnectCalls).To(Equal([]string{"wlan0"}))

					// More failures immediately after must not trigger another recovery
					// because the cooldown has not elapsed.
					probe(2)
					Expect(reconnectCalls).To(Equal([]string{"wlan0"}), "cooldown should suppress repeated recovery")
					Expect(restartCalls).To(Equal(0))
				})
			})

			Context("state reset", func() {
				It("resets failure count after a recovery attempt", func() {
					probe(2)
					Expect(mgr.consecutiveFailures).To(Equal(0), "recovered state should be reset")

					// A single subsequent failure must not recover, proving the counter
					// was reset to zero rather than left at two.
					probe(1)
					Expect(reconnectCalls).To(Equal([]string{"wlan0"}), "no recovery after a single post-reset failure")
				})

				It("resets failure count when the gateway becomes reachable", func() {
					probe(1)
					Expect(mgr.consecutiveFailures).To(Equal(1))

					mgr.recoverConnection("wlan0", true, reconnect, restart)

					Expect(mgr.consecutiveFailures).To(Equal(0), "reachable gateway clears failure state")
					Expect(reconnectCalls).To(BeEmpty(), "no recovery action on a healthy gateway")
					Expect(restartCalls).To(Equal(0))
				})
			})
		})
	})
})
