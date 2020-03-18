// +build ignore

package multishipping

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
			ID:        "multishipping",
			Label:     `Multishipping Settings`,
			SortOrder: 311,
			Scopes:    scope.PermWebsite,
			Resource:  0, // Magento_Multishipping::config_multishipping
			Groups: element.MakeGroups(
				element.Group{
					ID:        "options",
					Label:     `Options`,
					SortOrder: 2,
					Scopes:    scope.PermWebsite,
					Fields: element.MakeFields(
						element.Field{
							// Path: multishipping/options/checkout_multiple
							ID:        "checkout_multiple",
							Label:     `Allow Shipping to Multiple Addresses`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   true,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: multishipping/options/checkout_multiple_maximum_qty
							ID:        "checkout_multiple_maximum_qty",
							Label:     `Maximum Qty Allowed for Shipping to Multiple Addresses`,
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   100,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
