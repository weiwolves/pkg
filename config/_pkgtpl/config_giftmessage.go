// +build ignore

package giftmessage

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
			ID: "sales",
			Groups: element.MakeGroups(
				element.Group{
					ID:        "gift_options",
					Label:     `Gift Options`,
					SortOrder: 100,
					Scopes:    scope.PermWebsite,
					Fields: element.MakeFields(
						element.Field{
							// Path: sales/gift_options/allow_order
							ID:        "allow_order",
							Label:     `Allow Gift Messages on Order Level`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: sales/gift_options/allow_items
							ID:        "allow_items",
							Label:     `Allow Gift Messages for Order Items`,
							Type:      element.TypeSelect,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		element.Section{
			ID: "sales",
			Groups: element.MakeGroups(
				element.Group{
					ID: "gift_messages",
					Fields: element.MakeFields(
						element.Field{
							// Path: sales/gift_messages/allow_items
							ID:      `allow_items`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},

						element.Field{
							// Path: sales/gift_messages/allow_order
							ID:      `allow_order`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
