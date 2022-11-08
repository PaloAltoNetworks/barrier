package barrier

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/buger/goterm"
	"github.com/smartystreets/goconvey/convey"
)

type assertionError struct {
	msg         string
	description string
	Expected    interface{}
	Actual      interface{}
}

func newassertionError(msg string) assertionError {
	return assertionError{
		msg: msg,
	}
}

func (e assertionError) Error() string {
	if e.Expected != nil && e.Actual != nil {
		return goterm.Color(fmt.Sprintf("[FAIL] %s: expected: '%s', actual '%s'", e.msg, e.Expected, e.Actual), goterm.RED)
	}
	return goterm.Color(fmt.Sprintf("[FAIL] %s: %s", e.msg, e.description), goterm.RED)
}

// Assert can use goconvey function to perform an assertion.
func Assert(t TestInfo, message string, actual interface{}, f func(interface{}, ...interface{}) string, expected ...interface{}) {

	if msg := f(actual, expected...); msg != "" {

		r := newassertionError(message)

		if err := json.Unmarshal([]byte(msg), &r); err != nil {
			r.description = strings.Replace(strings.Replace(msg, "\n", ", ", -1), "\t", " ", -1)
		}

		panic(r)
	}

	fmt.Fprint(t, goterm.Color(fmt.Sprintf("- [PASS] %s", message), goterm.GREEN)) // nolint
	fmt.Fprintln(t)                                                                // nolint
}

// Step runs a particular step.
func Step(t TestInfo, name string, step func() error) {

	start := time.Now()
	fmt.Fprintf(t, "%s\n", name) // nolint
	if err := step(); err != nil {
		fmt.Fprintf(t, "%s\n", goterm.Color(fmt.Sprintf("took: %s", time.Since(start).Round(time.Millisecond)), goterm.BLUE)) // nolint
		Assert(t, "step should not return any error", err, convey.ShouldBeNil)
	}

	fmt.Fprintf(t, "%s\n\n", goterm.Color(fmt.Sprintf("took: %s", time.Since(start).Round(time.Millisecond)), goterm.BLUE)) // nolint
}
