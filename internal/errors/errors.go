package errors

import (
	"fmt"
)

func New[T ~string](format T, values ...any) Error[T] {
	return newError(format, values...)
}

func As[T ~string](err error) Error[T] {
	return err.(Error[T])
}

type Error[T ~string] string

func newError[T ~string](format T, values ...any) Error[T] {
	var err string
	if len(values) != 0 {
		// Do not call fmt.Sprintf() if not necessary.
		// Major performance improvement.
		// Only necessary if there are any values.
		err = fmt.Sprintf(string(format), values...)
	} else {
		err = string(format)
	}

	return Error[T](err)
}

func (self Error[T]) Unpack() T {
	return T(self)
}

func (self Error[T]) Error() string {
	return string(self)
}
