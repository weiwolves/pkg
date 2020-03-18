// +build ignore

package translation

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
					ID: "js",
					Fields: element.MakeFields(
						element.Field{
							// Path: dev/js/translate_strategy
							ID:        "translate_strategy",
							Label:     `Translation Strategy`,
							Comment:   text.Long(`Please put your store into maintenance mode and redeploy static files after changing strategy`),
							Type:      element.TypeSelect,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							Default:   `dictionary`,
							// SourceModel: Magento\Translation\Model\Js\Config\Source\Strategy
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
					ID: "translate_inline",
					Fields: element.MakeFields(
						element.Field{
							// Path: dev/translate_inline/active
							ID:      `active`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},

						element.Field{
							// Path: dev/translate_inline/active_admin
							ID:      `active_admin`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},

						element.Field{
							// Path: dev/translate_inline/invalid_caches
							ID:      `invalid_caches`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"block_html":null}`,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
