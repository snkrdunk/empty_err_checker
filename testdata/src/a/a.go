package a

import (
	"errors"
	"math/rand"
	"time"
)

func a() {
	rand.Seed(time.Now().UnixNano())
	f := func() error {
		_, err := do()
		if err != nil {
			return err
		}
		isInvalid := checkSomething()
		if isInvalid {
			return err // want "returned error is not checked."
		}
		return nil
	}
	if err := f(); err != nil {
		panic(err)
	}

	if err := returnEmptyErr(); err != nil {
		panic(err)
	}
}

func returnEmptyErr() error {
	_, err := do()
	if err != nil {
		return err
	}
	isInvalid := checkSomething()
	if isInvalid {
		return err // want "returned error is not checked."
	}
	return nil
}

func do() (int, error) {
	i := rand.Intn(10)
	if i == 0 {
		return i, errors.New("error")
	} else {
		return i, nil
	}
}

func checkSomething() bool {
	return true
}
