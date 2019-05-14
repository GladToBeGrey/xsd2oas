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

// writeExample
// Take the populated data structures and output a sample message

package main

import (
	"fmt"
	"io"

	"github.com/lucasjones/reggen"
)

const tab = "  "

// entry point for writing
func writeExample(f io.Writer, ctxt *context) {

	indent := ""
	path := ""
	doc := ctxt.complexTypes["Document"]
	// fmt.Printf("Got Document%v\n", doc)
	fmt.Fprintf(f, "%v{\n", indent)
	writeOne(f, ctxt, doc, path, indent)
	fmt.Fprintf(f, "%v}\n", indent)
}

func writeOne(f io.Writer, ctxt *context, cplx complexType, path string, indent string) {
	for idx, el := range cplx.elems {
		arOpen, arClose := "", ""
		if idx > 0 {
			fmt.Fprintf(f, ",\n")
		}
		if el.maxOccurs > 1 {
			arOpen, arClose = "[", "]"
		}
		if t, ok := ctxt.complexTypes[el.etype]; ok {
			//process complex type
			// fmt.Printf("Path:%v(%v)\n", path+"/"+el.name, el.etype)
			// fmt.Printf("Path:%v\n", path+"/"+el.name)
			fmt.Fprintf(f, "%v\"%v\": %v{\n", indent+tab, el.name, arOpen)
			writeOne(f, ctxt, t, path+"/"+el.name, indent+tab)
			fmt.Fprintf(f, "%v", arClose)
		} else {
			//process simple type
			s := ctxt.simpleTypes[el.etype]
			if len(s.attrs) == 0 {
				// fmt.Printf("Path:%v(%v)\n", path+"/"+el.name, s.base)
				fmt.Fprintf(f, "%v\"%v\": %v %v %v", indent+tab, el.name, arOpen, sampleData(s), arClose)
			} else {
				// fmt.Printf("Path:%v\n", path+"/"+el.name)
				fmt.Fprintf(f, "%v\"%v\": %v{\n", indent+tab, el.name, arOpen)
				// fmt.Printf("Path:%v(%v)\n", path+"/"+el.name+"/value", s.base)
				fmt.Fprintf(f, "%v\"%v\": %v,\n", indent+tab+tab, "value", sampleData(s))
				for idx, attr := range s.attrs {
					if idx > 0 {
						fmt.Fprintf(f, ",\n")
					}
					// fmt.Printf("Path:%v(%v)\n", path+"/"+el.name+"/@"+attr.name, "string")
					atype := ctxt.simpleTypes[attr.atype]
					fmt.Fprintf(f, "%v\"%v\": %v", indent+tab+tab, "@"+attr.name, sampleData(atype))
				}
				fmt.Fprintf(f, "\n%v}%v", indent+tab, arClose)
			}
		}
	}
	fmt.Fprintf(f, "\n%v}", indent)
}

func sampleData(s simpleType) string {
	jname, _ := mapTypename(s.base)
	switch jname {
	case "boolean":
		return "true"
	case "number":
		return "123456"
	case "string":
		switch {
		case s.pattern != "":
			str, err := reggen.Generate(s.pattern, 10)
			if err != nil {
				panic(err)
			}
			return fmt.Sprintf("\"%v\"", str)
		case len(s.enum) > 0:
			return fmt.Sprintf("\"%v\"", s.enum[0])
		default:
			min := s.minLength
			max := s.maxLength
			if min < 1 {
				min = 1
			}
			if max < 1 {
				max = 10
			}
			if max > 1000 {
				max = 1000
			} // undocumented golang regex limit!
			patt := fmt.Sprintf("[0-9A-Fa-f]{%v,%v}", min, max)
			str, err := reggen.Generate(patt, 10)
			if err != nil {
				panic(err)
			}
			return fmt.Sprintf("\"%v\"", str)
		}
	}
	return s.base
}
