package gotest

// Suite is a collection of tests
type Suite struct {
	Name        string
	Description string
	Setup       SetupSuiteFunction
}
