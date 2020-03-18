// +build ignore

package wishlist

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
			ID:        "wishlist",
			Label:     `Wish List`,
			SortOrder: 140,
			Scopes:    scope.PermStore,
			Resource:  0, // Magento_Wishlist::config_wishlist
			Groups: element.MakeGroups(
				element.Group{
					ID:        "email",
					Label:     `Share Options`,
					SortOrder: 2,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: wishlist/email/email_identity
							ID:        "email_identity",
							Label:     `Email Sender`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `general`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
						},

						element.Field{
							// Path: wishlist/email/email_template
							ID:        "email_template",
							Label:     `Email Template`,
							Comment:   text.Long(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `wishlist_email_email_template`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Template
						},

						element.Field{
							// Path: wishlist/email/number_limit
							ID:        "number_limit",
							Label:     `Max Emails Allowed to be Sent`,
							Comment:   text.Long(`10 by default. Max - 10000`),
							Type:      element.TypeText,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   10,
						},

						element.Field{
							// Path: wishlist/email/text_limit
							ID:        "text_limit",
							Label:     `Email Text Length Limit`,
							Comment:   text.Long(`255 by default`),
							Type:      element.TypeText,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   255,
						},
					),
				},

				element.Group{
					ID:        "general",
					Label:     `General Options`,
					SortOrder: 1,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: wishlist/general/active
							ID:        "active",
							Label:     `Enabled`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   true,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},

				element.Group{
					ID:        "wishlist_link",
					Label:     `My Wish List Link`,
					SortOrder: 3,
					Scopes:    scope.PermWebsite,
					Fields: element.MakeFields(
						element.Field{
							// Path: wishlist/wishlist_link/use_qty
							ID:        "use_qty",
							Label:     `Display Wish List Summary`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							// SourceModel: Magento\Wishlist\Model\Config\Source\Summary
						},
					),
				},
			),
		},
		element.Section{
			ID: "rss",
			Groups: element.MakeGroups(
				element.Group{
					ID:        "wishlist",
					Label:     `Wish List`,
					SortOrder: 2,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: rss/wishlist/active
							ID:        "active",
							Label:     `Enable RSS`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
