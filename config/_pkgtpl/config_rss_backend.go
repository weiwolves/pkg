// +build ignore

package rss

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
	// RssConfigActive => Enable RSS.
	// Path: rss/config/active
	// BackendModel: Magento\Rss\Model\System\Config\Backend\Links
	// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
	RssConfigActive cfgmodel.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.Sections) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.Sections) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.RssConfigActive = cfgmodel.NewBool(`rss/config/active`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
