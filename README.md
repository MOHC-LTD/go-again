# retry

Package `retry` provides functionality for retrying code.

## Contents

- [💻 Installation](#installation)
- [✨‍ Usage](#usage)
- [⚒️ How to contribute](#how-to-contribute)

## Installation

Install the module using

```sh
go get -u github.com/MOHC-LTD/go-again
```

## Usage

Use this package to retry code.

```go
type RetrySomething struct {
	retry retry.Specification
	repo  Repository
}

func Retry(retries uint64, repo Repository) RetrySomething {
	return RetrySomething{retry.Configure().WithMaxRetries(retries), repo}
}

func (retrySomething RetrySomething) Do() error {
	err := retrySomething.retry.
		OnAttemptFailure(func(err error) {
			// The attempt failed...
		}).
		Try(func() error {
			var err error
			err = retrySomething.repo.Do()
			if err != nil {
				return err
			}

			return nil
		})

	return err
}
```

## How to contribute

Something missing or not working as expected? See our [contribution guide](./CONTRIBUTING.md).
