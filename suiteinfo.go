package barrier

import (
	"fmt"
	"hash/fnv"
	"io"
)

// suiteInfo is runtime information for the suite
type suiteInfo struct {
	Name        string
	Description string
	Setup       SetupSuiteFunction
	tests       testsMap
	writer      io.Writer
	data        interface{}
	stash       interface{}
}

// SuiteInfo is the interface for the test writer
type SuiteInfo interface {
	RegisterTest(t Test)
	// SetupInfo provides the user data
	SetupInfo() interface{}
	// Stash provides the callers data that may store configs/runtime they need
	Stash() interface{}
	// Write performs a write
	Write(p []byte) (n int, err error)
}

// RegisterTest register a test in the main suite.
func (s *suiteInfo) RegisterTest(t Test) {

	if t.Name == "" {
		panic("test is missing name")
	}

	if t.Description == "" {
		panic("test is missing description")
	}

	if s.Name != defaultSuiteName {
		if t.SuiteName != "" {
			panic("suite name in a test can not be provided while registering it as a part of a suite")
		}
		t.SuiteName = s.Name
	}

	if t.Author == "" {
		panic("test is missing author")
	}

	if t.Function == nil {
		panic("test is missing function")
	}

	if len(t.Tags) == 0 {
		panic("test is missing tags")
	}

	h := fnv.New32()
	if _, err := h.Write([]byte(s.Name + s.Description + t.Name + t.Description + t.Author)); err != nil {
		panic(err)
	}
	t.id = fmt.Sprintf("%x", h.Sum32())

	if _, ok := s.tests[t.Name]; ok {
		panic("a test of the same name was previously registered: " + t.Name)
	}

	s.tests[t.Name] = t
}

// Stash provides the callers data that may store configs/runtime they need
func (s *suiteInfo) Stash() interface{} {
	return s.stash
}

// SetupInfo returns the eventual object stored by the Setup function.
func (s *suiteInfo) SetupInfo() interface{} {
	return s.data
}

// Write performs a write
func (s *suiteInfo) Write(p []byte) (n int, err error) {
	return s.writer.Write(p)
}

func (s *suiteInfo) testsWithArgs(verbose, matchAll bool, tags []string) *suiteInfo {

	ts := testsMap{}

	if verbose {
		fmt.Println("Running Tests:")
	}

	for _, t := range s.tests {

		if !t.matchTags(tags, matchAll) {
			continue
		}

		if verbose {
			fmt.Println(" - " + t.Name)
		}

		ts[t.Name] = t
	}

	if verbose && len(ts) == 0 {
		fmt.Println("No matching tests found.")
	}

	s.tests = ts
	return s
}

func (s *suiteInfo) testsWithIDs(verbose bool, ids []string) *suiteInfo {
	if len(ids) == 0 {
		return s
	}

	ts := testsMap{}

	if verbose {
		fmt.Println("Running Tests:")
	}

	for _, t := range s.tests {
		for _, id := range ids {
			if id == t.id {

				if verbose {
					fmt.Println(" - " + t.Name)
				}

				ts[t.Name] = t
			}
		}
	}

	if verbose && len(ts) == 0 {
		fmt.Println("No matching tests found.")
	}

	s.tests = ts
	return s
}

func (s *suiteInfo) String() string {
	return fmt.Sprintf(`suite name : %s
suite desc : %s
`, s.Name, s.Description)
}

func (s *suiteInfo) listTests() {

	fmt.Printf("%s\n", s)
	for _, test := range s.tests.sorted() {
		fmt.Printf("%s\n", test)
	}
}
