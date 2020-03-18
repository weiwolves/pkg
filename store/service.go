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

package store

import (
	"sort"
	"sync"
	"sync/atomic"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/weiwolves/pkg/store/scope"
)

//go:generate go run codegen_main.go

// Finder depends on the runMode from package scope and finds the active store
// depending on the run mode. The Hash argument will be provided via
// scope.RunMode type or the scope.FromContextRunMode(ctx) function. runMode is
// named in Mage world: MAGE_RUN_CODE and MAGE_RUN_TYPE. The MAGE_RUN_TYPE can
// be either website or store scope and MAGE_RUN_CODE any defined website or
// store code from the database. In our case we must pass an ID and not a code
// string.
type Finder interface {
	// DefaultStoreID returns the default active store ID and its website ID
	// depending on the run mode. Error behaviour is mostly of type NotValid.
	DefaultStoreID(runMode scope.TypeID) (websiteID, storeID uint32, err error)
	// StoreIDbyCode returns, depending on the runMode, for a storeCode its
	// active store ID and its website ID. An empty runMode hash falls back to
	// select the default website, with its default group, and the slice of
	// default stores. A not-found error behaviour gets returned if the code
	// cannot be found. If the runMode equals to scope.DefaultTypeID, the
	// returned ID is always 0 and error is nil.
	StoreIDbyCode(runMode scope.TypeID, storeCode string) (websiteID, storeID uint32, err error)
}

// Service represents type which handles the underlying storage and takes care
// of the default stores. A Service is bound a specific scope.Scope. Depending
// on the scope it is possible or not to switch stores. A Service contains also
// a config.Getter which gets passed to the scope of a Store(), Group() or
// Website() so that you always have the possibility to access a scoped based
// configuration value. This Service uses three internal maps to cache Websites,
// Groups and Stores.

type Service struct {
	// defaultStore someone must be always the default guy. Handled via atomic
	// package.
	defaultStoreID      int64
	chanClose           chan struct{}
	chanEventSubscriber []chan int
	log                 log.Logger

	// mu protects the following fields ... maybe use more mutexes
	mu sync.RWMutex
	// in general these caches can be optimized
	websites StoreWebsites
	groups   StoreGroups
	stores   Stores

	// uint32 key identifies a website, group or store
	cacheWebsite map[uint32]*StoreWebsite
	cacheGroup   map[uint32]*StoreGroup
	cacheStore   map[uint32]*Store
}

// Option type to pass options to the service type.
type Option struct {
	sortOrder int
	fn        func(*Service) error
}

const (
	eventOptionsApplied = iota + 1
	eventClose
)

// NewService creates a new store Service which handles websites, groups and
// stores. You must either provide the functional options or call LoadFromDB()
// to setup the internal cache.
func NewService(opts ...Option) (*Service, error) {
	srv := &Service{
		chanClose:      make(chan struct{}),
		defaultStoreID: -1, // means not set, because 0 can be admin store.
		cacheWebsite:   make(map[uint32]*StoreWebsite),
		cacheGroup:     make(map[uint32]*StoreGroup),
		cacheStore:     make(map[uint32]*Store),
	}
	if err := srv.Options(opts...); err != nil {
		return nil, errors.WithStack(err)
	}
	return srv, nil
}

// MustNewService same as NewService, but panics on error.
func MustNewService(opts ...Option) *Service {
	m, err := NewService(opts...)
	if err != nil {
		panic(err)
	}
	return m
}

// Options applies various options to the running store service.
func (s *Service) Options(opts ...Option) error {
	sort.Slice(opts, func(i, j int) bool {
		return opts[i].sortOrder < opts[j].sortOrder // ascending 0-9 sorting ;-)
	})
	for _, opt := range opts {
		if err := opt.fn(s); err != nil {
			return errors.WithStack(err)
		}
	}
	s.sort()
	s.apply2ndLevelData()
	if err := s.validate(); err != nil {
		return errors.WithStack(err)
	}
	s.dispatchEvent(eventOptionsApplied)
	return nil
}

func (s *Service) sort() {
	s.mu.Lock()
	sort.Stable(sortNaturallyWebsites{&s.websites})
	sort.Stable(sortNaturallyGroups{&s.groups})
	sort.Stable(sortNaturallyStores{&s.stores})
	s.mu.Unlock()
}

func (s *Service) apply2ndLevelData() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Websites
	for _, w := range s.websites.Data {
		// Stores
		if w.Stores == nil {
			w.Stores = NewStores()
		} else {
			for i := range w.Stores.Data {
				w.Stores.Data[i] = nil
			}
			w.Stores.Data = w.Stores.Data[:0]
		}
		for _, st := range s.stores.Data {
			if w.WebsiteID == st.WebsiteID {
				st2 := st.Copy()
				st2.StoreWebsite = nil
				st2.StoreGroup = nil
				w.Stores.Append(st2)
			}
		}
		// Groups
		if w.StoreGroups == nil {
			w.StoreGroups = NewStoreGroups()
		} else {
			for i := range w.StoreGroups.Data {
				w.StoreGroups.Data[i] = nil
			}
			w.StoreGroups.Data = w.StoreGroups.Data[:0]
		}
		for _, g := range s.groups.Data {
			if w.WebsiteID == g.WebsiteID {
				g2 := g.Copy()
				g2.StoreWebsite = nil
				w.StoreGroups.Append(g2)
			}
		}
	}

	// Groups
	for _, g := range s.groups.Data {
		for _, w := range s.websites.Data {
			if w.WebsiteID == g.WebsiteID {
				w2 := w.Copy()
				w2.Stores = nil
				w2.StoreGroups = nil
				g.StoreWebsite = w2
			}
		}
	}
	// Stores
	for _, st := range s.stores.Data {
		for _, g := range s.groups.Data {
			if st.GroupID == g.GroupID {
				g2 := g.Copy()
				g2.StoreWebsite = nil
				st.StoreGroup = g2
			}
		}
		for _, w := range s.websites.Data {
			if st.WebsiteID == w.WebsiteID {
				w2 := w.Copy()
				w2.Stores = nil
				w2.StoreGroups = nil
				st.StoreWebsite = w2
			}
		}
	}
}

func (s *Service) validate() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// These checks are usually database constraints and logic checks

	// Website: Check for default website and DefaultGroupIDs
	var defaultWebsiteCounter int
	for _, w := range s.websites.Data {
		if w.IsDefault {
			defaultWebsiteCounter++
		}

		var foundGroupID bool
		for _, g := range s.groups.Data {
			if w.DefaultGroupID == g.GroupID && !foundGroupID {
				foundGroupID = true
			}
		}
		if !foundGroupID {
			return errors.NotValid.Newf("[store] Website[%d].DefaultGroupID[%d] not found in Groups.", w.WebsiteID, w.DefaultGroupID)
		}
	} // end each websites
	if len(s.websites.Data) > 0 && defaultWebsiteCounter > 1 {
		return errors.NotValid.Newf("[store] Only one Website can be the default Website. Found %d default websites", defaultWebsiteCounter)
	}

	// Groups: Check for website ID and default store ID
	for _, g := range s.groups.Data {
		var foundWebsiteID bool
		for _, w := range s.websites.Data {
			if g.WebsiteID == w.WebsiteID && !foundWebsiteID {
				foundWebsiteID = true
			}
		}
		if !foundWebsiteID {
			return errors.NotValid.Newf("[store] Group[%d].WebsiteID[%d] not found in Websites.", g.GroupID, g.WebsiteID)
		}
		var foundStoreID bool
		for _, s := range s.stores.Data {
			if g.DefaultStoreID == s.StoreID && !foundStoreID {
				foundStoreID = true
			}
		}
		if !foundStoreID {
			return errors.NotValid.Newf("[store] Group[%d].DefaultStoreID[%d] not found in Websites.", g.GroupID, g.DefaultStoreID)
		}
	}

	for _, st := range s.stores.Data {
		var foundWebsiteID bool
		for _, w := range s.websites.Data {
			if st.WebsiteID == w.WebsiteID && !foundWebsiteID {
				foundWebsiteID = true
			}
		}
		if !foundWebsiteID {
			return errors.NotValid.Newf("[store] Store[%d].WebsiteID[%d] not found in Websites.", st.StoreID, st.WebsiteID)
		}
		var foundGroupID bool
		for _, g := range s.groups.Data {
			if st.GroupID == g.GroupID && !foundGroupID {
				foundGroupID = true
			}
		}
		if !foundGroupID {
			return errors.NotValid.Newf("[store] Store[%d].GroupID[%d] not found in Group.", st.StoreID, st.GroupID)
		}
	}

	return nil
}

func (s *Service) dispatchEvent(id int) {
	for _, ces := range s.chanEventSubscriber {
		ces <- id
	}
}

func (s *Service) Close() error {
	s.dispatchEvent(eventClose)
	close(s.chanClose)
	return nil
}

// IsAllowedStoreID checks if the store ID is allowed within the runMode.
// Returns true on success. An error may occur when the default website and
// store can't be selected. An empty scope.Hash checks the default website with
// its default group and its default stores.
func (s *Service) IsAllowedStoreID(runMode scope.TypeID, storeID uint32) (isAllowed bool, storeCode string, _ error) {
	scp, scpID := runMode.Unpack()

	switch scp {
	case scope.Store:
		st, err := s.Store(storeID)
		if err != nil {
			return false, "", errors.WithStack(err)
		}
		return st.IsActive && st.Code != "", st.Code, nil

	case scope.Group:
		for _, st := range s.stores.Data {
			if st.IsActive && st.GroupID == scpID && st.StoreID == storeID && st.Code != "" {
				return true, st.Code, nil
			}
		}
		return false, "", nil
	case scope.Website:
		for _, st := range s.stores.Data {
			if st.IsActive && st.WebsiteID == scpID && st.StoreID == storeID && st.Code != "" {
				return true, st.Code, nil
			}
		}
		return false, "", nil
	default:
		w, err := s.websites.Default()
		if err != nil {
			return false, "", errors.WithStack(err)
		}
		g, err := w.DefaultGroup()
		if err != nil {
			return false, "", errors.WithStack(err)
		}
		for _, st := range s.stores.Data {
			if st.IsActive && st.WebsiteID == w.WebsiteID && st.GroupID == g.GroupID && st.StoreID == storeID {
				return true, st.Code, nil
			}
		}
		return false, "", nil
	}
}

// DefaultStoreView returns the overall default store view depending on the
// website, group and store settings. It traverses websites, looks for
// is_default, returns the website.default_group_id, then group with
// default_store_id and finally loads the active store. If not an active store
// can be found and NotFound error gets returned.
func (s *Service) DefaultStoreView() (*Store, error) {
	s.mu.RLock()

	dsID := atomic.LoadInt64(&s.defaultStoreID)
	if dsID >= 0 {
		cs, ok := s.cacheStore[uint32(dsID)]
		s.mu.RUnlock()
		if ok {
			return cs, nil
		}
	}

	var id int64 = -1
WebsiteLoop:
	for _, w := range s.websites.Data {
		if w.IsDefault {
			for _, g := range s.groups.Data {
				if g.GroupID == w.DefaultGroupID {
					for _, t := range s.stores.Data {
						if g.DefaultStoreID == t.StoreID && t.IsActive {
							id = int64(t.StoreID)
							break WebsiteLoop
						}
					}
				}
			}
		}
	}
	s.mu.RUnlock()
	if id < 0 {
		return nil, errors.NotFound.Newf("[store] Default Store ID not found.")
	}

	atomic.StoreInt64(&s.defaultStoreID, id)
	return s.Store(uint32(id))
}

// DefaultStoreID returns the default active store ID depending on the run mode.
// Error behaviour is mostly of type NotValid.
func (s *Service) DefaultStoreID(runMode scope.TypeID) (websiteID, storeID uint32, _ error) {
	switch scp, id := runMode.Unpack(); scp {
	case scope.Store:
		st, err := s.Store(id)
		if err != nil {
			return 0, 0, errors.Wrapf(err, "[store] DefaultStoreID Scope %s ID %d", scp, id)
		}
		if !st.IsActive {
			return 0, 0, errors.NotValid.Newf("[store] DefaultStoreID %s the store ID %d is not active", runMode, st.StoreID)
		}
		return st.WebsiteID, st.StoreID, nil

	case scope.Group:
		g, err := s.Group(id)
		if err != nil {
			return 0, 0, errors.Wrapf(err, "[store] DefaultStoreID Scope %s ID %d", scp, id)
		}
		st, err := s.Store(g.DefaultStoreID)
		if err != nil {
			return 0, 0, errors.Wrapf(err, "[store] DefaultStoreID Scope %s ID %d", scp, id)
		}
		if !st.IsActive {
			return 0, 0, errors.NotValid.Newf("[store] DefaultStoreID %s the store ID %d is not active", runMode, st.StoreID)
		}
		return st.WebsiteID, st.StoreID, nil

	case scope.Website:
		w, err := s.Website(id)
		if err != nil {
			return 0, 0, errors.Wrapf(err, "[store] DefaultStoreID.Website Scope %s ID %d", scp, id)
		}

		// Special case for admin scope, all zero
		if w.WebsiteID == 0 && w.DefaultGroupID == 0 {
			st, err := s.Store(0)
			if err != nil {
				return 0, 0, errors.Wrapf(err, "[store] DefaultStoreID.Website.Store Scope %s ID %d", scp, id)
			}
			return st.WebsiteID, st.StoreID, nil
		}

		st, err := w.DefaultStore()
		if err != nil {
			return 0, 0, errors.Wrapf(err, "[store] DefaultStoreID.Website.DefaultStore Scope %s ID %d", scp, id)
		}
		return st.WebsiteID, st.StoreID, nil

	default:
		w, err := s.websites.Default()
		if err != nil {
			return 0, 0, errors.Wrapf(err, "[store] DefaultStoreID.Website.Default Scope %s ID %d", scp, id)
		}
		st, err := w.DefaultStore()
		if err != nil {
			return 0, 0, errors.Wrapf(err, "[store] DefaultStoreID.Website.DefaultStore Scope %s ID %d", scp, id)
		}
		return st.WebsiteID, st.StoreID, nil
	}
}

// StoreIDbyCode returns, depending on the runMode, for a storeCode its active
// store ID and its website ID. An empty runMode hash falls back to select the
// default website, with its default group, and the slice of default stores. A
// not-found error behaviour gets returned if the code cannot be found. If the
// runMode equals to scope.DefaultTypeID, the returned ID is always 0 and error
// is nil. Implements interface Finder.
func (s *Service) StoreIDbyCode(runMode scope.TypeID, storeCode string) (websiteID, storeID uint32, err error) {
	if storeCode == "" {
		wID, sID, err := s.DefaultStoreID(0)
		return wID, sID, errors.WithStack(err)
	}

	runModeID, err := runMode.ID()
	if err != nil {
		return 0, 0, errors.WithStack(err)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	switch runMode.Type() {
	case scope.Store:
		for _, st := range s.stores.Data {
			if st.IsActive && st.Code == storeCode {
				return st.WebsiteID, st.StoreID, nil
			}
		}
	case scope.Group:
		for _, st := range s.stores.Data {
			if st.IsActive && st.GroupID == runModeID && st.Code == storeCode {
				return st.WebsiteID, st.StoreID, nil
			}
		}
	case scope.Website:
		for _, st := range s.stores.Data {
			if st.IsActive && st.WebsiteID == runModeID && st.Code == storeCode {
				return st.WebsiteID, st.StoreID, nil
			}
		}
	default:
		w, err := s.websites.Default()
		if err != nil {
			return 0, 0, errors.WithStack(err)
		}
		st, err := w.DefaultStore()
		if err != nil {
			return 0, 0, errors.WithStack(err)
		}
		if st.Code != "" && st.Code == storeCode {
			return st.WebsiteID, st.StoreID, nil
		}
	}
	return 0, 0, errors.NotFound.Newf("[store] Code %q not found for runMode %s", storeCode, runMode)
}

// AllowedStores creates a new slice containing all active stores depending on
// the current runMode. The returned slice and its pointers are owned by the
// callee.
func (s *Service) AllowedStores(runMode scope.TypeID) (*Stores, error) {
	scp, scpID := runMode.Unpack()

	// copy because Filter function reuses backing slice array
	storeCol := &Stores{
		Data: make([]*Store, len(s.stores.Data)),
	}
	copy(storeCol.Data, s.stores.Data)

	switch scp {
	case scope.Default:
		return storeCol.Filter(func(st *Store) bool {
			return st.IsActive && st.StoreID == 0
		}), nil

	case scope.Store:
		return storeCol.Filter(func(st *Store) bool {
			return st.IsActive
		}), nil

	case scope.Group:
		return storeCol.Filter(func(st *Store) bool {
			return st.IsActive && st.GroupID == scpID
		}), nil

	case scope.Website:
		return storeCol.Filter(func(st *Store) bool {
			return st.IsActive && st.WebsiteID == scpID
		}), nil

	default:
		return nil, errors.NotImplemented.Newf("[store] Scope %s not yet implemented.", scp)
	}
}

// Website returns the cached Website from an ID including all of its groups and
// all related stores.
// The returned pointer is owned by the Service. You must copy it for further modifications.
func (s *Service) Website(id uint32) (*StoreWebsite, error) {
	s.mu.RLock()
	w, ok := s.cacheWebsite[id]
	s.mu.RUnlock()
	if ok {
		return w, nil
	}
	for _, w := range s.websites.Data {
		if w.WebsiteID == id {
			s.mu.Lock()
			s.cacheWebsite[id] = w
			s.mu.Unlock()
			return w, nil
		}
	}
	return nil, errors.NotFound.Newf("[store] Cannot find Website ID %d", id)
}

// Websites returns a cached slice containing all Websites with its associated
// groups and stores. You shall not modify the returned slice.
func (s *Service) Websites() StoreWebsites {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.websites
}

// Group returns a cached Group which contains all related stores and its website.
// The returned pointer is owned by the Service. You must copy it for further modifications.
func (s *Service) Group(id uint32) (*StoreGroup, error) {
	s.mu.RLock()
	g, ok := s.cacheGroup[id]
	s.mu.RUnlock()
	if ok {
		return g, nil
	}
	for _, g := range s.groups.Data {
		if g.GroupID == id {
			s.mu.Lock()
			s.cacheGroup[id] = g
			s.mu.Unlock()
			return g, nil
		}
	}
	return nil, errors.NotFound.Newf("[store] Cannot find Group ID %d", id)
}

// Groups returns a cached slice containing all  Groups with its associated
// stores and websites. You shall not modify the returned slice.
func (s *Service) Groups() StoreGroups {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.groups
}

// Store returns the cached Store view containing its group and its website.
// The returned pointer is owned by the Service. You must copy it for further modifications.
func (s *Service) Store(id uint32) (*Store, error) {
	s.mu.RLock()
	d, ok := s.cacheStore[id]
	s.mu.RUnlock()
	if ok {
		return d, nil
	}
	for _, d := range s.stores.Data {
		if d.StoreID == id {
			s.mu.Lock()
			s.cacheStore[id] = d
			s.mu.Unlock()
			return d, nil
		}
	}
	return nil, errors.NotFound.Newf("[store] Cannot find Store ID %d", id)
}

// Stores returns a cached Store slice containing all related websites and groups.
// You shall not modify the returned slice.
func (s *Service) Stores() Stores {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.stores
}

// ClearCache resets the internal caches which stores the pointers to Websites,
// Groups or Stores. The ReInit() also uses this method to clear caches before
// the Storage gets reloaded.
func (s *Service) ClearCache() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.cacheWebsite) > 0 {
		for k := range s.cacheWebsite {
			delete(s.cacheWebsite, k)
		}
	}
	if len(s.cacheGroup) > 0 {
		for k := range s.cacheGroup {
			delete(s.cacheGroup, k)
		}
	}
	if len(s.cacheStore) > 0 {
		for k := range s.cacheStore {
			delete(s.cacheStore, k)
		}
	}
	s.defaultStoreID = -1

	for i := range s.websites.Data {
		s.websites.Data[i] = nil
	}
	s.websites.Data = s.websites.Data[:0]

	for i := range s.groups.Data {
		s.groups.Data[i] = nil
	}
	s.groups.Data = s.groups.Data[:0]

	for i := range s.stores.Data {
		s.stores.Data[i] = nil
	}
	s.stores.Data = s.stores.Data[:0]
}

// IsCacheEmpty returns true if the internal cache is empty.
func (s *Service) IsCacheEmpty() bool {
	return len(s.cacheWebsite) == 0 && len(s.cacheGroup) == 0 && len(s.cacheStore) == 0 &&
		s.defaultStoreID == -1
}
