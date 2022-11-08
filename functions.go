package barrier

import (
	"context"
)

// SetupSuiteFunction is the type of function that can be run from a suite setup.
// The returned data will be available to the main test function using TestInfo.SuiteSetupInfo() function.
// The returned function will be run at the end of the suite.
//
// If SetupSuiteFunction returns an error, the entire suite of test is stopped.
type SetupSuiteFunction func(context.Context, SuiteInfo) (interface{}, TearDownFunction, error)

// A TestFunction is the type of a function that is run by a Test.
type TestFunction func(context.Context, TestInfo) error

// SetupFunction is the type of function that can be run a a test setup.
// The returned data will be available to the main test function using TestInfo.SetupInfo() function.
// The returned function will be run at the end of the test.
//
// If SetupFunction returns an error, the entire suite of test is stopped.
type SetupFunction func(context.Context, TestInfo) (interface{}, TearDownFunction, error)

// A TearDownFunction is the type of function returned by a SetupFunction or SetupSuiteFunction.
type TearDownFunction func()
