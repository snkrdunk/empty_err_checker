package a

import (
	"errors"
	"fmt"
)

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
		err = fmt.Errorf("error")
		return err
	}
	return nil
}

func validErrChecker3() error {
	err := verifySomething()
	if err != nil {
		if !isValid() {
			return err
		}
		return err
	}
	return nil
}

func validErrChecker4() error {
	err := verifySomething()
	if err != nil {
		return err
	}
	if err == nil {
		err := errors.New("error")
		if !isValid() {
			return err
		}
		return err
	}
	return nil
}

func validErrChecker5() error {
	err := verifySomething()
	if err != nil && !isValid() || true {
		return err
	}
	return nil
}

func validErrChecker6() error {
	err := verifySomething()
	if !isValid() && err != nil || true {
		return err
	}
	return nil
}

func validErrChecker7() error {
	err := verifySomething()
	if !isValid() && true || err != nil {
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

func inValidErrChecker2() error {
	err := verifySomething()
	if err != nil {
		return err
	}
	if err == nil {
		if !isValid() {
			return err // want "returned error is not checked."
		}
		return err // want "returned error is not checked."
	}
	return nil
}

func invalidErrChecker3() error {
	err := verifySomething()
	isValid := isValid()
	if isValid {
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
