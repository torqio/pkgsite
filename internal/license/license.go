// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package license

import "path"

// Metadata holds information extracted from a license file: its license Type
// and FilePath relative to the contents directory.
type Metadata struct {
	// Types is the set of license types, as determined by the licensecheck package.
	Types []string
	// FilePath is the '/'-separated path to the license file in the module zip,
	// relative to the contents directory.
	FilePath string
}

// A License is a classified license file path and its contents.
type License struct {
	*Metadata
	Contents []byte
}

// redistributableLicenses defines the set of (detectable) licenses that are
// considered permissive of redistribution.
var redistributableLicenses = map[string]bool{
	// Licenses acceptable by OSI
	"AGPL-3.0":             true,
	"Apache-2.0":           true,
	"Artistic-2.0":         true,
	"BSD-2-Clause":         true,
	"BSD-2-Clause-FreeBSD": true,
	"BSD-3-Clause":         true,
	"BSL-1.0":              true,
	"GPL2":                 true,
	"GPL3":                 true,
	"ISC":                  true,
	"LGPL-2.1":             true,
	"LGPL-3.0":             true,
	"MIT":                  true,
	"MPL-2.0":              true,
	"Zlib":                 true,
}

// AreRedistributable reports whether content subject to the given licenses
// should be considered redistributable. The current algorithm for this is to
// ensure that:
//   1. Every applicable license is redistributable, and
//   2. There is at least one license in the root directory.
func AreRedistributable(licenses []*Metadata) bool {
	haveRootLicense := false
	for _, lm := range licenses {
		if path.Dir(lm.FilePath) == "." {
			haveRootLicense = true
		}
		if !lm.IsRedistributable() {
			return false
		}
	}
	return haveRootLicense
}

// IsRedistributable reports whether content subject to the given
// license should be considered redistributable. The license
// must have at least one type, and all its types must be redistributable.
func (lm *Metadata) IsRedistributable() bool {
	if len(lm.Types) == 0 {
		return false
	}
	for _, typ := range lm.Types {
		if !redistributableLicenses[typ] {
			return false
		}
	}
	return true
}

// ToMetadatas converts a slice of Licenses to a slice of Metadatas.
func ToMetadatas(lics []*License) []*Metadata {
	var ms []*Metadata
	for _, l := range lics {
		ms = append(ms, l.Metadata)
	}
	return ms
}
