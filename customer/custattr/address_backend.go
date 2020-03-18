// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package custattr

import "github.com/weiwolves/pkg/eav"

// AddressDataPostcode post code data model @todo
// @see magento2/site/app/code/Magento/Customer/Model/Attribute/Data/Postcode.php
func AddressBackendRegion() *eav.AttributeBackend {
	return eav.NewAttributeBackend()
}

// AddressBackendStreet handles multiline street address @todo
// @see Mage_Customer_Model_Resource_Address_Attribute_Backend_Street
func AddressBackendStreet() *eav.AttributeBackend {
	return eav.NewAttributeBackend()
}
