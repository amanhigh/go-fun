package manager

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// These specs register on Ginkgo's global suite, which the package's
// TestManager entrypoint runs. They live in the internal `manager` package so
// they can reach the unexported recoverConnection and recovery state directly.

// testConnectionName is the connection identifier exercised by the recovery
// policy specs. The value is a synthetic placeholder; only the policy's failure
// counting is under test, not the connection itself.
const testConnectionName = "test-wifi-connection"

var _ = Describe("Auto", func() {
	Context("OS network manager", func() {
		Describe("recoverConnection", func() {
			var (
				mgr          *OSManagerImpl
				restartCalls int
			)

			// restart counts invocations and always succeeds.
			restart := func() error {
				restartCalls++
				return nil
			}

			BeforeEach(func() {
				// Minimal manager: recovery state lives on the struct and is exercised
				// directly via recoverConnection, so no scheduler is required.
				mgr = &OSManagerImpl{}
				restartCalls = 0
			})

			// probe drives n consecutive unreachable gateway checks through the policy.
			probe := func(n int) {
				for range n {
					mgr.recoverConnection(testConnectionName, false, restart)
				}
			}

			Context("restart-only recovery policy", func() {
				It("does not restart after failures just below the threshold", func() {
					probe(networkManagerRestartAfterFailures - 1)

					Expect(restartCalls).To(Equal(0), "restart must not run before the failure threshold")
					Expect(mgr.consecutiveFailures).To(Equal(networkManagerRestartAfterFailures-1), "failure count continues to accumulate")
				})

				It("restarts NetworkManager exactly once at the threshold and resets the counter", func() {
					probe(networkManagerRestartAfterFailures)

					Expect(restartCalls).To(Equal(1), "restart must run once at the failure threshold")
					Expect(mgr.consecutiveFailures).To(Equal(0), "failure count resets after restart")
				})

				It("begins a fresh cycle after a restart so the next threshold of failures does not restart", func() {
					probe(networkManagerRestartAfterFailures)
					Expect(restartCalls).To(Equal(1), "restart must run at the failure threshold")

					// One less than the threshold of failures after the restart must not
					// trigger another restart because the counter resets to zero.
					probe(networkManagerRestartAfterFailures - 1)
					Expect(restartCalls).To(Equal(1), "no restart within the new cycle before the threshold")
					Expect(mgr.consecutiveFailures).To(Equal(networkManagerRestartAfterFailures-1), "failure count accumulates in the new cycle")
				})

				It("resets the failure counter and performs no restart when the gateway is reachable", func() {
					probe(1)
					Expect(mgr.consecutiveFailures).To(Equal(1))

					mgr.recoverConnection(testConnectionName, true, restart)

					Expect(mgr.consecutiveFailures).To(Equal(0), "reachable gateway clears failure state")
					Expect(restartCalls).To(Equal(0), "no recovery action on a healthy gateway")
				})
			})
		})
	})
})
