package logrushookmonitor

import (
	"context"
	"github.com/sirupsen/logrus"
	"sync"
)

type HookMonitor struct {
	hook          logrus.Hook
	onError       func(error) error
	ctx           context.Context
	ctxCancelFunc context.CancelFunc

	lock sync.RWMutex
}

func New() (self *HookMonitor) {
	self = new(HookMonitor)
	self.onError = func(err error) error { return err }
	self.ctx, self.ctxCancelFunc = context.WithCancel(context.Background())
	return
}

func (self *HookMonitor) SetHook(hook logrus.Hook) *HookMonitor {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.hook = hook
	return self
}

func (self *HookMonitor) SetOnError(onError func(error) error) *HookMonitor {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.onError = onError
	return self
}

func (self *HookMonitor) GetContext() context.Context {
	self.lock.Lock()
	defer self.lock.Unlock()

	return self.ctx
}

// Levels returns all logrus levels.
func (self *HookMonitor) Levels() []logrus.Level {
	self.lock.RLock()
	defer self.lock.RUnlock()

	return self.hook.Levels()
}

func (self *HookMonitor) Fire(entry *logrus.Entry) error {
	self.lock.RLock()
	defer self.lock.RUnlock()

	err := self.hook.Fire(entry)
	err = self.onError(err)
	if err != nil {
		self.ctxCancelFunc()
	}
	return err
}
