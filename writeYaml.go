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

// writeYaml
// Take the populated data structures and output OAS in Yaml

package main

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
)

const tsz = 2 // tab size
var regupr *regexp.Regexp

func init() {
	regupr = regexp.MustCompile("[a-z]") // compiled regexp tp find uppercase
}

// entry point for writing
func writeYaml(f io.Writer, ctxt *context) {

	writeHdrs(f, ctxt, 0)
	writeComponents(f, ctxt, 0)
}

// write the name of an element or type
func writeName(n named, f io.Writer, ctxt *context, indent int) {
	inPrintf(f, indent, "%s:\n", n.getName())
}

// write an element
// if multiple occurrences are allowed, make it an array of items
// of the specified type
func writeElement(el *element, f io.Writer, ctxt *context, indent int) {
	name := el.getName()
	// if fixup flag set, convert param name that are all upppercase to camelcase
	if ctxt.fixUppercase {
		if regupr.FindStringIndex(name) == nil {
			name = name[0:1] + strings.ToLower(name[1:])
		}
	}
	if el.maxOccurs > 1 {
		inPrintf(f, indent, "%s:\n", name)
		inPrintf(f, indent+tsz, "type: array\n")
		inPrintf(f, indent+tsz, "items:\n")
		inPrintf(f, indent+tsz+tsz, "$ref: '#/components/schemas/%s'\n", el.etype)
	} else {
		inPrintf(f, indent, "%s:\n", name)
		inPrintf(f, indent+tsz, "$ref: '#/components/schemas/%s'\n", el.etype)
	}
}

// write the body of a simple type
// if it has attributes, turn it into an object
// the value element represents the base type
// each attribute forms a separate element named @Attributename
func writeSimpleBody(simple *simpleType, f io.Writer, ctxt *context, indent int) {
	if len(simple.attrs) > 0 {
		inPrintf(f, indent, "type: object\n")
		inPrintf(f, indent, "properties:\n")
		inPrintf(f, indent+tsz, "\"value\":\n")
		writeSimpleProperties(simple, f, ctxt, indent+tsz+tsz)
		required := writeAttrs(simple, f, ctxt, indent+tsz)
		inPrintf(f, indent, "required: %s\n", arrayString(required))
		inPrintf(f, indent, "additionalProperties: false\n")
	} else {
		writeSimpleProperties(simple, f, ctxt, indent)
	}
}

// write the properties of a simple type
func writeSimpleProperties(simple *simpleType, f io.Writer, ctxt *context, indent int) {
	jtype, mapped := mapTypename(simple.base)
	inPrintf(f, indent, "type: %s\n", jtype)
	if mapped {
		inPrintf(f, indent, "# XML datatype was %s\n", simple.base)
	}
	// string constraints
	if simple.minLength > -1 {
		inPrintf(f, indent, "minLength: %d\n", simple.minLength)
	}
	if simple.maxLength > -1 {
		inPrintf(f, indent, "maxLength: %d\n", simple.maxLength)
	}
	if simple.length > -1 {
		inPrintf(f, indent, "minLength: %d\n", simple.length)
		inPrintf(f, indent, "maxLength: %d\n", simple.length)
	}
	if len(simple.enum) > 0 {
		inPrintf(f, indent, "enum: %s\n", arrayString(simple.enum))
	}
	if simple.pattern != "" {
		// double all slashes to make valid JSON escapes
		escaped := strings.Replace(simple.pattern, "\\", "\\\\", -1)
		inPrintf(f, indent, "pattern: '%s'\n", escaped)
	}
	// number constraints
	if simple.minInclusive > -1 {
		inPrintf(f, indent, "minimum: %d\n", simple.minInclusive)
	}
	if simple.minExclusive > -1 {
		inPrintf(f, indent, "exclusiveMinimum: %d\n", simple.minExclusive)
	}
	if simple.maxInclusive > -1 {
		inPrintf(f, indent, "maximum: %d\n", simple.maxInclusive)
	}
	if simple.maxExclusive > -1 {
		inPrintf(f, indent, "exclusiveMaximum: %d\n", simple.maxExclusive)
	}
	// JSON schema can't handle these rules
	if simple.totalDigits > -1 {
		inPrintf(f, indent, "# XML specified totalDigits=%d\n", simple.totalDigits)
	}
	if simple.fractionDigits > -1 {
		inPrintf(f, indent, "# XML specified fractionDigits=%d\n", simple.fractionDigits)
	}
	if simple.whiteSpace != "" {
		inPrintf(f, indent, "# XML specified whiteSpace=%s\n", simple.whiteSpace)
	}
}

// write the file headers
func writeHdrs(f io.Writer, ctxt *context, indent int) {
	domain := "https://example.com"
	// when := time.Now().Format(time.RFC1123)
	if ctxt.domain != "" {
		domain = ctxt.domain
	}
	rootType := ctxt.complexTypes[ctxt.root.getName()]
	hdrs := [...]string{
		"openapi: 3.0.0\n",
		"info:\n",
		"\ttitle: '" + ctxt.outFileBase + "'\n",
		"\tversion: '0.1'\n",
		"\n",
		"servers:\n",
		"\t- url: '" + domain + "'\n",
		"\n",
		"paths:\n",
		"\t'/" + ctxt.outFileBase + "':\n",
		"\t\tput:\n",
		"\t\t\trequestBody:\n",
		"\t\t\t\tcontent:\n",
		"\t\t\t\t\tapplication/json:\n",
		"\t\t\t\t\t\tschema:\n",
		"\t\t\t\t\t\t\t$ref: '#/components/schemas/" + rootType.elems[0].etype + "'\n",
		"\t\t\tresponses:\n",
		"\t\t\t\t'200':\n",
		"\t\t\t\t\tdescription: Happy path\n",
		"\t\t\t\t'400':\n",
		"\t\t\t\t\tdescription: Bad request (body describes why)\n",
		"\t\t\t\t'410':\n",
		"\t\t\t\t\tdescription: Unauthorised\n",
		"\t\t\t\t'504':\n",
		"\t\t\t\t\tdescription: Gateway timeout (server did not respond)\n",
		"\t\t\t\t'5XX':\n",
		"\t\t\t\t\tdescription: Server Error\n",
	}
	for _, str := range hdrs {
		str = detab(str)
		inPrintf(f, indent, str)
	}
	inPrintf(f, 0, "\n")
}

// write all the component definitions
func writeComponents(f io.Writer, ctxt *context, indent int) {
	inPrintf(f, indent, "# ---Component definitions---\n")
	inPrintf(f, indent, "components:\n")
	writeSchemas(f, ctxt, indent+tsz)
}

// write all the schema definitions
func writeSchemas(f io.Writer, ctxt *context, indent int) {

	inPrintf(f, indent, "schemas:\n\n")

	// merge simple and complex together and sort them
	cmb := make([]string, 0)
	for _, simple := range ctxt.simpleTypes {
		if simple.include {
			cmb = append(cmb, simple.name)
		}
	}
	for _, cmplx := range ctxt.complexTypes {
		if cmplx.include {
			cmb = append(cmb, cmplx.name)
		}
	}
	sort.Strings(cmb)

	for _, nm := range cmb {
		if simple, ok := ctxt.simpleTypes[nm]; ok {
			writeSimple(simple, f, ctxt, indent+tsz)
		} else {
			writeComplex(ctxt.complexTypes[nm], f, ctxt, indent+tsz)
		}
	}
}

// write a simple type definition
func writeSimple(simple *simpleType, f io.Writer, ctxt *context, indent int) {
	writeName(simple, f, ctxt, indent)
	writeSimpleBody(simple, f, ctxt, indent+tsz)
}

// write a complex type definition
func writeComplex(cmplx *complexType, f io.Writer, ctxt *context, indent int) {
	writeName(cmplx, f, ctxt, indent)
	writeComplexBody(cmplx, f, ctxt, indent+tsz)
}

// write the body of a complex type
func writeComplexBody(cmplx *complexType, f io.Writer, ctxt *context, indent int) {
	// if it's based on simple, do simple body
	if cmplx.simpleBase != nil {
		fmt.Printf("Doing simple body for %s: %v\n", cmplx.name, *cmplx.simpleBase)
		writeSimpleBody(cmplx.simpleBase, f, ctxt, indent+tsz)
		return
	}

	switch cmplx.etype {
	case "choice":
		// XSD choice maps to YAML schema thus:
		//   "type": "object"
		//   "properties":
		//     "Pty":
		//       "$ref": "#/components/schemas/PartyTypeDef",
		//     "Agt":
		//       "$ref": "#/components/schemas/AgentTypeDef",
		//   "oneOf":
		//   - required: [Pty]
		//   - required: [Agt]
		inPrintf(f, indent, "type: object\n")
		inPrintf(f, indent, "properties:\n")
		for _, el := range cmplx.elems {
			if el.include {
				writeElement(el, f, ctxt, indent+tsz)
			}
		}
		inPrintf(f, indent, "oneOf:\n")
		for _, el := range cmplx.elems {
			if el.include {
				inPrintf(f, indent, "- required: [%v]\n", el.getName())
			}
		}

	default:
		inPrintf(f, indent, "type: object\n")
		if len(cmplx.attrs)+len(cmplx.elems) > 0 {
			inPrintf(f, indent, "properties:\n")
			if len(cmplx.attrs) > 0 {
				fmt.Printf("Doing attrs for complex %s\n", cmplx.name)
				writeAttrs(cmplx, f, ctxt, indent+tsz)
			}
			required := make([]string, 0)
			for _, el := range cmplx.elems {
				if el.include {
					writeElement(el, f, ctxt, indent+tsz)
					if el.minOccurs != 0 {
						required = append(required, el.getName())
					}
				}
			}
			if len(required) > 0 {
				inPrintf(f, indent, "required: %s\n", arrayString(required))
			}
		}
	}
	if !cmplx.anyFlag {
		inPrintf(f, indent, "additionalProperties: false\n")
	} else {
		inPrintf(f, indent, "# XSD allows 'any', so properties not restricted\n")
	}
}

func writeAttrs(attd attributed, f io.Writer, ctxt *context, indent int) []string {
	attrs := attd.getAttrs()
	required := []string{"value"}
	for _, attr := range attrs {
		if attr.required {
			required = append(required, "@"+attr.name)
		}
		inPrintf(f, indent, "'@%s':\n", attr.name)
		// atype must be either builtin or simple ...
		if _, ok := ctxt.simpleTypes[attr.atype]; ok {
			inPrintf(f, indent+tsz, "$ref: '#/components/schemas/%s'\n", attr.atype)
		} else {
			inPrintf(f, indent+tsz, "type: %s\n", attr.atype)
		}
		if attr.adefault != "" {
			inPrintf(f, indent+tsz, "default: \"%s\"\n", attr.adefault)
		}
		if attr.fixed != "" {
			// fmt.Fprinf(f, ",\n")
			inPrintf(f, indent+tsz, "# XML specified fixed value %s\n", attr.fixed)
		}
	}
	return required
}

// print an arbitrary thing with an indent
func inPrintf(f io.Writer, indent int, s string, v ...interface{}) (int, error) {
	var n1, n2 int
	var err error

	n1, err = fmt.Fprintf(f, "%*s", indent, "")
	if err == nil {
		n2, err = fmt.Fprintf(f, s, v...)
	}
	if err != nil {
		fmt.Printf("Write failed: %v\n", err)
	}
	return n1 + n2, err
}

func arrayString(strs []string) string {
	s := "["
	for _, val := range strs {
		s += fmt.Sprintf("'%s',", val)
	}
	s = s[:len(s)-1]
	s += "]"
	return s
}

func detab(str string) string {
	spc := strings.Repeat(" ", tsz)
	return strings.ReplaceAll(str, "\t", spc)
}
