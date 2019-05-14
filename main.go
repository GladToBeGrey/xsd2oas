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

// xsd2oas project main.go
// main function for xsd2oas

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	var exf, maskf *os.File
	// license notice
	fmt.Println("xsd2oas Copyright (C) 2019  Tom Hay")
	fmt.Println("This program comes with ABSOLUTELY NO WARRANTY")
	fmt.Println("This is free software, and you are welcome to redistribute it")
	fmt.Println("under certain conditions; see COPYING.txt for details.")

	// initialise
	ctxt := newContext()
	cmdLineParse(&ctxt)

	// open the input file
	fname := ctxt.inFile
	inf, err := os.Open(fname)
	if err != nil {
		fmt.Printf("File %v open err %v", fname, err)
		os.Exit(2)
	}
	defer inf.Close()

	// open the output file
	fname = ctxt.outFile
	outf, err := os.Create(fname)
	if err != nil {
		fmt.Printf("File %v open err %v", fname, err)
		os.Exit(2)
	}
	defer outf.Close()

	// open the mask file
	if ctxt.maskFile != "" {
		ctxt.mask = true
		fname := ctxt.maskFile
		maskf, err := os.Open(fname)
		if err != nil {
			fmt.Printf("File %v open err %v", fname, err)
			os.Exit(2)
		}
		scanner := bufio.NewScanner(maskf)
		for scanner.Scan() {
			ctxt.maskLines = append(ctxt.maskLines, scanner.Text())
		}
		if scanner.Err() != nil {
			fmt.Printf("File %v scan err %v", fname, scanner.Err())
			os.Exit(2)
		}
		fmt.Printf("File %v scanned OK - %v lines\n", fname, len(ctxt.maskLines))
	}
	defer maskf.Close()

	// open the example file
	if ctxt.exFile != "" {
		fname := ctxt.exFile
		f, err := os.Create(fname)
		if err != nil {
			fmt.Printf("File %v open err %v", fname, err)
			os.Exit(2)
		}
		exf = f
		fmt.Printf("File %v open OK\n", fname)
	}
	defer exf.Close()

	parseXml(inf, &ctxt)
	tagInclude(&ctxt)
	writeYaml(outf, &ctxt)
	if ctxt.exFile != "" {
		writeExample(exf, &ctxt)
	}
}
