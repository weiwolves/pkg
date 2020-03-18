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

package ratelimit

import (
	"context"
	"fmt"
	"io"
	"sort"
	"sync"

	"github.com/weiwolves/pkg/config"
	"github.com/weiwolves/pkg/net/mw"
	"github.com/weiwolves/pkg/store/scope"
	"github.com/weiwolves/pkg/sync/singleflight"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

// Auto generated: Do not edit. See net/internal/scopedService package for more details.

type service struct {
	// useWebsite internal flag used in configByContext(w,r) to tell the
	// currenct handler if the scoped configuration is store or website based.
	useWebsite bool
	// optionAfterApply allows to set a custom function which runs every time
	// after the options have been applied. Gets only executed if not nil.
	optionAfterApply func() error

	// rwmu protects all fields below
	rwmu sync.RWMutex
	// scopeCache internal cache for configurations.
	scopeCache map[scope.TypeID]*ScopedConfig
	// optionFactory optional configuration closure, can be nil. It pulls out
	// the configuration settings from a slow backend during a request and
	// caches the settings in the internal map.  This function gets set via
	// WithOptionFactory()
	optionFactory OptionFactoryFunc
	// optionInflight checks on a per scope.TypeID basis if the configuration
	// loading process takes place. Stops the execution of other Goroutines (aka
	// incoming requests) with the same scope.TypeID until the configuration has
	// been fully loaded and applied for that specific scope. This function gets
	// set via WithOptionFactory()
	optionInflight *singleflight.Group
	// ErrorHandler gets called whenever a programmer makes an error. Most two
	// cases are: cannot extract scope from the context and scoped configuration
	// is not valid. The default handler prints the error to the client and
	// returns http.StatusServiceUnavailable
	mw.ErrorHandler
	// Log used for debugging. Defaults to black hole.
	Log log.Logger
	// rootConfig optional backend configuration. Gets only used while running
	// HTTP related middlewares.
	RootConfig config.Getter
}

func newService(opts ...Option) (*Service, error) {
	s := &Service{
		service: service{
			Log:          log.BlackHole{},
			ErrorHandler: defaultErrorHandler,
			scopeCache:   make(map[scope.TypeID]*ScopedConfig),
		},
	}
	if err := s.Options(WithDefaultConfig(scope.DefaultTypeID)); err != nil {
		return nil, errors.Wrap(err, "[ratelimit] Options WithDefaultConfig")
	}
	if err := s.Options(opts...); err != nil {
		return nil, errors.Wrap(err, "[ratelimit] Options any config")
	}
	return s, nil
}

// MustNew same as New() but panics on error. Use only during app start up process.
func MustNew(opts ...Option) *Service {
	c, err := New(opts...)
	if err != nil {
		panic(err)
	}
	return c
}

// Options applies option at creation time or refreshes them.
func (s *Service) Options(opts ...Option) error {
	for _, opt := range opts {
		// opt can be nil because of the backend options where we have an array instead
		// of a slice.
		if opt != nil {
			if err := opt(s); err != nil {
				return errors.Wrap(err, "[ratelimit] Service.Options")
			}
		}
	}
	if s.optionAfterApply != nil {
		return errors.Wrap(s.optionAfterApply(), "[ratelimit] optionValidation")
	}
	return nil
}

// ClearCache clears the internal map storing all scoped configurations. You
// must reapply all functional options.
// TODO(CyS) all previously applied options will be automatically reapplied.
func (s *Service) ClearCache() error {
	s.scopeCache = make(map[scope.TypeID]*ScopedConfig)
	return nil
}

// DebugCache uses Sprintf to write an ordered list (by scope.TypeID) into a
// writer. Only usable for debugging.
func (s *Service) DebugCache(w io.Writer) error {
	s.rwmu.RLock()
	defer s.rwmu.RUnlock()
	srtScope := make(scope.TypeIDs, len(s.scopeCache))
	var i int
	for scp := range s.scopeCache {
		srtScope[i] = scp
		i++
	}
	sort.Sort(srtScope)
	for _, scp := range srtScope {
		scpCfg := s.scopeCache[scp]
		if _, err := fmt.Fprintf(w, "%s => [%p]=%#v\n", scp, scpCfg, scpCfg); err != nil {
			return errors.Wrap(err, "[ratelimit] DebugCache Fprintf")
		}
	}
	return nil
}

// ConfigByScope creates a new scoped configuration depending on the
// Service.useWebsite flag. If useWebsite==true the scoped configuration
// contains only the website->default scope despite setting a store scope. If an
// OptionFactory is set the configuration gets loaded from the backend. A nil
// root config causes a panic.
func (s *Service) ConfigByScope(websiteID, storeID int64) (ScopedConfig, error) {
	cfg := s.RootConfig.NewScoped(websiteID, storeID)
	if s.useWebsite {
		cfg = s.RootConfig.NewScoped(websiteID, 0)
	}
	return s.ConfigByScopedGetter(cfg)
}

// configByContext extracts the scope (websiteID and storeID) from a  context.
// The scoped configuration gets initialized by configFromScope() and returned.
// It panics if rootConfig if nil. Errors get not logged.
func (s *Service) configByContext(ctx context.Context) (ScopedConfig, error) {
	// extract the scope out of the context and if not found a programmer made a
	// mistake.
	websiteID, storeID, scopeOK := scope.FromContext(ctx)
	if !scopeOK {
		return ScopedConfig{}, errors.NewNotFoundf("[ratelimit] configByContext: scope.FromContext not found")
	}

	scpCfg, err := s.ConfigByScope(websiteID, storeID)
	if err != nil {
		// the scoped configuration is invalid and hence a programmer or package user
		// made a mistake.
		return ScopedConfig{}, errors.Wrap(err, "[ratelimit] Service.configByContext.configFromScope") // rewrite error
	}
	return scpCfg, nil
}

// ConfigByScopedGetter returns the internal configuration depending on the
// ScopedGetter. Mainly used within the middleware.  If you have applied the
// option WithOptionFactory() the configuration will be pulled out only one time
// from the backend configuration service. The field optionInflight handles the
// guaranteed atomic single loading for each scope.
func (s *Service) ConfigByScopedGetter(scpGet config.Scoped) (ScopedConfig, error) {

	parent := scpGet.ParentID() // can be website or default
	current := scpGet.ScopeID() // can be store or website or default

	// 99.9999 % of the hits; 2nd argument must be zero because we must first
	// test if a direct entry can be found; if not we must apply either the
	// optionFactory function or do a fall back to the website scope and/or
	// default scope.
	if sCfg, err := s.ConfigByScopeID(current, 0); err == nil {
		if s.Log.IsDebug() {
			s.Log.Debug("ratelimit.Service.ConfigByScopedGetter.IsValid",
				log.Stringer("requested_scope", current),
				log.Stringer("requested_parent_scope", scope.TypeID(0)),
				log.Stringer("responded_scope", sCfg.ScopeID),
			)
		}
		return sCfg, nil
	}

	// load the configuration from the slow backend. optionInflight guarantees
	// that the closure will only be executed once but the returned result gets
	// returned to all waiting goroutines.
	if s.optionFactory != nil {
		res, ok := <-s.optionInflight.DoChan(current.String(), func() (interface{}, error) {
			if err := s.Options(s.optionFactory(scpGet)...); err != nil {
				return ScopedConfig{}, errors.Wrap(err, "[ratelimit] Options applied by OptionFactoryFunc")
			}
			sCfg, err := s.ConfigByScopeID(current, parent)
			if s.Log.IsDebug() {
				s.Log.Debug("ratelimit.Service.ConfigByScopedGetter.Inflight.Do",
					log.ErrWithKey("responded_scope_valid", err),
					log.Stringer("requested_scope", current),
					log.Stringer("requested_parent_scope", parent),
					log.Stringer("responded_scope", sCfg.ScopeID),
					log.Stringer("responded_parent", sCfg.ParentID),
				)
			}
			return sCfg, errors.Wrap(err, "[ratelimit] Options applied by OptionFactoryFunc")
		})
		if !ok { // unlikely to happen but you'll never know. how to test that?
			return ScopedConfig{}, errors.NewFatalf("[ratelimit] Inflight.DoChan returned a closed/unreadable channel")
		}
		if res.Err != nil {
			return ScopedConfig{}, errors.Wrap(res.Err, "[ratelimit] Inflight.DoChan.Error")
		}
		sCfg, ok := res.Val.(ScopedConfig)
		if !ok {
			return ScopedConfig{}, errors.NewFatalf("[ratelimit] Inflight.DoChan res.Val cannot be type asserted to scopedConfig")
		}
		return sCfg, nil
	}

	sCfg, err := s.ConfigByScopeID(current, parent)
	// under very high load: 20 users within 10 MicroSeconds this might get executed
	// 1-3 times. more thinking needed.
	if s.Log.IsDebug() {
		s.Log.Debug("ratelimit.Service.ConfigByScopedGetter.Parent",
			log.Stringer("requested_scope", current),
			log.Stringer("requested_parent_scope", parent),
			log.Stringer("responded_scope", sCfg.ScopeID),
			log.ErrWithKey("responded_scope_valid", err),
		)
	}
	return sCfg, errors.Wrap(err, "[ratelimit] Options applied and finaly validation")
}

// ConfigByScopeID returns the correct configuration for a scope and may fall
// back to the next higher scope: store -> website -> default. If `current`
// TypeID is Store, then the `parent` can only be Website or Default. If an
// entry for a scope cannot be found the next higher scope gets looked up and
// the pointer of the next higher scope gets assigned to the current scope. This
// prevents redundant configurations and enables us to change one scope
// configuration with an impact on all other scopes which depend on the parent
// scope. A zero `parent` triggers no further look ups. This function does not
// load any configuration (config.Getter related) from the backend and accesses
// the internal map of the Service directly.
//
// Important: a "current" scope cannot have multiple "parent" scopes.
func (s *Service) ConfigByScopeID(current scope.TypeID, parent scope.TypeID) (scpCfg ScopedConfig, _ error) {
	// "current" can be Store or Website scope and "parent" can be Website or
	// Default scope. If "parent" equals 0 then no fall back.

	if !current.ValidParent(parent) {
		return scpCfg, errors.NewNotValidf("[ratelimit] The current scope %s has an invalid parent scope %s", current, parent)
	}

	// pointer must get dereferenced in a lock to avoid race conditions while
	// reading in middleware the config values because we might execute the
	// functional options for another scope while one scope runs in the
	// middleware.

	// lookup store/website scope. this should hit 99% of the calls of this function.
	s.rwmu.RLock()
	pScpCfg, ok := s.scopeCache[current]
	if ok && pScpCfg != nil {
		scpCfg = *pScpCfg
	}
	s.rwmu.RUnlock()
	if ok {
		return scpCfg, errors.Wrap(scpCfg.isValid(), "[ratelimit] Validated directly found")
	}
	if parent == 0 {
		return scpCfg, errors.NewNotFoundf(errConfigNotFound, current)
	}

	// slow path: now lock everything until the fall back has been found.
	s.rwmu.Lock()
	defer s.rwmu.Unlock()

	// if the current scope cannot be found, fall back to parent scope and apply
	// the maybe found configuration to the current scope configuration.
	if !ok && parent.Type() == scope.Website {
		pScpCfg, ok = s.scopeCache[parent]
		if ok && pScpCfg != nil {
			pScpCfg.ParentID = parent
			scpCfg = *pScpCfg
			if err := scpCfg.isValid(); err != nil {
				return ScopedConfig{}, errors.Wrap(err, "[ratelimit] Error in Website scope configuration")
			}
			s.scopeCache[current] = pScpCfg // gets assigned a pointer so equal to parent
			return scpCfg, nil
		}
	}

	// if the current and parent scope cannot be found, fall back to default
	// scope and apply the maybe found configuration to the current scope
	// configuration.
	if !ok {
		pScpCfg, ok = s.scopeCache[scope.DefaultTypeID]
		if ok && pScpCfg != nil {
			pScpCfg.ParentID = scope.DefaultTypeID
			scpCfg = *pScpCfg
			if err := scpCfg.isValid(); err != nil {
				return ScopedConfig{}, errors.Wrap(err, "[ratelimit] error in default configuration")
			}
			s.scopeCache[current] = pScpCfg // gets assigned a pointer so equal to default
		} else {
			return scpCfg, errors.NewNotFoundf(errConfigNotFound, scope.DefaultTypeID)
		}
	}
	return scpCfg, nil
}

// findScopedConfig used in functional options to look up if a parent
// configuration exists and if not creates a newScopedConfig(). The
// scope.DefaultTypeID will always be appended to the end of the provided
// arguments. This function acquires a lock. You must call its buddy function
// updateScopedConfig() to close the lock.
func (s *Service) findScopedConfig(scopeIDs ...scope.TypeID) *ScopedConfig {
	s.rwmu.Lock() // Unlock() in updateScopedConfig()

	target, parents := scope.TypeIDs(scopeIDs).TargetAndParents()

	sc := s.scopeCache[target]
	if sc != nil {
		return sc
	}

	// "parents" contains now the next higher scopes, at least minimum the
	// DefaultTypeID. For example if we have as "target" scope Store then
	// "parents" would contain Website and/or Default, depending on how many
	// arguments have been applied in a functional option.
	for _, id := range parents {
		if sc, ok := s.scopeCache[id]; ok && sc != nil {
			shallowCopy := new(ScopedConfig)
			*shallowCopy = *sc
			shallowCopy.ParentID = id
			shallowCopy.ScopeID = target
			return shallowCopy
		}
	}
	// if parents[0] panics for being out of bounds then something is really wrong.
	return newScopedConfig(target, parents[0])
}

// updateScopedConfig used in functional options to store a scoped configuration
// in the internal cache. This function gets called in a function option at the
// end after applying the new configuration value. This function releases an
// already acquired lock. You can call its buddy function findScopedConfig() to
// acquire a lock.
func (s *Service) updateScopedConfig(sc *ScopedConfig) error {
	s.scopeCache[sc.ScopeID] = sc
	s.rwmu.Unlock()
	return nil
}
