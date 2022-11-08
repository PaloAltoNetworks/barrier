package barrier

import (
	"fmt"
)

const defaultSuiteName = "_default"

var mainSuites suitesMap

// RegisterTest register a test in the main suite.
func RegisterTest(t Test) {
	mainSuites[defaultSuiteName].RegisterTest(t)
}

// RegisterSuite registers a test suite.
func RegisterSuite(s Suite) SuiteInfo {

	if s.Name == "" {
		panic("suite is missing name")
	}

	if s.Description == "" {
		panic("suite is missing description")
	}

	if _, ok := mainSuites[s.Name]; ok {
		panic(fmt.Sprintf("suite already registered with name %s", s.Name))
	}

	si := &suiteInfo{
		Name:        s.Name,
		Description: s.Description,
		Setup:       s.Setup,
		tests:       testsMap{},
	}
	mainSuites[si.Name] = si
	return si
}

func init() {
	mainSuites = map[string]*suiteInfo{}
	RegisterSuite(Suite{
		Name:        defaultSuiteName,
		Description: defaultSuiteName,
	})
}
