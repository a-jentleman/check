// Package check provides utilities for checking the values of variables against
// common conditions returning specific error types.
package check

import (
	"cmp"
	"fmt"
)

// OutOfRangeError represents an error that occurred when a value was outside an expected range.
type OutOfRangeError string

func (o OutOfRangeError) Error() string {
	return string(o)
}

// Is reports whether err is of type OutOfRangeError in accordance with [errors.Is].
func (o OutOfRangeError) Is(err error) bool {
	_, ok := err.(OutOfRangeError) // per errors.Is documentation, this should only do a shallow comparison (and not unwrap err)
	return ok
}

// InRange returns an [OutOfRangeError] if !(min <= actual <= max).
// InRange panics if max < min or if min or max are NaN.
func InRange[T cmp.Ordered](actual, min, max T) error {
	switch {
	case max != max || min != min:
		panic("check: range boundary is NaN")
	case max < min:
		panic("check: max < min")
	case cmp.Less(actual, min) || cmp.Less(max, actual):
		return OutOfRangeError(fmt.Sprintf("out-of-range (expected %v <= actual <= %v, but actual=%v)", min, max, actual))
	default:
		return nil
	}
}
