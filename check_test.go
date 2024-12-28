package check_test

import (
	"errors"
	"math"
	"testing"

	"github.com/a-jentleman/check"
	"github.com/stretchr/testify/assert"
)

func TestInRange(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		tests := []struct {
			name     string
			actual   int
			min      int
			max      int
			wantErr  bool
			panicErr bool
		}{
			{"in-range", 5, 1, 10, false, false},
			{"at-lower-bound", 1, 1, 10, false, false},
			{"at-upper-bound", 10, 1, 10, false, false},
			{"below-lower-bound", 0, 1, 10, true, false},
			{"above-upper-bound", 11, 1, 10, true, false},
			{"max-less-than-min", 5, 10, 1, false, true},
			{"actual-equals-min-and-max", 5, 5, 5, false, false},
			{"min-max-equal-actual-not-equal", 4, 5, 5, true, false},
			{"negative-range-valid", -3, -5, -1, false, false},
			{"negative-range-out-of-lower-bound", -6, -5, -1, true, false},
			{"negative-range-out-of-upper-bound", 0, -5, -1, true, false},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				defer func() {
					assert.Equal(t, tt.panicErr, recover() != nil)
				}()

				actualErr := check.InRange[int](tt.actual, tt.min, tt.max)
				if tt.wantErr {
					assert.ErrorIs(t, actualErr, check.OutOfRangeError(""))
					return
				}
			})
		}
	})

	t.Run("float", func(t *testing.T) {
		tests := []struct {
			name     string
			actual   float64
			min      float64
			max      float64
			wantErr  bool
			panicErr bool
		}{
			{"in-range", 5, 1, 10, false, false},
			{"at-lower-bound", 1, 1, 10, false, false},
			{"at-upper-bound", 10, 1, 10, false, false},
			{"below-lower-bound", 0, 1, 10, true, false},
			{"above-upper-bound", 11, 1, 10, true, false},
			{"max-less-than-min", 5, 10, 1, false, true},
			{"actual-equals-min-and-max", 5, 5, 5, false, false},
			{"min-max-equal-actual-not-equal", 4, 5, 5, true, false},
			{"negative-range-valid", -3, -5, -1, false, false},
			{"negative-range-out-of-lower-bound", -6, -5, -1, true, false},
			{"negative-range-out-of-upper-bound", 0, -5, -1, true, false},
			{"lower-bound-NaN", 0, math.NaN(), 1, false, true},
			{"upper-bound-NaN", 0, 0, math.NaN(), false, true},
			{"actual-NaN", math.NaN(), 0, 1, true, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				defer func() {
					assert.Equal(t, tt.panicErr, recover() != nil)
				}()

				actualErr := check.InRange[float64](tt.actual, tt.min, tt.max)
				if tt.wantErr {
					assert.ErrorIs(t, actualErr, check.OutOfRangeError(""))
					return
				}
				assert.NoError(t, actualErr)
			})
		}
	})
}

func TestIndex(t *testing.T) {
	tests := []struct {
		name    string
		index   int
		slice   []int
		wantErr bool
	}{
		{"valid-index", 2, []int{1, 2, 3}, false},
		{"index-0-valid", 0, []int{1, 2, 3}, false},
		{"index-out-of-range-upper", 3, []int{1, 2, 3}, true},
		{"index-out-of-range-lower", -1, []int{1, 2, 3}, true},
		{"index-at-upper-bound-valid", 2, []int{1, 2, 3}, false},
		{"empty-slice", 0, []int{}, true},
		{"index-out-of-range-empty-slice", 1, []int{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := check.Index(tt.index, tt.slice)
			if tt.wantErr {
				assert.ErrorIs(t, actual, check.OutOfRangeError(""))
				return
			}
			assert.NoError(t, actual)
		})
	}
}

type NotOutOfRangeError string

func (e NotOutOfRangeError) Error() string {
	return string(e)
}

func TestOutOfRangeError_Is(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		target error
		want   bool
	}{
		{"same-types", check.OutOfRangeError("out of range"), check.OutOfRangeError("another error"), true},
		{"different-types", check.OutOfRangeError("out of range"), errors.New("not out of range"), false},
		{"different-string-types", check.OutOfRangeError("out of range"), NotOutOfRangeError("not out of range"), false},
		{"nil-target", check.OutOfRangeError("out of range"), nil, false},
		{"nil-error", nil, check.OutOfRangeError("out of range"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want {
				assert.ErrorIs(t, tt.err, tt.target)
				return
			}
			assert.NotErrorIs(t, tt.err, tt.target)
		})
	}
}

func TestOutOfRangeError_Error(t *testing.T) {
	tests := []struct {
		name     string
		inputErr check.OutOfRangeError
		expected string
	}{
		{"basic-error-message", check.OutOfRangeError("out of range"), "out of range"},
		{"empty-error-message", check.OutOfRangeError(""), ""},
		{"long-error-message", check.OutOfRangeError("this is a very long error message to test bounds"), "this is a very long error message to test bounds"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.inputErr.Error())
		})
	}
}
