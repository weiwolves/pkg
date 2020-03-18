// +build ignore

package ups

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
			ID: "carriers",
			Groups: element.MakeGroups(
				element.Group{
					ID:        "ups",
					Label:     `UPS`,
					SortOrder: 100,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: carriers/ups/access_license_number
							ID:        "access_license_number",
							Label:     `Access License Number`,
							Type:      element.TypeObscure,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   nil,
							// BackendModel: Magento\Config\Model\Config\Backend\Encrypted @todo Magento\Config\Model\Config\Backend\Encrypted
						},

						element.Field{
							// Path: carriers/ups/active
							ID:        "active",
							Label:     `Enabled for Checkout`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: carriers/ups/active_rma
							ID:        "active_rma",
							Label:     `Enabled for RMA`,
							Type:      element.TypeSelect,
							SortOrder: 15,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: carriers/ups/allowed_methods
							ID:         "allowed_methods",
							Label:      `Allowed Methods`,
							Type:       element.TypeMultiselect,
							SortOrder:  170,
							Visible:    element.VisibleYes,
							Scopes:     scope.PermWebsite,
							CanBeEmpty: true,
							Default:    `1DM,1DML,1DA,1DAL,1DAPI,1DP,1DPL,2DM,2DML,2DA,2DAL,3DS,GND,GNDCOM,GNDRES,STD,XPR,WXS,XPRL,XDM,XDML,XPD,01,02,03,07,08,11,12,14,54,59,65`,
							// SourceModel: Magento\Ups\Model\Config\Source\Method
						},

						element.Field{
							// Path: carriers/ups/shipment_requesttype
							ID:        "shipment_requesttype",
							Label:     `Packages Request Type`,
							Type:      element.TypeSelect,
							SortOrder: 47,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   false,
							// SourceModel: Magento\Shipping\Model\Config\Source\Online\Requesttype
						},

						element.Field{
							// Path: carriers/ups/container
							ID:        "container",
							Label:     `Container`,
							Type:      element.TypeSelect,
							SortOrder: 50,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `CP`,
							// SourceModel: Magento\Ups\Model\Config\Source\Container
						},

						element.Field{
							// Path: carriers/ups/free_shipping_enable
							ID:        "free_shipping_enable",
							Label:     `Free Shipping Amount Threshold`,
							Type:      element.TypeSelect,
							SortOrder: 210,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
						},

						element.Field{
							// Path: carriers/ups/free_shipping_subtotal
							ID:        "free_shipping_subtotal",
							Label:     `Free Shipping Amount Threshold`,
							Type:      element.TypeText,
							SortOrder: 220,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},

						element.Field{
							// Path: carriers/ups/dest_type
							ID:        "dest_type",
							Label:     `Destination Type`,
							Type:      element.TypeSelect,
							SortOrder: 60,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `RES`,
							// SourceModel: Magento\Ups\Model\Config\Source\DestType
						},

						element.Field{
							// Path: carriers/ups/free_method
							ID:        "free_method",
							Label:     `Free Method`,
							Type:      element.TypeSelect,
							SortOrder: 200,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `GND`,
							// SourceModel: Magento\Ups\Model\Config\Source\Freemethod
						},

						element.Field{
							// Path: carriers/ups/gateway_url
							ID:        "gateway_url",
							Label:     `Gateway URL`,
							Type:      element.TypeText,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `http://www.ups.com/using/services/rave/qcostcgi.cgi`,
						},

						element.Field{
							// Path: carriers/ups/gateway_xml_url
							ID:        "gateway_xml_url",
							Label:     `Gateway XML URL`,
							Type:      element.TypeText,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `https://onlinetools.ups.com/ups.app/xml/Rate`,
						},

						element.Field{
							// Path: carriers/ups/handling_type
							ID:        "handling_type",
							Label:     `Calculate Handling Fee`,
							Type:      element.TypeSelect,
							SortOrder: 110,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `F`,
							// SourceModel: Magento\Shipping\Model\Source\HandlingType
						},

						element.Field{
							// Path: carriers/ups/handling_action
							ID:        "handling_action",
							Label:     `Handling Applied`,
							Type:      element.TypeSelect,
							SortOrder: 120,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `O`,
							// SourceModel: Magento\Shipping\Model\Source\HandlingAction
						},

						element.Field{
							// Path: carriers/ups/handling_fee
							ID:        "handling_fee",
							Label:     `Handling Fee`,
							Type:      element.TypeText,
							SortOrder: 130,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},

						element.Field{
							// Path: carriers/ups/max_package_weight
							ID:        "max_package_weight",
							Label:     `Maximum Package Weight (Please consult your shipping carrier for maximum supported shipping weight)`,
							Type:      element.TypeText,
							SortOrder: 80,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   150,
						},

						element.Field{
							// Path: carriers/ups/min_package_weight
							ID:        "min_package_weight",
							Label:     `Minimum Package Weight (Please consult your shipping carrier for minimum supported shipping weight)`,
							Type:      element.TypeText,
							SortOrder: 90,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   0.1,
						},

						element.Field{
							// Path: carriers/ups/origin_shipment
							ID:        "origin_shipment",
							Label:     `Origin of the Shipment`,
							Type:      element.TypeSelect,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `Shipments Originating in United States`,
							// SourceModel: Magento\Ups\Model\Config\Source\OriginShipment
						},

						element.Field{
							// Path: carriers/ups/password
							ID:        "password",
							Label:     `Password`,
							Type:      element.TypeObscure,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   nil,
							// BackendModel: Magento\Config\Model\Config\Backend\Encrypted @todo Magento\Config\Model\Config\Backend\Encrypted
						},

						element.Field{
							// Path: carriers/ups/pickup
							ID:        "pickup",
							Label:     `Pickup Method`,
							Type:      element.TypeSelect,
							SortOrder: 80,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `CC`,
							// SourceModel: Magento\Ups\Model\Config\Source\Pickup
						},

						element.Field{
							// Path: carriers/ups/sort_order
							ID:        "sort_order",
							Label:     `Sort Order`,
							Type:      element.TypeText,
							SortOrder: 1000,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},

						element.Field{
							// Path: carriers/ups/title
							ID:        "title",
							Label:     `Title`,
							Type:      element.TypeText,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `United Parcel Service`,
						},

						element.Field{
							// Path: carriers/ups/tracking_xml_url
							ID:        "tracking_xml_url",
							Label:     `Tracking XML URL`,
							Type:      element.TypeText,
							SortOrder: 60,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `https://www.ups.com/ups.app/xml/Track`,
						},

						element.Field{
							// Path: carriers/ups/type
							ID:        "type",
							Label:     `UPS Type`,
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `UPS`,
							// SourceModel: Magento\Ups\Model\Config\Source\Type
						},

						element.Field{
							// Path: carriers/ups/is_account_live
							ID:        "is_account_live",
							Label:     `Live Account`,
							Type:      element.TypeSelect,
							SortOrder: 25,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: carriers/ups/unit_of_measure
							ID:        "unit_of_measure",
							Label:     `Weight Unit`,
							Type:      element.TypeSelect,
							SortOrder: 60,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `LBS`,
							// SourceModel: Magento\Ups\Model\Config\Source\Unitofmeasure
						},

						element.Field{
							// Path: carriers/ups/username
							ID:        "username",
							Label:     `User ID`,
							Type:      element.TypeObscure,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   nil,
							// BackendModel: Magento\Config\Model\Config\Backend\Encrypted @todo Magento\Config\Model\Config\Backend\Encrypted
						},

						element.Field{
							// Path: carriers/ups/negotiated_active
							ID:        "negotiated_active",
							Label:     `Enable Negotiated Rates`,
							Type:      element.TypeSelect,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: carriers/ups/shipper_number
							ID:        "shipper_number",
							Label:     `Shipper Number`,
							Comment:   text.Long(`Required for negotiated rates; 6-character UPS`),
							Type:      element.TypeText,
							SortOrder: 50,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},

						element.Field{
							// Path: carriers/ups/sallowspecific
							ID:        "sallowspecific",
							Label:     `Ship to Applicable Countries`,
							Type:      element.TypeSelect,
							SortOrder: 900,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   false,
							// SourceModel: Magento\Shipping\Model\Config\Source\Allspecificcountries
						},

						element.Field{
							// Path: carriers/ups/specificcountry
							ID:         "specificcountry",
							Label:      `Ship to Specific Countries`,
							Type:       element.TypeMultiselect,
							SortOrder:  910,
							Visible:    element.VisibleYes,
							Scopes:     scope.PermWebsite,
							CanBeEmpty: true,
							// SourceModel: Magento\Directory\Model\Config\Source\Country
						},

						element.Field{
							// Path: carriers/ups/showmethod
							ID:        "showmethod",
							Label:     `Show Method if Not Applicable`,
							Type:      element.TypeSelect,
							SortOrder: 920,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: carriers/ups/specificerrmsg
							ID:        "specificerrmsg",
							Label:     `Displayed Error Message`,
							Type:      element.TypeTextarea,
							SortOrder: 800,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `This shipping method is currently unavailable. If you would like to ship using this shipping method, please contact us.`,
						},

						element.Field{
							// Path: carriers/ups/mode_xml
							ID:        "mode_xml",
							Label:     `Mode`,
							Comment:   text.Long(`This enables or disables SSL verification of the Magento server by UPS.`),
							Type:      element.TypeSelect,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   true,
							// SourceModel: Magento\Shipping\Model\Config\Source\Online\Mode
						},

						element.Field{
							// Path: carriers/ups/debug
							ID:        "debug",
							Label:     `Debug`,
							Type:      element.TypeSelect,
							SortOrder: 920,
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
			ID: "carriers",
			Groups: element.MakeGroups(
				element.Group{
					ID: "ups",
					Fields: element.MakeFields(
						element.Field{
							// Path: carriers/ups/cutoff_cost
							ID:      `cutoff_cost`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
						},

						element.Field{
							// Path: carriers/ups/handling
							ID:      `handling`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},

						element.Field{
							// Path: carriers/ups/cfgmodel
							ID:      `model`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Magento\Ups\Model\Carrier`,
						},

						element.Field{
							// Path: carriers/ups/is_online
							ID:      `is_online`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
