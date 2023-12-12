# empty_err_checker

[![Test](https://github.com/snkrdunk/empty_err_checker/actions/workflows/test.yml/badge.svg)](https://github.com/snkrdunk/empty_err_checker/actions/workflows/test.yml)

empty_err_checker is checking whether the return value 'err' is nil.

example
```go
func inValidErrChecker() error {
	var err error
	isValid := isValid()
	if !isValid {
		return err // report this return err as invalid
	}
	return nil
}
```

Check [the test code](https://github.com/snkrdunk/empty_err_checker/blob/main/testdata/src/a/a.go) for detailed detection examples.

# Installation

```console
go install github.com/snkrdunk/empty_err_checker/cmd/empty_err_checker@latest
```
