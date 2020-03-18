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

package backendjwt

import (
	"github.com/weiwolves/pkg/config/cfgpath"
	"github.com/weiwolves/pkg/config/element"
	"github.com/weiwolves/pkg/net/jwt"
	"github.com/weiwolves/pkg/storage/text"
	"github.com/weiwolves/pkg/store/scope"
)

// NewConfigStructure global configuration structure for this package.
// Used in frontend (to display the user all the settings) and in
// backend (scope checks and default values). See the source code
// of this function for the overall available sections, groups and fields.
func NewConfigStructure() (element.Sections, error) {
	return element.MakeSectionsValidated(
		element.Section{
			ID: cfgpath.MakeRoute("net"),
			Groups: element.MakeGroups(
				element.Group{
					ID:        cfgpath.MakeRoute("jwt"),
					Label:     text.Chars(`JSON Web Token (JWT)`),
					SortOrder: 40,
					Scopes:    scope.PermWebsite,
					Fields: element.MakeFields(
						element.Field{
							// Path: net/jwt/disabled
							ID:        cfgpath.MakeRoute("disabled"),
							Label:     text.Chars(`JSON Webtoken is disabled`),
							Comment:   text.Chars(`Disables completely the JWT validation. Set to true/enable to activate the disabling.`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   false,
						},
						element.Field{
							// Path: net/jwt/expiration
							ID:        cfgpath.MakeRoute("expiration"),
							Label:     text.Chars(`Token Expiration`),
							Comment:   text.Chars(`Per second (s), minute (i), hour (h) or day (d)`),
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   jwt.DefaultExpire.String(),
						},
						element.Field{
							// Path: net/jwt/skew
							ID:        cfgpath.MakeRoute("skew"),
							Label:     text.Chars(`Max time skew`),
							Comment:   text.Chars(`How much time skew we allow between signer and verifier. Per second (s), minute (i), hour (h) or day (d). Must be a positive value.`),
							Type:      element.TypeText,
							SortOrder: 25,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   jwt.DefaultSkew.String(),
						},
						element.Field{
							// Path: net/jwt/single_usage
							ID:        cfgpath.MakeRoute("single_usage"),
							Label:     text.Chars(`Enable single token usage`),
							Comment:   text.Chars(`If set to true for each request a token can only be used once. The JTI (JSON Token Identifier) gets added to the blacklist until it expires.`),
							Type:      element.TypeSelect,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `false`,
						},
						element.Field{
							// Path: net/jwt/signing_method
							ID:        cfgpath.MakeRoute("signing_method"),
							Label:     text.Chars(`Token Signing Algorithm`),
							Type:      element.TypeSelect,
							SortOrder: 35,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   jwt.DefaultSigningMethod,
						},
						element.Field{
							// Path: net/jwt/hmac_password
							ID:        cfgpath.MakeRoute("hmac_password"),
							Label:     text.Chars(`Global HMAC Token Password. If empty, a random very long password will be generated.`),
							Type:      element.TypeObscure,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},
						element.Field{
							// Path: net/jwt/hmac_password_per_user
							ID:        cfgpath.MakeRoute("hmac_password_per_user"),
							Label:     text.Chars(`Enable per user HMAC token password`),
							Comment:   text.Chars(`A random HMAC password will be generated for each user who logs in`),
							Type:      element.TypeSelect,
							SortOrder: 45,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   0,
						},
						element.Field{
							// Path: net/jwt/rsa_key
							ID:        cfgpath.MakeRoute("rsa_key"),
							Label:     text.Chars(`Private RSA Key`),
							Type:      element.TypeObscure,
							SortOrder: 50,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},
						element.Field{
							// Path: net/jwt/rsa_key_password
							ID:        cfgpath.MakeRoute("rsa_key_password"),
							Label:     text.Chars(`Private RSA Key Password`),
							Comment:   text.Chars(`If the key has been secured via a password, provide it here.`),
							Type:      element.TypeObscure,
							SortOrder: 60,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},
						element.Field{
							// Path: net/jwt/ecdsa_key
							ID:        cfgpath.MakeRoute("ecdsa_key"),
							Label:     text.Chars(`Private ECDSA Key`),
							Comment:   text.Chars(`Elliptic Curve Digital Signature Algorithm, as defined in FIPS 186-3.`),
							Type:      element.TypeObscure,
							SortOrder: 70,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},
						element.Field{
							// Path: net/jwt/ecdsa_key_password
							ID:        cfgpath.MakeRoute("ecdsa_key_password"),
							Label:     text.Chars(`Private ECDSA Key Password`),
							Comment:   text.Chars(`If the key has been secured via a password, provide it here.`),
							Type:      element.TypeObscure,
							SortOrder: 80,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},
					),
				},
			),
		},
	)
}
