// +build ignore

package productalert

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
					ID:        "productalert",
					Label:     `Product Alerts`,
					SortOrder: 250,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: catalog/productalert/allow_price
							ID:        "allow_price",
							Label:     `Allow Alert When Product Price Changes`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: catalog/productalert/allow_stock
							ID:        "allow_stock",
							Label:     `Allow Alert When Product Comes Back in Stock`,
							Type:      element.TypeSelect,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: catalog/productalert/email_price_template
							ID:        "email_price_template",
							Label:     `Price Alert Email Template`,
							Comment:   text.Long(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `catalog_productalert_email_price_template`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Template
						},

						element.Field{
							// Path: catalog/productalert/email_stock_template
							ID:        "email_stock_template",
							Label:     `Stock Alert Email Template`,
							Comment:   text.Long(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `catalog_productalert_email_stock_template`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Template
						},

						element.Field{
							// Path: catalog/productalert/email_identity
							ID:        "email_identity",
							Label:     `Alert Email Sender`,
							Type:      element.TypeSelect,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `general`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
						},
					),
				},

				element.Group{
					ID:        "productalert_cron",
					Label:     `Product Alerts Run Settings`,
					SortOrder: 260,
					Scopes:    scope.PermDefault,
					Fields: element.MakeFields(
						element.Field{
							// Path: catalog/productalert_cron/frequency
							ID:        "frequency",
							Label:     `Frequency`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// BackendModel: Magento\Cron\Model\Config\Backend\Product\Alert
							// SourceModel: Magento\Cron\Model\Config\Source\Frequency
						},

						element.Field{
							// Path: catalog/productalert_cron/time
							ID:        "time",
							Label:     `Start Time`,
							Type:      element.TypeTime,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
						},

						element.Field{
							// Path: catalog/productalert_cron/error_email
							ID:        "error_email",
							Label:     `Error Email Recipient`,
							Type:      element.TypeText,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
						},

						element.Field{
							// Path: catalog/productalert_cron/error_email_identity
							ID:        "error_email_identity",
							Label:     `Error Email Sender`,
							Type:      element.TypeSelect,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							Default:   `general`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
						},

						element.Field{
							// Path: catalog/productalert_cron/error_email_template
							ID:        "error_email_template",
							Label:     `Error Email Template`,
							Comment:   text.Long(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							Default:   `catalog_productalert_cron_error_email_template`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Template
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		element.Section{
			ID: "catalog",
			Groups: element.MakeGroups(
				element.Group{
					ID: "productalert_cron",
					Fields: element.MakeFields(
						element.Field{
							// Path: catalog/productalert_cron/error_email
							ID:      `error_email`,
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
