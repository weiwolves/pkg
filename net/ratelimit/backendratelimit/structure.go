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

package backendratelimit

import (
	"github.com/weiwolves/pkg/config/cfgpath"
	"github.com/weiwolves/pkg/config/element"
	"github.com/weiwolves/pkg/storage/text"
	"github.com/weiwolves/pkg/store/scope"
)

// todo(CS): add config values and path for ratelimit.VaryBy type

// NewConfigStructure global configuration structure for this package. Used in
// frontend (to display the user all the settings) and in backend (scope checks
// and default values). See the source code of this function for the overall
// available sections, groups and fields.
func NewConfigStructure() (element.Sections, error) {
	sortIdx := 10
	var iter = func() int {
		sortIdx += 10
		return sortIdx
	}
	return element.MakeSectionsValidated(
		element.Section{
			ID: cfgpath.MakeRoute("net"),
			Groups: element.MakeGroups(
				element.Group{
					ID:        cfgpath.MakeRoute("ratelimit"),
					Label:     text.Chars(`Rate throtteling`),
					SortOrder: 130,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: net/ratelimit/disabled
							ID:        cfgpath.MakeRoute("disabled"),
							Label:     text.Chars(`Disabled`),
							Comment:   text.Chars(`Set to true to disable rate limiting.`),
							Type:      element.TypeSelect,
							SortOrder: iter(),
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},
						element.Field{
							// Path: net/ratelimit/burst
							ID:        cfgpath.MakeRoute("burst"),
							Label:     text.Chars(`Burst`),
							Comment:   text.Chars(`Defines the number of requests that will be allowed to exceed the rate in a single burst and must be greater than or equal to zero`),
							Type:      element.TypeText,
							SortOrder: iter(),
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   20,
						},
						element.Field{
							// Path: net/ratelimit/requests
							ID:        cfgpath.MakeRoute("requests"),
							Label:     text.Chars(`Requests`),
							Comment:   text.Chars(`Number of requests allowed per time period`),
							Type:      element.TypeText,
							SortOrder: iter(),
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   100,
						},
						element.Field{
							// Path: net/ratelimit/duration
							ID:    cfgpath.MakeRoute("duration"),
							Label: text.Chars(`Duration`),
							Comment: text.Chars(`Per second (s), minute (i), hour (h) or day (d). For example, PerMin(60) permits
60 requests instantly per key followed by one request per second indefinitely
whereas PerSec(1) only permits one request per second with no tolerance for
bursts.`),
							Type:      element.TypeText,
							SortOrder: iter(),
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `h`,
						},
					),
				},
				element.Group{
					ID:        cfgpath.MakeRoute("ratelimit_storage"),
					Label:     text.Chars(`Rate throtteling storage`),
					SortOrder: 140,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: net/ratelimit_storage/gcra_name
							ID:        cfgpath.MakeRoute("gcra_name"),
							Label:     text.Chars(`Name of the registered GCRA`),
							Comment:   text.Chars(`Insert the name of the registered GCRA during program initialization with the function Backend.Register().`),
							Type:      element.TypeText,
							SortOrder: iter(),
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},
						element.Field{
							// Path: net/ratelimit_storage/enable_gcra_memory
							ID:        cfgpath.MakeRoute("enable_gcra_memory"),
							Label:     text.Chars(`Use GCRA in-memory (max keys)`),
							Comment:   text.Chars(`If maxKeys > 0 in-memory key storage will be enabled. The max keys  number of different keys is restricted to the specified amount (65536). In this case, it uses an LRU algorithm to evict older keys to make room for newer ones.`),
							Type:      element.TypeText,
							SortOrder: iter(),
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   0,
						},
						element.Field{
							// Path: net/ratelimit_storage/enable_gcra_redis
							ID:        cfgpath.MakeRoute("enable_gcra_redis"),
							Label:     text.Chars(`Use GCRA Redis`),
							Comment:   text.Chars(`If a Redis URL is provided a Redis server will be used for key storage. Setting both entries (in-memory and Redis) then only Redis will be applied. URLs should follow the draft IANA specification for the scheme (https://www.iana.org/assignments/uri-schemes/prov/redis). For example: redis://localhost:6379/3 |  redis://:6380/0 => connects to localhost:6380 | redis:// => connects to localhost:6379 with DB 0 | redis://empty:myPassword@clusterName.xxxxxx.0001.usw2.cache.amazonaws.com:6379/0`),
							Type:      element.TypeText,
							SortOrder: iter(),
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},
					),
				},
			),
		},
	)
}
