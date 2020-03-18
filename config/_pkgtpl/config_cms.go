// +build ignore

package cms

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
			ID: "web",
			Groups: element.MakeGroups(
				element.Group{
					ID: "default",
					Fields: element.MakeFields(
						element.Field{
							// Path: web/default/cms_home_page
							ID:        "cms_home_page",
							Label:     `CMS Home Page`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `home`,
							// SourceModel: Magento\Cms\Model\Config\Source\Page
						},

						element.Field{
							// Path: web/default/cms_no_route
							ID:        "cms_no_route",
							Label:     `CMS No Route Page`,
							Type:      element.TypeSelect,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `no-route`,
							// SourceModel: Magento\Cms\Model\Config\Source\Page
						},

						element.Field{
							// Path: web/default/cms_no_cookies
							ID:        "cms_no_cookies",
							Label:     `CMS No Cookies Page`,
							Type:      element.TypeSelect,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `enable-cookies`,
							// SourceModel: Magento\Cms\Model\Config\Source\Page
						},

						element.Field{
							// Path: web/default/show_cms_breadcrumbs
							ID:        "show_cms_breadcrumbs",
							Label:     `Show Breadcrumbs for CMS Pages`,
							Type:      element.TypeSelect,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   true,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},

				element.Group{
					ID:        "browser_capabilities",
					Label:     `Browser Capabilities Detection`,
					SortOrder: 200,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: web/browser_capabilities/cookies
							ID:        "cookies",
							Label:     `Redirect to CMS-page if Cookies are Disabled`,
							Type:      element.TypeSelect,
							SortOrder: 100,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: web/browser_capabilities/javascript
							ID:        "javascript",
							Label:     `Show Notice if JavaScript is Disabled`,
							Type:      element.TypeSelect,
							SortOrder: 200,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: web/browser_capabilities/local_storage
							ID:        "local_storage",
							Label:     `Show Notice if Local Storage is Disabled`,
							Type:      element.TypeSelect,
							SortOrder: 300,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},
		element.Section{
			ID:        "cms",
			Label:     `Content Management`,
			SortOrder: 1001,
			Scopes:    scope.PermStore,
			Resource:  0, // Magento_Cms::config_cms
			Groups: element.MakeGroups(
				element.Group{
					ID:        "wysiwyg",
					Label:     `WYSIWYG Options`,
					SortOrder: 100,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: cms/wysiwyg/enabled
							ID:        "enabled",
							Label:     `Enable WYSIWYG Editor`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `enabled`,
							// SourceModel: Magento\Cms\Model\Config\Source\Wysiwyg\Enabled
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		element.Section{
			ID: "web",
			Groups: element.MakeGroups(
				element.Group{
					ID: "default",
					Fields: element.MakeFields(
						element.Field{
							// Path: web/default/front
							ID:      `front`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `cms`,
						},

						element.Field{
							// Path: web/default/no_route
							ID:      `no_route`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `cms/noroute/index`,
						},
					),
				},
			),
		},
		element.Section{
			ID: "system",
			Groups: element.MakeGroups(
				element.Group{
					ID: "media_storage_configuration",
					Fields: element.MakeFields(
						element.Field{
							// Path: system/media_storage_configuration/allowed_resources
							ID:      `allowed_resources`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"wysiwyg_image_folder":"wysiwyg"}`,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
