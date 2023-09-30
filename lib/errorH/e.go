package errorH

import "fmt"

func Wrap(message string, err error) error {
	return fmt.Errorf(message, err)
}

func WrapIfErr(message string, err error) error {
	if err != nil {
		return Wrap(message, err)
	}

	return nil
}
