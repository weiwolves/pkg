// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config_test

// TODO use github.com/fortytw2/leaktest

import (
	"io/ioutil"
	goLog "log"
	"sync"
	"testing"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/log/logw"
	"github.com/weiwolves/pkg/config"
	"github.com/weiwolves/pkg/config/storage"
	"github.com/weiwolves/pkg/store/scope"
	"github.com/weiwolves/pkg/util/assert"
	"github.com/fortytw2/leaktest"
)

// those tests cannot run in  because of reading and writing the debug log :-(

var _ config.MessageReceiver = (*testSubscriber)(nil)

type testSubscriber struct {
	t *testing.T
	f func(p config.Path) error
}

func (ts *testSubscriber) MessageConfig(p config.Path) error {
	// ts.t.Logf("Message: %s ScopeGroup %s ScopeID %d", p.String(), p.Scope.String(), p.ID)
	return ts.f(p)
}

func initLogger() (*log.MutexBuffer, log.Logger) {
	var debugBuf log.MutexBuffer
	lg := logw.NewLog(
		logw.WithDebug(&debugBuf, "testDebug: ", goLog.Lshortfile),
		logw.WithInfo(ioutil.Discard, "testInfo: ", goLog.Lshortfile),
		logw.WithLevel(logw.LevelDebug),
	)
	return &debugBuf, lg
}

func TestPubSubBubbling(t *testing.T) {
	testPath := config.MustMakePath("aa/bb/cc")

	s := config.MustNewService(storage.NewMap(), config.Options{
		EnablePubSub: true,
	})

	_, err := s.Subscribe("", nil)
	assert.True(t, errors.Empty.Match(err), "Error: %s", err)

	subID, err := s.Subscribe(testPath.String(), &testSubscriber{
		t: t,
		f: func(p config.Path) error {
			assert.Exactly(t, testPath.BindWebsite(123).String(), p.String(), "In closure Exactly")
			scp, id := p.ScopeID.Unpack()
			if scp == scope.Default {
				assert.Equal(t, int64(0), id)
			} else {
				assert.Equal(t, int64(123), id)
			}
			return nil
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, subID, "The very first subscription ID should be 1")

	assert.NoError(t, s.Set(testPath.BindWebsite(123), []byte(`1`)))
	assert.NoError(t, s.Close())

	// send on closed channel
	assert.NoError(t, s.Set(testPath.BindWebsite(3), []byte(`1`)))
	err = s.Close()
	assert.True(t, errors.AlreadyClosed.Match(err), "Error: %s", err)
}

func TestPubSubPanicSimple(t *testing.T) {
	defer leaktest.Check(t)()

	debugBuf, logger := initLogger()
	s := config.MustNewService(storage.NewMap(), config.Options{
		EnablePubSub: true,
		Log:          logger,
	})

	testPath := config.MustMakePath("xx/yy/zz")

	subID, err := s.Subscribe(testPath.BindStore(123).String(), &testSubscriber{
		t: t,
		f: func(_ config.Path) error {
			panic("Don't panic!")
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, subID, "The very first subscription ID should be 1")
	assert.NoError(t, s.Set(testPath.BindStore(123), []byte(`321`)), "Writing value 123 should not fail")
	assert.NoError(t, s.Close(), "Closing the service should not fail.")
	time.Sleep(time.Millisecond)
	assert.Contains(t, debugBuf.String(), `config.pubSub.publish.recover.r pubSub: true recover: "Don't panic!"`)
}

func TestPubSubPanicError(t *testing.T) {
	// defer leaktest.Check(t)()

	debugBuf, logger := initLogger()
	s := config.MustNewService(storage.NewMap(), config.Options{
		EnablePubSub: true,
		Log:          logger,
	})

	testPath := config.MustMakePath("aa/bb/cc")

	pErr := errors.New("OMG! Panic!")

	subID, err := s.Subscribe(testPath.BindStore(123).String(), &testSubscriber{
		t: t,
		f: func(_ config.Path) error {
			panic(pErr)
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, subID, "The very first subscription ID should be 1")
	assert.NoError(t, s.Set(testPath.BindStore(123), []byte(`321`)))

	assert.NoError(t, s.Close())
	time.Sleep(time.Millisecond)
	assert.Contains(t, debugBuf.String(), `config.pubSub.publish.recover.err pubSub: true error: "OMG! Panic!"`)
}

func TestPubSubPanicMultiple(t *testing.T) {
	debugBuf, logger := initLogger()
	s := config.MustNewService(storage.NewMap(), config.Options{
		EnablePubSub: true,
		Log:          logger,
	})

	subID, err := s.Subscribe("default/0/xx", &testSubscriber{
		t: t,
		f: func(p config.Path) error {
			assert.Equal(t, `default/0/xx/yy/zz`, p.String())
			id, err := p.ScopeID.ID()
			assert.NoError(t, err)
			assert.Exactly(t, uint32(0), id)
			panic("One: Don't panic!")
		},
	})
	assert.NoError(t, err)
	assert.True(t, subID > 0)

	subID, err = s.Subscribe("default/0/xx/yy", &testSubscriber{
		t: t,
		f: func(p config.Path) error {
			assert.Equal(t, "default/0/xx/yy/zz", p.String())
			id, err := p.ScopeID.ID()
			assert.NoError(t, err)
			assert.Exactly(t, uint32(0), id)
			panic("Two: Don't panic!")
		},
	})
	assert.NoError(t, err)
	assert.True(t, subID > 0)

	subID, err = s.Subscribe("default/0/xx/yy/zz", &testSubscriber{
		t: t,
		f: func(p config.Path) error {
			assert.Equal(t, "default/0/xx/yy/zz", p.String())
			id, err := p.ScopeID.ID()
			assert.NoError(t, err)
			assert.Exactly(t, uint32(0), id)
			panic("Three: Don't panic!")
		},
	})
	assert.NoError(t, err)
	assert.True(t, subID > 0)

	assert.NoError(t, s.Set(config.MustMakePath("xx/yy/zz"), []byte(`any kind of data`)))
	assert.NoError(t, s.Close())
	time.Sleep(time.Millisecond)
	// t.Log(debugBuf.String())
	assert.Contains(t, debugBuf.String(), `config.pubSub.publish.recover.r pubSub: true recover: "One: Don't panic!"`)
	assert.Contains(t, debugBuf.String(), `config.pubSub.publish.recover.r pubSub: true recover: "Two: Don't panic!"`)
	assert.Contains(t, debugBuf.String(), `config.pubSub.publish.recover.r pubSub: true recover: "Three: Don't panic!"`)
}

func TestPubSubUnsubscribe(t *testing.T) {
	debugBuf, logger := initLogger()
	s := config.MustNewService(storage.NewMap(), config.Options{
		EnablePubSub: true,
		Log:          logger,
	})

	p := config.MustMakePath("xx/yy/zz").BindStore(123)
	pErr := errors.New("WTF? Panic!")
	subID, err := s.Subscribe(p.String(), &testSubscriber{
		t: t,
		f: func(_ config.Path) error {
			panic(pErr)
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, subID, "The very first subscription ID should be 1")
	assert.NoError(t, s.Unsubscribe(subID))
	assert.NoError(t, s.Set(p, []byte(`any kind of data`)))
	assert.NoError(t, s.Close())
	assert.Regexp(t, `config.Service.Set duration: [0-9]+ path: "stores/123/xx/yy/zz" data_length: 16`, debugBuf.String())
}

type levelCalls struct {
	sync.Mutex
	level2Calls int
	level3Calls int
}

func TestPubSubEvict(t *testing.T) {
	debugBuf, logger := initLogger()
	s := config.MustNewService(storage.NewMap(), config.Options{
		EnablePubSub: true,
		Log:          logger,
	})

	levelCall := new(levelCalls)

	pErr := errors.New("WTF Eviction? Panic!")

	subID, err := s.Subscribe("stores/123/xx/yy", &testSubscriber{
		t: t,
		f: func(p config.Path) error {
			assert.Contains(t, p.String(), "xx/yy")
			// this function gets called 3 times
			levelCall.Lock()
			levelCall.level2Calls++
			levelCall.Unlock()
			return nil
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, subID)

	subID, err = s.Subscribe("stores/123/xx/yy/zz", &testSubscriber{
		t: t,
		f: func(p config.Path) error {
			assert.Contains(t, p.String(), "xx/yy/zz")
			levelCall.Lock()
			levelCall.level3Calls++
			levelCall.Unlock()
			// this function gets called 1 times and then gets removed
			panic(pErr)
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, subID)

	assert.NoError(t, s.Set(config.MustMakePath("xx/yy/zz").BindStore(123), []byte(`321`)))
	assert.NoError(t, s.Set(config.MustMakePath("xx/yy/aa").BindStore(123), []byte(`321`)))
	assert.NoError(t, s.Set(config.MustMakePath("xx/yy/zz").BindStore(123), []byte(`321`)))

	assert.NoError(t, s.Close())
	time.Sleep(time.Millisecond)
	assert.Contains(t, debugBuf.String(), `config.pubSub.publish.recover.err pubSub: true error: "WTF Eviction? Panic!" path: "stores/123/xx/yy/zz"`)

	levelCall.Lock()
	// t.Log(levelCall.level2Calls,levelCall.level3Calls)
	assert.Equal(t, 3, levelCall.level2Calls)
	assert.Equal(t, 1, levelCall.level3Calls)
	levelCall.Unlock()
	err = s.Close()
	assert.True(t, errors.AlreadyClosed.Match(err), "Error: %s", err)
}
