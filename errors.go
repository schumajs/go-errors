/*
  Copyright 2016 Jens Schumann <schumajs@gmail.com>

  Use of this source code is governed by the MIT license that can be found in
  the LICENSE file.
*/

package errors

import (
	"fmt"
	"log"
	"os"

	"github.com/go-errors/errors"
)

var IllegalArgumentError = New("illegal argument")

var Logger = log.New(os.Stderr, "[error] ", 0)

type Error struct {
	Err    error
	suffix string
}

func (e *Error) Error() string {
	msg := e.Err.Error()

	if e.suffix != "" {
		msg = fmt.Sprintf("%v: %v", msg, e.suffix)
	}

	return msg
}

func (e *Error) ErrorStack() string {
	return fmt.Sprintf(
		"%v %v\n%v",
		e.Err.(*errors.Error).TypeName(),
		e.Error(),
		string(e.Err.(*errors.Error).Stack()))
}

func New(errOrStr interface{}) *Error {
	return &Error{
		Err: errors.Wrap(errOrStr, 1),
	}
}

func NewSuffix(errOrStr interface{}, suffix string) *Error {
	return &Error{
		Err:    errors.Wrap(errOrStr, 1),
		suffix: suffix,
	}
}

func Errorf(format string, a ...interface{}) *Error {
	return &Error{
		Err: errors.Wrap(fmt.Errorf(format, a...), 1),
	}
}

func Wrap(errOrStr interface{}, skip ...int) *Error {
	var err error

	switch errOrStr := errOrStr.(type) {
	case *Error:
		return errOrStr
	case error:
		err = errOrStr
	default:
		err = fmt.Errorf("%v", errOrStr)
	}

	if len(skip) == 0 {
		err = errors.Wrap(err, 1)
	} else {
		err = errors.Wrap(err, skip[0])
	}

	return &Error{
		Err: err,
	}
}

func Is(err error, original error) bool {
	if err == original {
		return true
	}

	if err, ok := err.(*errors.Error); ok {
		return Is(err.Err, original)
	}

	if original, ok := original.(*errors.Error); ok {
		return Is(err, original.Err)
	}

	if err, ok := err.(*Error); ok {
		return Is(err.Err, original)
	}

	if original, ok := original.(*Error); ok {
		return Is(err, original.Err)
	}

	return false
}

func Print(errOrStr interface{}) {
	print(Wrap(errOrStr, 2))
}

func Printf(format string, a ...interface{}) {
	print(Wrap(fmt.Errorf(format, a...), 2))
}

func Fatal(errOrStr interface{}) {
	fatal(Wrap(errOrStr, 2))
}

func Fatalf(format string, a ...interface{}) {
	fatal(Wrap(fmt.Errorf(format, a...), 2))
}

func print(err error) {
	switch e := err.(type) {
	case *Error:
		Logger.Print(e.ErrorStack())
	default:
		panic("should never happen")
	}
}

func fatal(err error) {
	switch e := err.(type) {
	case *Error:
		Logger.Fatal(e.ErrorStack())
	default:
		panic("should never happen")
	}
}
