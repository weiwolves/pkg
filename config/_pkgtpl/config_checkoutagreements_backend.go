// +build ignore

package checkoutagreements

import (
	"github.com/weiwolves/pkg/config/cfgmodel"
	"github.com/weiwolves/pkg/config/element"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	cfgmodel.PkgBackend
	// CheckoutOptionsEnableAgreements => Enable Terms and Conditions.
	// Path: checkout/options/enable_agreements
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CheckoutOptionsEnableAgreements cfgmodel.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.Sections) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.Sections) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CheckoutOptionsEnableAgreements = cfgmodel.NewBool(`checkout/options/enable_agreements`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
