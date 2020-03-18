// +build ignore

package layerednavigation

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
			ID:        "catalog",
			SortOrder: 40,
			Scopes:    scope.PermStore,
			Groups: element.MakeGroups(
				element.Group{
					ID:        "layered_navigation",
					Label:     `Layered Navigation`,
					SortOrder: 490,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: catalog/layered_navigation/display_product_count
							ID:        "display_product_count",
							Label:     `Display Product Count`,
							Type:      element.TypeSelect,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   true,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: catalog/layered_navigation/price_range_calculation
							ID:        "price_range_calculation",
							Label:     `Price Navigation Step Calculation`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `auto`,
							// SourceModel: Magento\Catalog\Model\Config\Source\Price\Step
						},

						element.Field{
							// Path: catalog/layered_navigation/price_range_step
							ID:        "price_range_step",
							Label:     `Default Price Navigation Step`,
							Type:      element.TypeText,
							SortOrder: 15,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   100,
						},

						element.Field{
							// Path: catalog/layered_navigation/price_range_max_intervals
							ID:        "price_range_max_intervals",
							Label:     `Maximum Number of Price Intervals`,
							Comment:   text.Long(`Maximum number of price intervals is 100`),
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   10,
						},

						element.Field{
							// Path: catalog/layered_navigation/one_price_interval
							ID:        "one_price_interval",
							Label:     `Display Price Interval as One Price`,
							Comment:   text.Long(`This setting will be applied when all prices in the specific price interval are equal.`),
							Type:      element.TypeSelect,
							SortOrder: 15,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: catalog/layered_navigation/interval_division_limit
							ID:        "interval_division_limit",
							Label:     `Interval Division Limit`,
							Comment:   text.Long(`Please specify the number of products, that will not be divided into subintervals.`),
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   9,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
