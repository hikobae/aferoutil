package zip

import (
	"errors"
	"testing"
)

func TestEr(t *testing.T) {
	var (
		oldErr = errors.New("Old error")
		newErr = errors.New("New error")
	)

	testcases := []struct {
		init     error
		input    error
		expected error
	}{
		{nil, nil, nil},
		{nil, newErr, newErr},
		{oldErr, nil, oldErr},
		{oldErr, newErr, oldErr},
	}

	for _, tc := range testcases {
		err := tc.init
		f := func() error {
			return tc.input
		}
		er(f, &err)
		if err != tc.expected {
			t.Errorf("expected %v, but %v", tc.expected, err)
		}
	}
}
