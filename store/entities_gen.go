// Code generated by codegen. DO NOT EDIT.
// Generated by sql/dmlgen. DO NOT EDIT.
package store

import (
	"fmt"
	"io"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/storage/null"
)

// Store represents a single row for DB table store. Auto generated.
//easyjson:json
type Store struct {
	StoreID      uint32        `max_len:"5"`   // store_id smallint(5) unsigned NOT NULL PRI  auto_increment "Store ID"
	Code         null.String   `max_len:"64"`  // code varchar(64) NULL UNI DEFAULT 'NULL'  "Code"
	WebsiteID    uint32        `max_len:"5"`   // website_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'  "Website ID"
	GroupID      uint32        `max_len:"5"`   // group_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'  "Group ID"
	Name         string        `max_len:"255"` // name varchar(255) NOT NULL    "Store Name"
	SortOrder    uint32        `max_len:"5"`   // sort_order smallint(5) unsigned NOT NULL  DEFAULT '0'  "Store Sort Order"
	IsActive     bool          `max_len:"5"`   // is_active smallint(5) unsigned NOT NULL MUL DEFAULT '0'  "Store Activity"
	StoreGroup   *StoreGroup   // 1:1 store.group_id => store_group.group_id
	StoreWebsite *StoreWebsite // 1:1 store.website_id => store_website.website_id
}

// Empty empties all the fields of the current object. Also known as Reset.
func (e *Store) Empty() *Store { *e = Store{}; return e }

// Copy copies the struct and returns a new pointer
func (e *Store) Copy() *Store {
	e2 := new(Store)
	*e2 = *e // for now a shallow copy
	return e2
}

// WriteTo implements io.WriterTo and writes the field names and their values to
// w. This is especially useful for debugging or or generating a hash of the
// struct.
func (e *Store) WriteTo(w io.Writer) (n int64, err error) {
	// for now this printing is good enough. If you need better swap out with your code.
	n2, err := fmt.Fprint(w,
		"store_id:", e.StoreID, "\n",
		"code:", e.Code, "\n",
		"website_id:", e.WebsiteID, "\n",
		"group_id:", e.GroupID, "\n",
		"name:", e.Name, "\n",
		"sort_order:", e.SortOrder, "\n",
		"is_active:", e.IsActive, "\n",
	)
	return int64(n2), err
}

// StoreCollection represents a collection type for DB table store
// Not thread safe. Auto generated.
//easyjson:json
type StoreCollection struct {
	Data []*Store `json:"data,omitempty"`
}

// NewStoreCollection  creates a new initialized collection. Auto generated.
func NewStoreCollection() *StoreCollection {
	return &StoreCollection{
		Data: make([]*Store, 0, 5),
	}
}

// StoreIDs returns a slice with the data or appends it to a slice.
// Auto generated.
func (cc *StoreCollection) StoreIDs(ret ...uint32) []uint32 {
	if ret == nil {
		ret = make([]uint32, 0, len(cc.Data))
	}
	for _, e := range cc.Data {
		ret = append(ret, e.StoreID)
	}
	return ret
}

// Codes returns a slice with the data or appends it to a slice.
// Auto generated.
func (cc *StoreCollection) Codes(ret ...null.String) []null.String {
	if ret == nil {
		ret = make([]null.String, 0, len(cc.Data))
	}
	for _, e := range cc.Data {
		ret = append(ret, e.Code)
	}
	return ret
}

// WriteTo implements io.WriterTo and writes the field names and their values to
// w. This is especially useful for debugging or or generating a hash of the
// struct.
func (cc *StoreCollection) WriteTo(w io.Writer) (n int64, err error) {
	for i, d := range cc.Data {
		n2, err := d.WriteTo(w)
		if err != nil {
			return 0, errors.Wrapf(err, "[store] WriteTo failed at index %d", i)
		}
		n += n2
	}
	return n, nil
}

// Filter filters the current slice by predicate f without memory allocation.
// Auto generated via dmlgen.
func (cc *StoreCollection) Filter(f func(*Store) bool) *StoreCollection {
	b, i := cc.Data[:0], 0
	for _, e := range cc.Data {
		if f(e) {
			b = append(b, e)
			cc.Data[i] = nil // this avoids the memory leak
		}
		i++
	}
	cc.Data = b
	return cc
}

// Each will run function f on all items in []* Store . Auto generated via
// dmlgen.
func (cc *StoreCollection) Each(f func(*Store)) *StoreCollection {
	for i := range cc.Data {
		f(cc.Data[i])
	}
	return cc
}

// Cut will remove items i through j-1. Auto generated via dmlgen.
func (cc *StoreCollection) Cut(i, j int) *StoreCollection {
	z := cc.Data // copy slice header
	copy(z[i:], z[j:])
	for k, n := len(z)-j+i, len(z); k < n; k++ {
		z[k] = nil // this avoids the memory leak
	}
	z = z[:len(z)-j+i]
	cc.Data = z
	return cc
}

// Swap will satisfy the sort.Interface. Auto generated via dmlgen.
func (cc *StoreCollection) Swap(i, j int) { cc.Data[i], cc.Data[j] = cc.Data[j], cc.Data[i] }

// Len will satisfy the sort.Interface. Auto generated via dmlgen.
func (cc *StoreCollection) Len() int { return len(cc.Data) }

// Delete will remove an item from the slice. Auto generated via dmlgen.
func (cc *StoreCollection) Delete(i int) *StoreCollection {
	z := cc.Data // copy the slice header
	end := len(z) - 1
	cc.Swap(i, end)
	copy(z[i:], z[i+1:])
	z[end] = nil // this should avoid the memory leak
	z = z[:end]
	cc.Data = z
	return cc
}

// Insert will place a new item at position i. Auto generated via dmlgen.
func (cc *StoreCollection) Insert(n *Store, i int) *StoreCollection {
	z := cc.Data // copy the slice header
	z = append(z, &Store{})
	copy(z[i+1:], z[i:])
	z[i] = n
	cc.Data = z
	return cc
}

// Append will add a new item at the end of * StoreCollection . Auto generated
// via dmlgen.
func (cc *StoreCollection) Append(n ...*Store) *StoreCollection {
	cc.Data = append(cc.Data, n...)
	return cc
}

// StoreGroup represents a single row for DB table store_group. Auto generated.
//easyjson:json
type StoreGroup struct {
	GroupID        uint32        `max_len:"5"`   // group_id smallint(5) unsigned NOT NULL PRI  auto_increment "Group ID"
	WebsiteID      uint32        `max_len:"5"`   // website_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'  "Website ID"
	Name           string        `max_len:"255"` // name varchar(255) NOT NULL    "Store Group Name"
	RootCategoryID uint32        `max_len:"10"`  // root_category_id int(10) unsigned NOT NULL  DEFAULT '0'  "Root Category ID"
	DefaultStoreID uint32        `max_len:"5"`   // default_store_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'  "Default Store ID"
	Code           null.String   `max_len:"64"`  // code varchar(64) NULL UNI DEFAULT 'NULL'  "Store group unique code"
	StoreWebsite   *StoreWebsite // 1:1 store_group.website_id => store_website.website_id
}

// Empty empties all the fields of the current object. Also known as Reset.
func (e *StoreGroup) Empty() *StoreGroup { *e = StoreGroup{}; return e }

// Copy copies the struct and returns a new pointer
func (e *StoreGroup) Copy() *StoreGroup {
	e2 := new(StoreGroup)
	*e2 = *e // for now a shallow copy
	return e2
}

// WriteTo implements io.WriterTo and writes the field names and their values to
// w. This is especially useful for debugging or or generating a hash of the
// struct.
func (e *StoreGroup) WriteTo(w io.Writer) (n int64, err error) {
	// for now this printing is good enough. If you need better swap out with your code.
	n2, err := fmt.Fprint(w,
		"group_id:", e.GroupID, "\n",
		"website_id:", e.WebsiteID, "\n",
		"name:", e.Name, "\n",
		"root_category_id:", e.RootCategoryID, "\n",
		"default_store_id:", e.DefaultStoreID, "\n",
		"code:", e.Code, "\n",
	)
	return int64(n2), err
}

// StoreGroupCollection represents a collection type for DB table store_group
// Not thread safe. Auto generated.
//easyjson:json
type StoreGroupCollection struct {
	Data []*StoreGroup `json:"data,omitempty"`
}

// NewStoreGroupCollection  creates a new initialized collection. Auto generated.
func NewStoreGroupCollection() *StoreGroupCollection {
	return &StoreGroupCollection{
		Data: make([]*StoreGroup, 0, 5),
	}
}

// GroupIDs returns a slice with the data or appends it to a slice.
// Auto generated.
func (cc *StoreGroupCollection) GroupIDs(ret ...uint32) []uint32 {
	if ret == nil {
		ret = make([]uint32, 0, len(cc.Data))
	}
	for _, e := range cc.Data {
		ret = append(ret, e.GroupID)
	}
	return ret
}

// Codes returns a slice with the data or appends it to a slice.
// Auto generated.
func (cc *StoreGroupCollection) Codes(ret ...null.String) []null.String {
	if ret == nil {
		ret = make([]null.String, 0, len(cc.Data))
	}
	for _, e := range cc.Data {
		ret = append(ret, e.Code)
	}
	return ret
}

// WriteTo implements io.WriterTo and writes the field names and their values to
// w. This is especially useful for debugging or or generating a hash of the
// struct.
func (cc *StoreGroupCollection) WriteTo(w io.Writer) (n int64, err error) {
	for i, d := range cc.Data {
		n2, err := d.WriteTo(w)
		if err != nil {
			return 0, errors.Wrapf(err, "[store] WriteTo failed at index %d", i)
		}
		n += n2
	}
	return n, nil
}

// Filter filters the current slice by predicate f without memory allocation.
// Auto generated via dmlgen.
func (cc *StoreGroupCollection) Filter(f func(*StoreGroup) bool) *StoreGroupCollection {
	b, i := cc.Data[:0], 0
	for _, e := range cc.Data {
		if f(e) {
			b = append(b, e)
			cc.Data[i] = nil // this avoids the memory leak
		}
		i++
	}
	cc.Data = b
	return cc
}

// Each will run function f on all items in []* StoreGroup . Auto generated via
// dmlgen.
func (cc *StoreGroupCollection) Each(f func(*StoreGroup)) *StoreGroupCollection {
	for i := range cc.Data {
		f(cc.Data[i])
	}
	return cc
}

// Cut will remove items i through j-1. Auto generated via dmlgen.
func (cc *StoreGroupCollection) Cut(i, j int) *StoreGroupCollection {
	z := cc.Data // copy slice header
	copy(z[i:], z[j:])
	for k, n := len(z)-j+i, len(z); k < n; k++ {
		z[k] = nil // this avoids the memory leak
	}
	z = z[:len(z)-j+i]
	cc.Data = z
	return cc
}

// Swap will satisfy the sort.Interface. Auto generated via dmlgen.
func (cc *StoreGroupCollection) Swap(i, j int) { cc.Data[i], cc.Data[j] = cc.Data[j], cc.Data[i] }

// Len will satisfy the sort.Interface. Auto generated via dmlgen.
func (cc *StoreGroupCollection) Len() int { return len(cc.Data) }

// Delete will remove an item from the slice. Auto generated via dmlgen.
func (cc *StoreGroupCollection) Delete(i int) *StoreGroupCollection {
	z := cc.Data // copy the slice header
	end := len(z) - 1
	cc.Swap(i, end)
	copy(z[i:], z[i+1:])
	z[end] = nil // this should avoid the memory leak
	z = z[:end]
	cc.Data = z
	return cc
}

// Insert will place a new item at position i. Auto generated via dmlgen.
func (cc *StoreGroupCollection) Insert(n *StoreGroup, i int) *StoreGroupCollection {
	z := cc.Data // copy the slice header
	z = append(z, &StoreGroup{})
	copy(z[i+1:], z[i:])
	z[i] = n
	cc.Data = z
	return cc
}

// Append will add a new item at the end of * StoreGroupCollection . Auto
// generated via dmlgen.
func (cc *StoreGroupCollection) Append(n ...*StoreGroup) *StoreGroupCollection {
	cc.Data = append(cc.Data, n...)
	return cc
}

// StoreWebsite represents a single row for DB table store_website. Auto
// generated.
//easyjson:json
type StoreWebsite struct {
	WebsiteID      uint32                `max_len:"5"`   // website_id smallint(5) unsigned NOT NULL PRI  auto_increment "Website ID"
	Code           null.String           `max_len:"64"`  // code varchar(64) NULL UNI DEFAULT 'NULL'  "Code"
	Name           null.String           `max_len:"128"` // name varchar(128) NULL  DEFAULT 'NULL'  "Website Name"
	SortOrder      uint32                `max_len:"5"`   // sort_order smallint(5) unsigned NOT NULL MUL DEFAULT '0'  "Sort Order"
	DefaultGroupID uint32                `max_len:"5"`   // default_group_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'  "Default Group ID"
	IsDefault      bool                  `max_len:"5"`   // is_default smallint(5) unsigned NOT NULL  DEFAULT '0'  "Defines Is Website Default"
	StoreGroup     *StoreGroupCollection // Reversed 1:M store_website.website_id => store_group.website_id
	Store          *StoreCollection      // Reversed 1:M store_website.website_id => store.website_id
}

// Empty empties all the fields of the current object. Also known as Reset.
func (e *StoreWebsite) Empty() *StoreWebsite { *e = StoreWebsite{}; return e }

// Copy copies the struct and returns a new pointer
func (e *StoreWebsite) Copy() *StoreWebsite {
	e2 := new(StoreWebsite)
	*e2 = *e // for now a shallow copy
	return e2
}

// WriteTo implements io.WriterTo and writes the field names and their values to
// w. This is especially useful for debugging or or generating a hash of the
// struct.
func (e *StoreWebsite) WriteTo(w io.Writer) (n int64, err error) {
	// for now this printing is good enough. If you need better swap out with your code.
	n2, err := fmt.Fprint(w,
		"website_id:", e.WebsiteID, "\n",
		"code:", e.Code, "\n",
		"name:", e.Name, "\n",
		"sort_order:", e.SortOrder, "\n",
		"default_group_id:", e.DefaultGroupID, "\n",
		"is_default:", e.IsDefault, "\n",
	)
	return int64(n2), err
}

// StoreWebsiteCollection represents a collection type for DB table store_website
// Not thread safe. Auto generated.
//easyjson:json
type StoreWebsiteCollection struct {
	Data []*StoreWebsite `json:"data,omitempty"`
}

// NewStoreWebsiteCollection  creates a new initialized collection. Auto
// generated.
func NewStoreWebsiteCollection() *StoreWebsiteCollection {
	return &StoreWebsiteCollection{
		Data: make([]*StoreWebsite, 0, 5),
	}
}

// WebsiteIDs returns a slice with the data or appends it to a slice.
// Auto generated.
func (cc *StoreWebsiteCollection) WebsiteIDs(ret ...uint32) []uint32 {
	if ret == nil {
		ret = make([]uint32, 0, len(cc.Data))
	}
	for _, e := range cc.Data {
		ret = append(ret, e.WebsiteID)
	}
	return ret
}

// Codes returns a slice with the data or appends it to a slice.
// Auto generated.
func (cc *StoreWebsiteCollection) Codes(ret ...null.String) []null.String {
	if ret == nil {
		ret = make([]null.String, 0, len(cc.Data))
	}
	for _, e := range cc.Data {
		ret = append(ret, e.Code)
	}
	return ret
}

// WriteTo implements io.WriterTo and writes the field names and their values to
// w. This is especially useful for debugging or or generating a hash of the
// struct.
func (cc *StoreWebsiteCollection) WriteTo(w io.Writer) (n int64, err error) {
	for i, d := range cc.Data {
		n2, err := d.WriteTo(w)
		if err != nil {
			return 0, errors.Wrapf(err, "[store] WriteTo failed at index %d", i)
		}
		n += n2
	}
	return n, nil
}

// Filter filters the current slice by predicate f without memory allocation.
// Auto generated via dmlgen.
func (cc *StoreWebsiteCollection) Filter(f func(*StoreWebsite) bool) *StoreWebsiteCollection {
	b, i := cc.Data[:0], 0
	for _, e := range cc.Data {
		if f(e) {
			b = append(b, e)
			cc.Data[i] = nil // this avoids the memory leak
		}
		i++
	}
	cc.Data = b
	return cc
}

// Each will run function f on all items in []* StoreWebsite . Auto generated via
// dmlgen.
func (cc *StoreWebsiteCollection) Each(f func(*StoreWebsite)) *StoreWebsiteCollection {
	for i := range cc.Data {
		f(cc.Data[i])
	}
	return cc
}

// Cut will remove items i through j-1. Auto generated via dmlgen.
func (cc *StoreWebsiteCollection) Cut(i, j int) *StoreWebsiteCollection {
	z := cc.Data // copy slice header
	copy(z[i:], z[j:])
	for k, n := len(z)-j+i, len(z); k < n; k++ {
		z[k] = nil // this avoids the memory leak
	}
	z = z[:len(z)-j+i]
	cc.Data = z
	return cc
}

// Swap will satisfy the sort.Interface. Auto generated via dmlgen.
func (cc *StoreWebsiteCollection) Swap(i, j int) { cc.Data[i], cc.Data[j] = cc.Data[j], cc.Data[i] }

// Len will satisfy the sort.Interface. Auto generated via dmlgen.
func (cc *StoreWebsiteCollection) Len() int { return len(cc.Data) }

// Delete will remove an item from the slice. Auto generated via dmlgen.
func (cc *StoreWebsiteCollection) Delete(i int) *StoreWebsiteCollection {
	z := cc.Data // copy the slice header
	end := len(z) - 1
	cc.Swap(i, end)
	copy(z[i:], z[i+1:])
	z[end] = nil // this should avoid the memory leak
	z = z[:end]
	cc.Data = z
	return cc
}

// Insert will place a new item at position i. Auto generated via dmlgen.
func (cc *StoreWebsiteCollection) Insert(n *StoreWebsite, i int) *StoreWebsiteCollection {
	z := cc.Data // copy the slice header
	z = append(z, &StoreWebsite{})
	copy(z[i+1:], z[i:])
	z[i] = n
	cc.Data = z
	return cc
}

// Append will add a new item at the end of * StoreWebsiteCollection . Auto
// generated via dmlgen.
func (cc *StoreWebsiteCollection) Append(n ...*StoreWebsite) *StoreWebsiteCollection {
	cc.Data = append(cc.Data, n...)
	return cc
}