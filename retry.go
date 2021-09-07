package retry

import (
	"context"
	"time"

	"github.com/pkg/errors"
	. "github.com/sethvargo/go-retry"
)

// Tryable function to retry, if an error is returned it will retry as specified by the config, if no error is returned, no further attempts will be made.
type Tryable func() error

// ErrorFunc function to intercept errors produced during attempts made. It will report every error captured during an attempt.
type ErrorFunc func(err error)

// Specification is the configuration specifying how to retry a function.
/*
In the case of the proceeding example, assuming that `failingFunction` always fails.

- The initial call of `failingFunction` fails.
- The `errorHandler` is called with the resulting error.
- `failingFunction` is called for a second time.
- The second call of `failingFunction` fails.
- The `errorHandler` is called with the resulting error.
- The `fallback` function is called.
- If the fallback function fails then an error is returned, else the error will be `nil`.

Example:
	err := retry.Configure().
		goFirstRetryDelay(1 * time.Millisecond). 	// After the first attempt delay this amount of time.
		WithMaxRetries(1). 							// Keep retrying this number of times.
		OnAttemptFailure(errorHandler). 			// Every time a failure happens report it here.
		WithFallBack(fallback). 					// If everything still failed then do this.
		Try(failingFunction) 						// Terminal operator for the function to retry.
*/
type Specification struct {
	fallback            Tryable       // The fallback function if all retries fail - Note: this can still produce an error.
	attemptErrorHandler ErrorFunc     // The error handler to capture errors that occur during the retry - this allows you to log each attempted failure even if you retry.
	maxRetries          uint64        // The maximum number of times to retry, note the total attempts made will be initial + retries.
	startDelay          time.Duration // The duration to start the exponential back off from, for example first run, wait this delay period then make first retry attempt.
	dontRetryErrorSet   []error       // A list of errors that we don't want to retry on, for example a 404 would not be worth retrying as you know it is missing.
}

// WithFallBack configures the retry to call this specified function when all attempts have failed.
func (s Specification) WithFallBack(fallback Tryable) Specification {
	s.fallback = fallback
	return s
}

// OnAttemptFailure configures an error handler that will receive any errors caught whilst retrying.
func (s Specification) OnAttemptFailure(errorHandler ErrorFunc) Specification {
	s.attemptErrorHandler = errorHandler
	return s
}

// Try configures a function to be run multiple times when a failure occurs until a configured limit.
func (s Specification) Try(functionToRetry Tryable) error {
	fib, err := NewFibonacci(s.startDelay)

	if err != nil {
		return errors.Wrap(err, "fibonacci start delay must be greater than 0, no attempt was made")
	}

	retrySpec := WithMaxRetries(s.maxRetries, fib)

	retryLibFunction := func(ctx context.Context) error {
		err := functionToRetry()

		if err == nil {
			return nil
		}

		for _, errType := range s.dontRetryErrorSet {
			if errors.Cause(err) == errType {
				return err
			}
		}

		if s.attemptErrorHandler != nil {
			s.attemptErrorHandler(err)
		}
		return RetryableError(err)
	}

	err = Do(context.Background(), retrySpec, retryLibFunction)

	if err != nil && s.fallback != nil {
		err = errors.Wrap(s.fallback(), "error in fallback")
	}

	return err
}
