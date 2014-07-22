package maintain_test

import (
	"errors"
	"syscall"
	"time"

	fake_client "github.com/cloudfoundry-incubator/executor/api/fakes"
	"github.com/cloudfoundry-incubator/rep/maintain"
	"github.com/cloudfoundry-incubator/runtime-schema/bbs/fake_bbs"
	"github.com/cloudfoundry-incubator/runtime-schema/models"
	steno "github.com/cloudfoundry/gosteno"
	"github.com/tedsuo/ifrit"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Maintain Presence", func() {
	var (
		executorPresence  models.ExecutorPresence
		heartbeatInterval = 500 * time.Millisecond

		fakeBBS    *fake_bbs.FakeRepBBS
		fakeClient *fake_client.FakeClient
		logger     *steno.Logger

		maintainer ifrit.Process

		presence           *fake_bbs.FakePresence
		maintainStatusChan chan bool
	)

	BeforeSuite(func() {
		steno.EnterTestMode(steno.LOG_DEBUG)
	})

	BeforeEach(func() {
		fakeClient = new(fake_client.FakeClient)

		presence = &fake_bbs.FakePresence{}
		maintainStatusChan = make(chan bool)

		executorPresence = models.ExecutorPresence{
			ExecutorID: "executor-id",
			Stack:      "lucid64",
		}

		fakeBBS = &fake_bbs.FakeRepBBS{}
		fakeBBS.MaintainExecutorPresenceReturns(presence, maintainStatusChan, nil)

		logger = steno.NewLogger("test-logger")

		maintainer = ifrit.Envoke(maintain.New(executorPresence, fakeClient, fakeBBS, logger, heartbeatInterval))
	})

	AfterEach(func() {
		maintainer.Signal(syscall.SIGTERM)
		<-maintainer.Wait()
	})

	Context("when running", func() {
		It("should already have started maintaining presence", func() {
			Ω(fakeBBS.MaintainExecutorPresenceCallCount()).Should(Equal(1))
			interval, maintainedPresence := fakeBBS.MaintainExecutorPresenceArgsForCall(0)
			Ω(interval).Should(Equal(heartbeatInterval))
			Ω(maintainedPresence).Should(Equal(executorPresence))
		})

		It("should ping the executor on each maintain tick", func() {
			maintainStatusChan <- true
			Eventually(fakeClient.PingCallCount).Should(Equal(1))

			maintainStatusChan <- true
			Eventually(fakeClient.PingCallCount).Should(Equal(2))
		})
	})

	Context("when the executor ping fails", func() {
		BeforeEach(func() {
			fakeClient.PingReturns(errors.New("bam"))
			maintainStatusChan <- true
			Eventually(fakeClient.PingCallCount).Should(Equal(1))
		})

		It("should remove presence", func() {
			Eventually(presence.Removed).Should(BeTrue())
		})

		It("should start pinging the executor without relying on its presence being maintained", func() {
			Eventually(fakeClient.PingCallCount, 10*heartbeatInterval).Should(Equal(2))
			Eventually(fakeClient.PingCallCount, 10*heartbeatInterval).Should(Equal(3))
		})

		Context("and then the executor ping succeeds", func() {
			var newMaintainStatusChan chan bool

			BeforeEach(func() {
				newMaintainStatusChan = make(chan bool)
				fakeBBS.MaintainExecutorPresenceReturns(presence, newMaintainStatusChan, nil)

				fakeClient.PingReturns(nil) //healthy again
				Eventually(fakeClient.PingCallCount).Should(Equal(2))
			})

			It("should attempt to reestablish presence", func() {
				Eventually(fakeBBS.MaintainExecutorPresenceCallCount).Should(Equal(2))
			})

			It("should ping the executor on each maintain tick", func() {
				Ω(fakeClient.PingCallCount()).Should(Equal(2))
				select {
				case newMaintainStatusChan <- true:
				case <-time.Tick(time.Second):
					Fail("newMaintainStatusChan not called in time")
				}
				Eventually(fakeClient.PingCallCount).Should(Equal(3))
			})
		})
	})

	Context("when we fail to maintain our presence", func() {
		BeforeEach(func() {
			maintainStatusChan <- false
		})

		It("does not shut down", func() {
			Consistently(maintainer.Wait()).ShouldNot(Receive(), "should not shut down")
		})

		It("continues to retry", func() {
			Ω(fakeClient.PingCallCount()).Should(Equal(0))
			maintainStatusChan <- true
			Eventually(fakeClient.PingCallCount).Should(Equal(1))
		})

		It("logs an error message", func() {
			testSink := steno.GetMeTheGlobalTestSink()

			records := []*steno.Record{}

			lockMessageIndex := 0
			Eventually(func() string {
				records = testSink.Records()

				if len(records) > 0 {
					lockMessageIndex := len(records) - 1
					return records[lockMessageIndex].Message
				}

				return ""
			}, 1.0, 0.1).Should(Equal("rep.maintain_presence.lost-lock"))

			Ω(records[lockMessageIndex].Level).Should(Equal(steno.LOG_ERROR))
		})
	})
})
