package barrier

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gofrs/uuid"
)

type testRun struct {
	buildID  string
	ctx      context.Context
	test     Test
	testInfo TestInfo
	verbose  bool
}

type testResult struct {
	err       error
	reader    io.Reader
	duration  time.Duration
	test      Test
	iteration int
	stack     []byte
}

type testRunner struct {
	concurrent    int
	buildID       string
	resultsChan   chan testRun
	setupErrs     chan error
	skipTeardown  bool
	status        map[string]testRun
	stopOnFailure bool
	stress        int
	suite         *suiteInfo
	teardowns     chan TearDownFunction
	timeout       time.Duration
	verbose       bool
}

func newTestRunner(
	buildID string,
	suite *suiteInfo,
	timeout time.Duration,
	concurrent int,
	stress int,
	verbose bool,
	skipTeardown bool,
	stopOnFailure bool,
) *testRunner {
	return &testRunner{
		concurrent:    concurrent,
		resultsChan:   make(chan testRun, concurrent*stress),
		setupErrs:     make(chan error),
		skipTeardown:  skipTeardown,
		status:        map[string]testRun{},
		stopOnFailure: stopOnFailure,
		stress:        stress,
		suite:         suite,
		timeout:       timeout,
		verbose:       verbose,
		buildID:       buildID,
	}
}

func (r *testRunner) executeIteration(ctx context.Context, currTest testRun, results chan testResult) {

	sem := make(chan struct{}, r.concurrent)

	for i := 0; i < r.stress; i++ {

		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			return
		}

		go func(t testRun, iteration int) {
			var data interface{}
			var td TearDownFunction
			var err error
			var start time.Time

			buf := &bytes.Buffer{}

			defer func() { <-sem }()

			ti := testResult{
				test:      t.test,
				reader:    buf,
				iteration: iteration,
			}

			defer func() {

				defer func() { results <- ti }()

				// recover remote code.
				r := recover()
				ti.duration = time.Since(start)
				if r == nil {
					return
				}

				err, ok := r.(assertionError)
				if ok {
					ti.err = err
					return
				}

				ti.err = fmt.Errorf("unhandled panic: %s", r)
				ti.stack = debug.Stack()
			}()

			subTestInfo := TestInfo{
				data:           data,
				iteration:      iteration,
				testID:         uuid.Must(uuid.NewV4()).String(),
				timeOfLastStep: t.testInfo.timeOfLastStep,
				timeout:        r.timeout,
				writer:         buf,
				suite:          r.suite,
			}

			if t.test.Setup != nil {
				data, td, err = t.test.Setup(t.ctx, subTestInfo)
				if err != nil {
					printSetupError(t.test.id, t.test.SuiteName, t.test.Name, nil, err)
					ti.err = err
					return
				}
				subTestInfo.data = data

				defer func() {
					if r.skipTeardown {
						subTestInfo.Write([]byte("Teardown skipped.")) //nolint
					} else if td != nil {
						td()
					}
				}()
			}

			start = time.Now()
			ti.err = t.test.Function(ctx, subTestInfo)

		}(currTest, i)
	}
}

func (r *testRunner) execute(ctx context.Context) error {

	sem := make(chan struct{}, r.concurrent)
	done := make(chan struct{})
	stop := make(chan struct{})

	var wg sync.WaitGroup
	var err error

L:
	for _, test := range r.suite.tests.sorted() {

		wg.Add(1)

		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			return err
		case <-stop:
			break L
		}

		go func(run testRun) {

			buf := &bytes.Buffer{}
			hdr := &bytes.Buffer{}

			run.testInfo.writer = buf
			run.testInfo.header = hdr

			defer func() { wg.Done(); <-sem }()

			resultsCh := make(chan testResult)

			go r.executeIteration(ctx, run, resultsCh)

			var results []testResult

		L2:
			for {
				select {
				case res := <-resultsCh:
					results = append(results, res)

					if res.err != nil {
						err = res.err

						if r.stopOnFailure {
							appendResults(run, results, r.verbose)
							fmt.Println(hdr.String())
							fmt.Println(buf.String())
							close(stop)

							return
						}
					}

					if len(results) == r.stress {
						appendResults(run, results, r.verbose)
						break L2
					}
				case <-ctx.Done():
					break L2
				}
			}

			if hdr.String() != "" {
				fmt.Println(hdr.String())
			}
			if buf.String() != "" {
				fmt.Println(buf.String())
			}
		}(testRun{
			ctx:     ctx,
			buildID: r.buildID,
			test:    test,
			verbose: r.verbose,
			testInfo: TestInfo{
				timeOfLastStep: time.Now(),
				timeout:        r.timeout,
				suite:          r.suite,
			},
		})
	}

	go func() {
		defer close(done)
		wg.Wait()
	}()

	select {
	case <-done:
	case <-stop:
	}

	return err
}

func (r *testRunner) run(ctx context.Context, suite *suiteInfo) error {

	if suite.Setup != nil {

		buf := &bytes.Buffer{}

		suite.writer = buf

		data, td, err := suite.Setup(ctx, suite)
		if err != nil {
			printSetupError("Suite", suite.Name, "", nil, err)
			return err
		}
		suite.data = data

		if r.verbose && buf.String() != "" {
			fmt.Println(buf.String())
			buf = &bytes.Buffer{}
			suite.writer = buf
		}

		defer func() {
			if r.skipTeardown {
				suite.Write([]byte("Teardown skipped.")) //nolint
			} else if td != nil {
				td()
			}

			if r.verbose && buf.String() != "" {
				fmt.Println(buf.String())
			}
		}()
	}

	r.teardowns = make(chan TearDownFunction, len(suite.tests))
	if err := r.execute(ctx); err != nil {
		return fmt.Errorf("failed test(s). please check logs")
	}

	if ctx.Err() != nil {
		return fmt.Errorf("deadline exceeded. Try giving a higher time limit using --limit option (%s)", ctx.Err())
	}

	return nil
}
