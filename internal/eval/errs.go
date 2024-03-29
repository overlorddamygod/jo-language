package eval

import (
	"errors"
	"fmt"
)

var (
	ErrUnexpected       = errors.New("UNEXPECTED ERROR")
	ErrIndexOutofBound  = errors.New("index out of bound")
	ErrArgLengthLess    = errors.New("arg length less than params length")
	ErrArgLengthGreater = errors.New("arg length greater than params length")
	ErrParseInt         = errors.New("cannot parse to int")
	ErrParseFloat       = errors.New("cannot parse to float")

	ErrBreak    = errors.New("Statement:Break")
	ErrContinue = errors.New("Statement:Continue")
	ErrThrow    = errors.New("Statement:Throw")
)

func ErrArgumentType(Type string) error {
	return fmt.Errorf("argument must be of type %s", Type)
}

func ErrArgumentLength(length int) error {
	return fmt.Errorf("argument length must be %d", length)
}

func ErrNoMethod(name string, Type string) error {
	return fmt.Errorf("no method named `%s` on type `%s`", name, Type)
}

func ErrReference(msg string) error {
	return fmt.Errorf("ReferenceError: %s", msg)
}
