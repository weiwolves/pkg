// +build ignore

package shipping

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
			ID:        "shipping",
			Label:     `Shipping Settings`,
			SortOrder: 310,
			Scopes:    scope.PermStore,
			Resource:  0, // Magento_Shipping::config_shipping
			Groups: element.MakeGroups(
				element.Group{
					ID:        "origin",
					Label:     `Origin`,
					SortOrder: 1,
					Scopes:    scope.PermWebsite,
					Fields: element.MakeFields(
						element.Field{
							// Path: shipping/origin/country_id
							ID:        "country_id",
							Label:     `Country`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `US`,
							// SourceModel: Magento\Directory\Model\Config\Source\Country
						},

						element.Field{
							// Path: shipping/origin/region_id
							ID:        "region_id",
							Label:     `Region/State`,
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   12,
						},

						element.Field{
							// Path: shipping/origin/postcode
							ID:        "postcode",
							Label:     `ZIP/Postal Code`,
							Type:      element.TypeText,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   90034,
						},

						element.Field{
							// Path: shipping/origin/city
							ID:        "city",
							Label:     `City`,
							Type:      element.TypeText,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},

						element.Field{
							// Path: shipping/origin/street_line1
							ID:        "street_line1",
							Label:     `Street Address`,
							Type:      element.TypeText,
							SortOrder: 50,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},

						element.Field{
							// Path: shipping/origin/street_line2
							ID:        "street_line2",
							Label:     `Street Address Line 2`,
							Type:      element.TypeText,
							SortOrder: 60,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},
					),
				},

				element.Group{
					ID:        "shipping_policy",
					Label:     `Shipping Policy Parameters`,
					SortOrder: 120,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: shipping/shipping_policy/enable_shipping_policy
							ID:        "enable_shipping_policy",
							Label:     `Apply custom Shipping Policy`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: shipping/shipping_policy/shipping_policy_content
							ID:        "shipping_policy_content",
							Label:     `Shipping Policy`,
							Type:      element.TypeTextarea,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},
					),
				},
			),
		},
		element.Section{
			ID:        "carriers",
			Label:     `Shipping Methods`,
			SortOrder: 320,
			Scopes:    scope.PermStore,
			Resource:  0, // Magento_Shipping::carriers
			Groups:    element.MakeGroups(),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
