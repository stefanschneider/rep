package models_test

import (
	. "github.com/cloudfoundry-incubator/runtime-schema/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LongRunningProcess", func() {
	Describe("LRPStartAuction", func() {
		var startAuction LRPStartAuction

		startAuctionPayload := `{
    "guid":"some-guid",
    "instance_guid":"some-instance-guid",
    "stack":"some-stack",
    "memory_mb" : 128,
    "disk_mb" : 512,
    "ports": [
      { "container_port": 8080 },
      { "container_port": 8081, "host_port": 1234 }
    ],
    "actions":[
      {
        "action":"download",
        "args":{
          "from":"old_location",
          "to":"new_location",
          "cache_key":"the-cache-key",
          "extract":true
        }
      }
    ],
    "log": {
      "guid": "123",
      "source_name": "APP",
      "index": 42
    },
    "index": 2,
    "state": 1
  }`

		BeforeEach(func() {
			index := 42

			startAuction = LRPStartAuction{
				Guid:         "some-guid",
				InstanceGuid: "some-instance-guid",
				Stack:        "some-stack",
				MemoryMB:     128,
				DiskMB:       512,
				Ports: []PortMapping{
					{ContainerPort: 8080},
					{ContainerPort: 8081, HostPort: 1234},
				},
				Actions: []ExecutorAction{
					{
						Action: DownloadAction{
							From:     "old_location",
							To:       "new_location",
							CacheKey: "the-cache-key",
							Extract:  true,
						},
					},
				},
				Log: LogConfig{
					Guid:       "123",
					SourceName: "APP",
					Index:      &index,
				},
				Index: 2,
				State: LRPStartAuctionStatePending,
			}
		})

		Describe("ToJSON", func() {
			It("should JSONify", func() {
				json := startAuction.ToJSON()
				Ω(string(json)).Should(MatchJSON(startAuctionPayload))
			})
		})

		Describe("NewLongRunningProcessFromJSON", func() {
			It("returns a LongRunningProcess with correct fields", func() {
				decodedStartAuction, err := NewLRPStartAuctionFromJSON([]byte(startAuctionPayload))
				Ω(err).ShouldNot(HaveOccurred())

				Ω(decodedStartAuction).Should(Equal(startAuction))
			})

			Context("with an invalid payload", func() {
				It("returns the error", func() {
					decodedStartAuction, err := NewLRPStartAuctionFromJSON([]byte("butts lol"))
					Ω(err).Should(HaveOccurred())

					Ω(decodedStartAuction).Should(BeZero())
				})
			})
		})
	})

	Describe("Transitional LRP", func() {
		var longRunningProcess TransitionalLongRunningProcess

		longRunningProcessPayload := `{
    "guid":"some-guid",
    "stack":"some-stack",
    "memory_mb" : 128,
    "disk_mb" : 512,
    "ports": [
      { "container_port": 8080 },
      { "container_port": 8081, "host_port": 1234 }
    ],
    "actions":[
      {
        "action":"download",
        "args":{
          "from":"old_location",
          "to":"new_location",
          "cache_key":"the-cache-key",
          "extract":true
        }
      }
    ],
    "log": {
      "guid": "123",
      "source_name": "APP",
      "index": 42
    },
    "state": 1
  }`

		BeforeEach(func() {
			index := 42

			longRunningProcess = TransitionalLongRunningProcess{
				Guid:     "some-guid",
				Stack:    "some-stack",
				MemoryMB: 128,
				DiskMB:   512,
				Ports: []PortMapping{
					{ContainerPort: 8080},
					{ContainerPort: 8081, HostPort: 1234},
				},
				Actions: []ExecutorAction{
					{
						Action: DownloadAction{
							From:     "old_location",
							To:       "new_location",
							CacheKey: "the-cache-key",
							Extract:  true,
						},
					},
				},
				Log: LogConfig{
					Guid:       "123",
					SourceName: "APP",
					Index:      &index,
				},
				State: TransitionalLRPStateDesired,
			}
		})

		Describe("ToJSON", func() {
			It("should JSONify", func() {
				json := longRunningProcess.ToJSON()
				Ω(string(json)).Should(MatchJSON(longRunningProcessPayload))
			})
		})

		Describe("NewLongRunningProcessFromJSON", func() {
			It("returns a LongRunningProcess with correct fields", func() {
				decodedLongRunningProcess, err := NewTransitionalLongRunningProcessFromJSON([]byte(longRunningProcessPayload))
				Ω(err).ShouldNot(HaveOccurred())

				Ω(decodedLongRunningProcess).Should(Equal(longRunningProcess))
			})

			Context("with an invalid payload", func() {
				It("returns the error", func() {
					decodedLongRunningProcess, err := NewTransitionalLongRunningProcessFromJSON([]byte("butts lol"))
					Ω(err).Should(HaveOccurred())

					Ω(decodedLongRunningProcess).Should(BeZero())
				})
			})
		})
	})
})