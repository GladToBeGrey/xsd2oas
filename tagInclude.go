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
	"io"
	"strings"
)

// entry point for tagging
func tagInclude(f io.Writer, ctxt *context) {

	path := ""
	doc := ctxt.complexTypes["Document"]
	// fmt.Printf("Got Document%v\n", doc)
	tagOne(ctxt, &doc, path, f)
}

func tagOne(ctxt *context, cplx *complexType, path string, f io.Writer) {
	// fmt.Printf("Tagging: %v\n", path)
	for idx, el := range cplx.elems {
		// fmt.Printf("Checking: %v (%v)\n", path+"/"+el.name, el.minOccurs)
		choice := cplx.etype == "choice"
		rqd := ctxt.all || // user specified all on command line
			!(choice || el.minOccurs == 0) || // XSD specifies mandatory: minOccurs -1 means unspecified, default 1
			isRequired(ctxt, path+"/"+el.name) // mask file requires inclusion
		if rqd {
			el.include = true
			cplx.elems[idx] = el
			if t, ok := ctxt.complexTypes[el.etype]; ok {
				//process complex type
				t.include = true
				tagOne(ctxt, &t, path+"/"+el.name, f)
				ctxt.complexTypes[el.etype] = t
				// fmt.Printf("Complex: %v\n", path+"/"+el.name)
			} else {
				if f != nil {
					fmt.Fprintf(f, "%v\n", path+"/"+el.name)
				}
				t := ctxt.simpleTypes[el.etype]
				t.include = true
				ctxt.simpleTypes[el.etype] = t
				// fmt.Printf("Simple: %v\n", path+"/"+el.name)
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
