// +build ignore

package swatches

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
					ID:        "frontend",
					SortOrder: 100,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: catalog/frontend/swatches_per_product
							ID:        "swatches_per_product",
							Label:     `Swatches per Product`,
							Type:      element.TypeText,
							SortOrder: 300,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   16,
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		element.Section{
			ID: "general",
			Groups: element.MakeGroups(
				element.Group{
					ID: "validator_data",
					Fields: element.MakeFields(
						element.Field{
							// Path: general/validator_data/input_types
							ID:      `input_types`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"swatch_visual":"swatch_visual","swatch_text":"swatch_text"}`,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
