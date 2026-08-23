package shipping

import "testing"

func TestShippingFee(t *testing.T) {
	tests := []struct {
		input int
		want  int
	}{
		{input: 4999, want: 500},
		{input: 5001, want: 0},
	}
	_ = tests
	_ = t
}
