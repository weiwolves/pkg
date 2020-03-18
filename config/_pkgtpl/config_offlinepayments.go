// +build ignore

package offlinepayments

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
			ID:        "payment",
			SortOrder: 400,
			Scopes:    scope.PermStore,
			Groups: element.MakeGroups(
				element.Group{
					ID:        "checkmo",
					Label:     `Check / Money Order`,
					SortOrder: 30,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: payment/checkmo/active
							ID:        "active",
							Label:     `Enabled`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   true,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: payment/checkmo/order_status
							ID:        "order_status",
							Label:     `New Order Status`,
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `pending`,
							// SourceModel: Magento\Sales\Model\Config\Source\Order\Status\NewStatus
						},

						element.Field{
							// Path: payment/checkmo/sort_order
							ID:        "sort_order",
							Label:     `Sort Order`,
							Type:      element.TypeText,
							SortOrder: 100,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},

						element.Field{
							// Path: payment/checkmo/title
							ID:        "title",
							Label:     `Title`,
							Type:      element.TypeText,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `Check / Money order`,
						},

						element.Field{
							// Path: payment/checkmo/allowspecific
							ID:        "allowspecific",
							Label:     `Payment from Applicable Countries`,
							Type:      element.TypeAllowspecific,
							SortOrder: 50,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							// SourceModel: Magento\Payment\Model\Config\Source\Allspecificcountries
						},

						element.Field{
							// Path: payment/checkmo/specificcountry
							ID:         "specificcountry",
							Label:      `Payment from Specific Countries`,
							Type:       element.TypeMultiselect,
							SortOrder:  51,
							Visible:    element.VisibleYes,
							Scopes:     scope.PermWebsite,
							CanBeEmpty: true,
							// SourceModel: Magento\Directory\Model\Config\Source\Country
						},

						element.Field{
							// Path: payment/checkmo/payable_to
							ID:        "payable_to",
							Label:     `Make Check Payable to`,
							Type:      element.Type,
							SortOrder: 61,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},

						element.Field{
							// Path: payment/checkmo/mailing_address
							ID:        "mailing_address",
							Label:     `Send Check to`,
							Type:      element.TypeTextarea,
							SortOrder: 62,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},

						element.Field{
							// Path: payment/checkmo/min_order_total
							ID:        "min_order_total",
							Label:     `Minimum Order Total`,
							Type:      element.TypeText,
							SortOrder: 98,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},

						element.Field{
							// Path: payment/checkmo/max_order_total
							ID:        "max_order_total",
							Label:     `Maximum Order Total`,
							Type:      element.TypeText,
							SortOrder: 99,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},

						element.Field{
							// Path: payment/checkmo/cfgmodel
							ID:      "model",
							Type:    element.Type,
							Visible: element.VisibleYes,
							Default: `Magento\OfflinePayments\Model\Checkmo`,
						},
					),
				},

				element.Group{
					ID:        "purchaseorder",
					Label:     `Purchase Order`,
					SortOrder: 32,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: payment/purchaseorder/active
							ID:        "active",
							Label:     `Enabled`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: payment/purchaseorder/order_status
							ID:        "order_status",
							Label:     `New Order Status`,
							Type:      element.TypeSelect,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `pending`,
							// SourceModel: Magento\Sales\Model\Config\Source\Order\Status\NewStatus
						},

						element.Field{
							// Path: payment/purchaseorder/sort_order
							ID:        "sort_order",
							Label:     `Sort Order`,
							Type:      element.TypeText,
							SortOrder: 100,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},

						element.Field{
							// Path: payment/purchaseorder/title
							ID:        "title",
							Label:     `Title`,
							Type:      element.TypeText,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `Purchase Order`,
						},

						element.Field{
							// Path: payment/purchaseorder/allowspecific
							ID:        "allowspecific",
							Label:     `Payment from Applicable Countries`,
							Type:      element.TypeAllowspecific,
							SortOrder: 50,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							// SourceModel: Magento\Payment\Model\Config\Source\Allspecificcountries
						},

						element.Field{
							// Path: payment/purchaseorder/specificcountry
							ID:         "specificcountry",
							Label:      `Payment from Specific Countries`,
							Type:       element.TypeMultiselect,
							SortOrder:  51,
							Visible:    element.VisibleYes,
							Scopes:     scope.PermWebsite,
							CanBeEmpty: true,
							// SourceModel: Magento\Directory\Model\Config\Source\Country
						},

						element.Field{
							// Path: payment/purchaseorder/min_order_total
							ID:        "min_order_total",
							Label:     `Minimum Order Total`,
							Type:      element.TypeText,
							SortOrder: 98,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},

						element.Field{
							// Path: payment/purchaseorder/max_order_total
							ID:        "max_order_total",
							Label:     `Maximum Order Total`,
							Type:      element.TypeText,
							SortOrder: 99,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},

						element.Field{
							// Path: payment/purchaseorder/cfgmodel
							ID:      "model",
							Type:    element.Type,
							Visible: element.VisibleYes,
							Default: `Magento\OfflinePayments\Model\Purchaseorder`,
						},
					),
				},

				element.Group{
					ID:        "banktransfer",
					Label:     `Bank Transfer Payment`,
					SortOrder: 30,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: payment/banktransfer/active
							ID:        "active",
							Label:     `Enabled`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: payment/banktransfer/title
							ID:        "title",
							Label:     `Title`,
							Type:      element.TypeText,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `Bank Transfer Payment`,
						},

						element.Field{
							// Path: payment/banktransfer/order_status
							ID:        "order_status",
							Label:     `New Order Status`,
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `pending`,
							// SourceModel: Magento\Sales\Model\Config\Source\Order\Status\NewStatus
						},

						element.Field{
							// Path: payment/banktransfer/allowspecific
							ID:        "allowspecific",
							Label:     `Payment from Applicable Countries`,
							Type:      element.TypeAllowspecific,
							SortOrder: 50,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							// SourceModel: Magento\Payment\Model\Config\Source\Allspecificcountries
						},

						element.Field{
							// Path: payment/banktransfer/specificcountry
							ID:         "specificcountry",
							Label:      `Payment from Specific Countries`,
							Type:       element.TypeMultiselect,
							SortOrder:  51,
							Visible:    element.VisibleYes,
							Scopes:     scope.PermWebsite,
							CanBeEmpty: true,
							// SourceModel: Magento\Directory\Model\Config\Source\Country
						},

						element.Field{
							// Path: payment/banktransfer/instructions
							ID:        "instructions",
							Label:     `Instructions`,
							Type:      element.TypeTextarea,
							SortOrder: 62,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},

						element.Field{
							// Path: payment/banktransfer/min_order_total
							ID:        "min_order_total",
							Label:     `Minimum Order Total`,
							Type:      element.TypeText,
							SortOrder: 98,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},

						element.Field{
							// Path: payment/banktransfer/max_order_total
							ID:        "max_order_total",
							Label:     `Maximum Order Total`,
							Type:      element.TypeText,
							SortOrder: 99,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},

						element.Field{
							// Path: payment/banktransfer/sort_order
							ID:        "sort_order",
							Label:     `Sort Order`,
							Type:      element.TypeText,
							SortOrder: 100,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},
					),
				},

				element.Group{
					ID:        "cashondelivery",
					Label:     `Cash On Delivery Payment`,
					SortOrder: 30,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: payment/cashondelivery/active
							ID:        "active",
							Label:     `Enabled`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: payment/cashondelivery/title
							ID:        "title",
							Label:     `Title`,
							Type:      element.TypeText,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `Cash On Delivery`,
						},

						element.Field{
							// Path: payment/cashondelivery/order_status
							ID:        "order_status",
							Label:     `New Order Status`,
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `pending`,
							// SourceModel: Magento\Sales\Model\Config\Source\Order\Status\NewStatus
						},

						element.Field{
							// Path: payment/cashondelivery/allowspecific
							ID:        "allowspecific",
							Label:     `Payment from Applicable Countries`,
							Type:      element.TypeAllowspecific,
							SortOrder: 50,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							// SourceModel: Magento\Payment\Model\Config\Source\Allspecificcountries
						},

						element.Field{
							// Path: payment/cashondelivery/specificcountry
							ID:         "specificcountry",
							Label:      `Payment from Specific Countries`,
							Type:       element.TypeMultiselect,
							SortOrder:  51,
							Visible:    element.VisibleYes,
							Scopes:     scope.PermWebsite,
							CanBeEmpty: true,
							// SourceModel: Magento\Directory\Model\Config\Source\Country
						},

						element.Field{
							// Path: payment/cashondelivery/instructions
							ID:        "instructions",
							Label:     `Instructions`,
							Type:      element.TypeTextarea,
							SortOrder: 62,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},

						element.Field{
							// Path: payment/cashondelivery/min_order_total
							ID:        "min_order_total",
							Label:     `Minimum Order Total`,
							Type:      element.TypeText,
							SortOrder: 98,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},

						element.Field{
							// Path: payment/cashondelivery/max_order_total
							ID:        "max_order_total",
							Label:     `Maximum Order Total`,
							Type:      element.TypeText,
							SortOrder: 99,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},

						element.Field{
							// Path: payment/cashondelivery/sort_order
							ID:        "sort_order",
							Label:     `Sort Order`,
							Type:      element.TypeText,
							SortOrder: 100,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},
					),
				},

				element.Group{
					ID:        "free",
					Label:     `Zero Subtotal Checkout`,
					SortOrder: 30,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: payment/free/active
							ID:        "active",
							Label:     `Enabled`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: payment/free/order_status
							ID:        "order_status",
							Label:     `New Order Status`,
							Type:      element.TypeSelect,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							// SourceModel: Magento\Sales\Model\Config\Source\Order\Status\Newprocessing
						},

						element.Field{
							// Path: payment/free/payment_action
							ID:        "payment_action",
							Label:     `Automatically Invoice All Items`,
							Type:      element.TypeSelect,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							// SourceModel: Magento\Payment\Model\Source\Invoice
						},

						element.Field{
							// Path: payment/free/sort_order
							ID:        "sort_order",
							Label:     `Sort Order`,
							Type:      element.TypeText,
							SortOrder: 100,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},

						element.Field{
							// Path: payment/free/title
							ID:        "title",
							Label:     `Title`,
							Type:      element.TypeText,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},

						element.Field{
							// Path: payment/free/allowspecific
							ID:        "allowspecific",
							Label:     `Payment from Applicable Countries`,
							Type:      element.TypeAllowspecific,
							SortOrder: 50,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							// SourceModel: Magento\Payment\Model\Config\Source\Allspecificcountries
						},

						element.Field{
							// Path: payment/free/specificcountry
							ID:         "specificcountry",
							Label:      `Payment from Specific Countries`,
							Type:       element.TypeMultiselect,
							SortOrder:  51,
							Visible:    element.VisibleYes,
							Scopes:     scope.PermWebsite,
							CanBeEmpty: true,
							// SourceModel: Magento\Directory\Model\Config\Source\Country
						},

						element.Field{
							// Path: payment/free/cfgmodel
							ID:      "model",
							Type:    element.Type,
							Visible: element.VisibleYes,
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		element.Section{
			ID: "payment",
			Groups: element.MakeGroups(
				element.Group{
					ID: "checkmo",
					Fields: element.MakeFields(
						element.Field{
							// Path: payment/checkmo/group
							ID:      `group`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `offline`,
						},
					),
				},

				element.Group{
					ID: "purchaseorder",
					Fields: element.MakeFields(
						element.Field{
							// Path: payment/purchaseorder/group
							ID:      `group`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `offline`,
						},
					),
				},

				element.Group{
					ID: "banktransfer",
					Fields: element.MakeFields(
						element.Field{
							// Path: payment/banktransfer/cfgmodel
							ID:      `model`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Magento\OfflinePayments\Model\Banktransfer`,
						},

						element.Field{
							// Path: payment/banktransfer/group
							ID:      `group`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `offline`,
						},
					),
				},

				element.Group{
					ID: "cashondelivery",
					Fields: element.MakeFields(
						element.Field{
							// Path: payment/cashondelivery/cfgmodel
							ID:      `model`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Magento\OfflinePayments\Model\Cashondelivery`,
						},

						element.Field{
							// Path: payment/cashondelivery/group
							ID:      `group`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `offline`,
						},
					),
				},

				element.Group{
					ID: "free",
					Fields: element.MakeFields(
						element.Field{
							// Path: payment/free/group
							ID:      `group`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `offline`,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
