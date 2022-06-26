package helpers

import (
	"github.com/onsi/gomega/types"
	"go.uber.org/atomic"
)

type FailInterceptor struct {
	ginkgoFail types.GomegaFailHandler
	didFail    *atomic.Bool
}

func NewFailInterceptor(fail types.GomegaFailHandler) *FailInterceptor {
	return &FailInterceptor{
		ginkgoFail: fail,
		didFail:    atomic.NewBool(false),
	}
}

func (f *FailInterceptor) Fail(message string, callerSkip ...int) {
	f.didFail.Store(true)
	if len(callerSkip) == 0 {
		f.ginkgoFail(message, 1)
	} else {
		f.ginkgoFail(message, callerSkip[0]+1)
	}
}

func (f *FailInterceptor) Reset() {
	f.didFail.Store(false)
}

func (f *FailInterceptor) DidFail() bool {
	return f.didFail.Load()
}
