// This file was generated by counterfeiter
package fake_snapshot

import (
	"sync"

	"github.com/cloudfoundry-incubator/rep/snapshot"
	"github.com/pivotal-golang/lager"
)

type FakeTaskProcessor struct {
	ProcessStub        func(lager.Logger, *snapshot.TaskSnapshot)
	processMutex       sync.RWMutex
	processArgsForCall []struct {
		arg1 lager.Logger
		arg2 *snapshot.TaskSnapshot
	}
}

func (fake *FakeTaskProcessor) Process(arg1 lager.Logger, arg2 *snapshot.TaskSnapshot) {
	fake.processMutex.Lock()
	fake.processArgsForCall = append(fake.processArgsForCall, struct {
		arg1 lager.Logger
		arg2 *snapshot.TaskSnapshot
	}{arg1, arg2})
	fake.processMutex.Unlock()
	if fake.ProcessStub != nil {
		fake.ProcessStub(arg1, arg2)
	}
}

func (fake *FakeTaskProcessor) ProcessCallCount() int {
	fake.processMutex.RLock()
	defer fake.processMutex.RUnlock()
	return len(fake.processArgsForCall)
}

func (fake *FakeTaskProcessor) ProcessArgsForCall(i int) (lager.Logger, *snapshot.TaskSnapshot) {
	fake.processMutex.RLock()
	defer fake.processMutex.RUnlock()
	return fake.processArgsForCall[i].arg1, fake.processArgsForCall[i].arg2
}

var _ snapshot.TaskProcessor = new(FakeTaskProcessor)