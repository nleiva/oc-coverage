/*
Disregard: Experimenting with different diff libraries
*/

package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// Testing diff libraries
func main() {

	b1, err := ioutil.ReadFile("../output/initial.cfg")
	if err != nil {
		log.Fatalf("Ooops, couldn't read the file")
	}
	s1 := string(b1)
	b2, err := ioutil.ReadFile("../output/openconfig.cfg")
	if err != nil {
		log.Fatalf("Ooops, couldn't read the file")
	}
	s2 := string(b2)

	// https://github.com/sergi/go-diff
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(s1, s2, true)
	fmt.Println(dmp.DiffPrettyText(diffs))

	// https://github.com/pmezard/go-difflib
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(s1),
		B:        difflib.SplitLines(s2),
		FromFile: "Original",
		ToFile:   "Current",
		Context:  3,
	}
	text, _ := difflib.GetUnifiedDiffString(diff)
	fmt.Printf(text)

}
