package a

import "errors"

func main() {
	err := validErrChecker()
	if err != nil {
		panic(err)
	}

	err = validErrChecker2()
	if err != nil {
		panic(err)
	}

	err = inValidErrChecker()
	if err != nil {
		panic(err)
	}
}

func validErrChecker() error {
	err := verifySomething()
	if err != nil {
		return err
	}
	return nil
}

func validErrChecker2() error {
	var err error
	isValid := isValid()
	if !isValid {
		err = errors.New("error")
		return err
	}
	return nil
}

func inValidErrChecker() error {
	var err error
	isValid := isValid()
	if !isValid {
		return err // want "returned error is not checked."
	}
	return nil
}

func isValid() bool {
	return false
}

func verifySomething() error {
	return errors.New("error")
}
