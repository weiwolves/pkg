// +build ignore

package developer

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
			ID: "dev",
			Groups: element.MakeGroups(
				element.Group{
					ID:        "front_end_development_workflow",
					Label:     `Frontend Development Workflow`,
					SortOrder: 8,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: dev/front_end_development_workflow/type
							ID:        "type",
							Label:     `Workflow type`,
							Comment:   text.Long(`Not available in production mode`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							Default:   `server_side_compilation`,
							// SourceModel: Magento\Developer\Model\Config\Source\WorkflowType
						},
					),
				},

				element.Group{
					ID:        "restrict",
					Label:     `Developer Client Restrictions`,
					SortOrder: 10,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: dev/restrict/allow_ips
							ID:        "allow_ips",
							Label:     `Allowed IPs (comma separated)`,
							Comment:   text.Long(`Leave empty for access from any location.`),
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Developer\Model\Config\Backend\AllowedIps
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		element.Section{
			ID: "dev",
			Groups: element.MakeGroups(
				element.Group{
					ID: "restrict",
					Fields: element.MakeFields(
						element.Field{
							// Path: dev/restrict/allow_ips
							ID:      `allow_ips`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
