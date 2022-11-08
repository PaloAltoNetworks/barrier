package gotest

import (
	"sort"
	"strings"
)

// suitesMap organizes suites in a map
type suitesMap map[string]*suiteInfo

func (s suitesMap) sorted() (out []*suiteInfo) {

	for _, t := range s {
		out = append(out, t)
	}

	sort.Slice(out, func(i int, j int) bool {
		return strings.Compare(out[i].Name, out[j].Name) == -1
	})

	return out
}

// testsMap organizes tests in a map
type testsMap map[string]Test

func (s testsMap) sorted() (out []Test) {

	for _, t := range s {
		out = append(out, t)
	}

	sort.Slice(out, func(i int, j int) bool {
		return strings.Compare(out[i].Name, out[j].Name) == -1
	})

	return out
}
