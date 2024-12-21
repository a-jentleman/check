package check_test

import (
	"errors"
	"math"
	"testing"

	"github.com/a-jentleman/check"
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
					if r := recover(); r != nil {
						if !tt.panicErr {
							t.Errorf("unexpected panic: %v", r)
						}
					} else if tt.panicErr {
						t.Errorf("expected panic, but did not panic")
					}
				}()

				err := check.InRange[int](tt.actual, tt.min, tt.max)
				if (err != nil) != tt.wantErr {
					t.Errorf("got error = %v, want error = %v", err != nil, tt.wantErr)
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
					if r := recover(); r != nil {
						if !tt.panicErr {
							t.Errorf("unexpected panic: %v", r)
						}
					} else if tt.panicErr {
						t.Errorf("expected panic, but did not panic")
					}
				}()

				err := check.InRange[float64](tt.actual, tt.min, tt.max)
				if (err != nil) != tt.wantErr {
					t.Errorf("got error = %v, want error = %v", err != nil, tt.wantErr)
				}
			})
		}
	})
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
			if got := errors.Is(tt.err, tt.target); got != tt.want {
				t.Errorf("errors.Is(%v, %v) = %v, want %v", tt.err, tt.target, got, tt.want)
			}
		})
	}
}

func TestOutOfRangeError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  check.OutOfRangeError
		want string
	}{
		{"simple-error-message", check.OutOfRangeError("value out of range"), "value out of range"},
		{"empty-error-message", check.OutOfRangeError(""), ""},
		{"numeric-error-message", check.OutOfRangeError("123"), "123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
