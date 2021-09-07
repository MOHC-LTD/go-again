package retry

import "time"

// Configure sets up an initial specification with some defaults. Use the builder methods to configure as required i.e retry.Configure().WithMaxRetries(10)
func Configure() Specification {
	return Specification{maxRetries: 1, startDelay: 1 * time.Second}
}

// WithMaxRetries sets up a specification with a maximum number of retries. The total attempts will be 1 + maxRetries
func (s Specification) WithMaxRetries(retries uint64) Specification {
	s.maxRetries = retries
	return s
}

// FirstRetryDelay sets up a specification with a given first retry delay, this is the gap between the initial attempt and the first retry
func (s Specification) FirstRetryDelay(t time.Duration) Specification {
	s.startDelay = t
	return s
}

// DontRetryOn sets up a specification that doesn't retry on a given error or list of errors
func (s Specification) DontRetryOn(err ...error) Specification {
	s.dontRetryErrorSet = append(s.dontRetryErrorSet, err...)
	return s
}
