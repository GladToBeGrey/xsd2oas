// xsd2oas - convert XSD files to OpenAPI Specification
// Copyright (C) 2019  Tom Hay

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.package main

// tagInclude
// Tag the elements to include

package main

import (
	"fmt"
	"strings"
)

// entry point for tagging
func tagInclude(ctxt *context) {

	path := ""
	doc := ctxt.complexTypes["Document"]
	// fmt.Printf("Got Document%v\n", doc)
	tagOne(ctxt, &doc, path)
}

func tagOne(ctxt *context, cplx *complexType, path string) {
	for idx, el := range cplx.elems {
		rqd := el.minOccurs > 0 || isRequired(ctxt, path+"/"+el.name)
		if rqd {
			el.include = true
			cplx.elems[idx] = el
			if ctxt.verbose {
				fmt.Printf("%v\n", path+"/"+el.name)
			}
			if t, ok := ctxt.complexTypes[el.etype]; ok {
				//process complex type
				t.include = true
				tagOne(ctxt, &t, path+"/"+el.name)
				ctxt.complexTypes[el.etype] = t
				// fmt.Printf("Complex: %v\n", path+"/"+el.name)
			} else {
				t := ctxt.simpleTypes[el.etype]
				t.include = true
				ctxt.simpleTypes[el.etype] = t
				// fmt.Printf("Simple: %v\n", path+"/"+t.name)
				for _, attr := range t.attrs {
					t := ctxt.simpleTypes[attr.atype]
					t.include = true
					ctxt.simpleTypes[attr.atype] = t
				}
			}
		}
	}
}

func isRequired(ctxt *context, path string) bool {
	if !ctxt.mask {
		return true
	}
	for _, s := range ctxt.maskLines {
		if strings.HasPrefix(s, path) {
			return true
		}
	}
	return false
}
