// +build ignore

package productvideo

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
			ID: "catalog",
			Groups: element.MakeGroups(
				element.Group{
					ID:        "product_video",
					Label:     `Product Video`,
					SortOrder: 350,
					Scopes:    scope.PermWebsite,
					Fields: element.MakeFields(
						element.Field{
							// Path: catalog/product_video/youtube_api_key
							ID:        "youtube_api_key",
							Label:     `YouTube API Key`,
							Type:      element.TypeText,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},

						element.Field{
							// Path: catalog/product_video/play_if_base
							ID:        "play_if_base",
							Label:     `Autostart base video`,
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: catalog/product_video/show_related
							ID:        "show_related",
							Label:     `Show related video`,
							Type:      element.TypeSelect,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: catalog/product_video/video_auto_restart
							ID:        "video_auto_restart",
							Label:     `Auto restart video`,
							Type:      element.TypeSelect,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
