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

// cmdLineParse
// parse and validate the command line arguments

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func cmdLineParse(ctxt *context) {
	inFilePtr := flag.String("in", "", "input file name")
	outFilePtr := flag.String("out", "", "output file name")
	exFilePtr := flag.String("ex", "", "example file name")
	maskFilePtr := flag.String("mask", "", "mask file name")
	licPtr := flag.Bool("lic", false, "print license info")
	verbPtr := flag.Bool("v", false, "verbose")

	flag.Parse()

	if *inFilePtr == "" || *outFilePtr == "" {
		fmt.Printf(
			`Usage: %s -in xsdfile -out yamlfile
Optional parameters:
-mask maskfile
-ex examplefile`, filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	ctxt.inFile = *inFilePtr
	ctxt.inFileBase = filepath.Base(ctxt.inFile)
	ctxt.outFile = *outFilePtr
	ctxt.outFileBase = filepath.Base(ctxt.outFile)
	ctxt.printLicense = *licPtr
	ctxt.verbose = *verbPtr

	if *exFilePtr != "" {
		ctxt.exFile = *exFilePtr
		ctxt.exFileBase = filepath.Base(ctxt.exFile)
	}
	if *maskFilePtr != "" {
		ctxt.maskFile = *maskFilePtr
		ctxt.maskFileBase = filepath.Base(ctxt.maskFile)
	}
}
