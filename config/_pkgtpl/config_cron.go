// +build ignore

package cron

import (
	"github.com/weiwolves/pkg/config/element"
	"github.com/weiwolves/pkg/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.Sections

func init() {
	ConfigStructure = element.MustMakeSectionsValidate(
		element.Section{
			ID: "system",
			Groups: element.MakeGroups(
				element.Group{
					ID:        "cron",
					Label:     `Cron (Scheduled Tasks) - all the times are in minutes`,
					Comment:   text.Long(`For correct URLs generated during cron runs please make sure that Web > Secure and Unsecure Base URLs are explicitly set.`),
					SortOrder: 15,
					Scopes:    scope.PermDefault,
					Fields:    element.MakeFields(),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
