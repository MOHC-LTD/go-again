package retry_test

import (
	retry "belgium/MOHC-LTD/go-again"
	"errors"
	"fmt"
	"testing"
	"time"
)

// Func test back off handler retries
func TestRetriesTwiceBeforeSuccessDoesNotReportError(t *testing.T) {
	// Given - a function that fails once before success

	failingFunction, getRunCount := failNTimes(1, nil)

	// When - try with max retry of 1 (so that the second attempt succeeds)

	err := retry.Configure().
		FirstRetryDelay(1 * time.Millisecond).
		WithMaxRetries(1).
		Try(failingFunction)

	// Then - expect no error
	if err != nil {
		t.Errorf("\nReceived unexpected error: \n\n%v", err)
	}

	// Then - expect to have been called twice
	if getRunCount() != 2 {
		t.Errorf(
			"\nReceived unexpected error: \n\n%v",
			fmt.Sprintf("Expected to be run twice but was run %v times", getRunCount()),
		)
	}
}

// Func test that the error is reported if we don't get a success after running out of retries
func TestErrorReportedAfterMultipleFailures(t *testing.T) {
	// Given - a function that fails twice before success

	failingFunction, _ := failNTimes(2, nil)

	// When - try with max retry of 1 (so that it never succeeds)

	err := retry.Configure().
		FirstRetryDelay(1 * time.Millisecond).
		WithMaxRetries(1).
		Try(failingFunction)

	// Then - should get an error

	if err == nil {
		t.Errorf(
			"\nReceived unexpected error: \n\n%v",
			"We should of had a error due to too many failed attempts and we didn't get one",
		)
	}
}

// Func Test that an error not to retry on can be specified to fail fast
func TestFailFastOnSpecificErrors(t *testing.T) {
	// Given - a function that fails once before success and returns a specified error

	errorNotToRetryOn := errors.New("error not to retry on")

	failingFunction, getRunCount := failNTimes(
		1,
		errorNotToRetryOn,
	)

	// When - try with max retry of 1 and error not to retry on

	err := retry.Configure().
		FirstRetryDelay(1 * time.Millisecond).
		WithMaxRetries(1).
		DontRetryOn(errorNotToRetryOn).
		Try(failingFunction)

	// Then - expect get error we specified not to retry on

	if err != errorNotToRetryOn {
		t.Errorf("\nReceived unexpected error: \n\n%v", "Didn't get the error we expected")
	}

	// Then - expect function not to have been retried

	if getRunCount() > 1 {
		t.Errorf("\nReceived unexpected error: \n\n%v", "We ran more than once and tried to retry")
	}
}

// Func Use a fallback function if all retry attempts fail
func TestFallBackFunctionIsCalled(t *testing.T) {
	// Given - a function that fails twice before success

	failingFunction, _ := failNTimes(2, nil)

	// Given - a fallback function that doesn't fail

	fallback, getFallbackRunCount := failNTimes(0, nil)

	// When - try with max retry of 1 (so that it never succeeds) and a fallback function

	retry.Configure().
		FirstRetryDelay(1 * time.Millisecond).
		WithMaxRetries(1).
		WithFallBack(fallback).
		Try(failingFunction)

	// Then - fallback should have been called

	if getFallbackRunCount() == 0 {
		t.Errorf("\nReceived unexpected error: \n\n%v", "fallback was not called")
	}
}

// Func Use a fallback function if all retry attempts fail
func TestErrorHandlerFunctionIsCalled(t *testing.T) {
	// Given - a function that fails once before success

	failingFunction, _ := failNTimes(1, nil)

	// Given - an error handler function

	errorHandlerCalled := false
	var errorHandlerErrorReceived error
	errorHandler := func(err error) {
		errorHandlerCalled = true

		errorHandlerErrorReceived = err
	}

	// When - try with max retry 1 (so that it succeeds on second attempt) and an error handler

	retry.Configure().
		FirstRetryDelay(1 * time.Millisecond).
		WithMaxRetries(1).
		OnAttemptFailure(errorHandler).
		Try(failingFunction)

	// Then - expect error handler to have been called

	if !errorHandlerCalled {
		t.Errorf("\nReceived unexpected error: \n\n%v", "error handler was not called")
	}

	// Then - expect error handler to have been called with an error

	if errorHandlerErrorReceived == nil {
		t.Errorf("\nReceived unexpected error: \n\n%v", "error handler was not called with an error")
	}
}

// failNTimes creates a function that will fail the first n calls. The second return is a function to ask for the number of times it was actually called.
func failNTimes(n int, errorToReturn error) (func() error, func() int) {
	runCount := 0

	failingFunction := func() error {
		runCount++

		if runCount <= n {
			if errorToReturn != nil {
				return errorToReturn
			}

			return fmt.Errorf("attempt number %v", runCount)
		}

		return nil
	}

	getRunCount := func() int {
		return runCount
	}

	return failingFunction, getRunCount
}
