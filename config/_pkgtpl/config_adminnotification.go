// +build ignore

package adminnotification

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
					ID:        "adminnotification",
					Label:     `Notifications`,
					SortOrder: 250,
					Scopes:    scope.PermDefault,
					Fields: element.MakeFields(
						element.Field{
							// Path: system/adminnotification/use_https
							ID:        "use_https",
							Label:     `Use HTTPS to Get Feed`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: system/adminnotification/frequency
							ID:        "frequency",
							Label:     `Update Frequency`,
							Type:      element.TypeSelect,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							Default:   true,
							// SourceModel: Magento\AdminNotification\Model\Config\Source\Frequency
						},

						element.Field{
							// Path: system/adminnotification/last_update
							ID:        "last_update",
							Label:     `Last Update`,
							Type:      element.TypeLabel,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		element.Section{
			ID: "system",
			Groups: element.MakeGroups(
				element.Group{
					ID: "adminnotification",
					Fields: element.MakeFields(
						element.Field{
							// Path: system/adminnotification/feed_url
							ID:      `feed_url`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `notifications.magentocommerce.com/magento2/community/notifications.rss`,
						},

						element.Field{
							// Path: system/adminnotification/popup_url
							ID:      `popup_url`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `widgets.magentocommerce.com/notificationPopup`,
						},

						element.Field{
							// Path: system/adminnotification/severity_icons_url
							ID:      `severity_icons_url`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `widgets.magentocommerce.com/%s/%s.gif`,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
