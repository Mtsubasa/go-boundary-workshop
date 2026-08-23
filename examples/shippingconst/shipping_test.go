package shippingconst

import "testing"

func TestShippingFee(t *testing.T) {
	tests := []struct {
		input int
		want  int
	}{
		{input: 4999, want: 500},
		{input: 5001, want: 0},
	}

	for _, tt := range tests {
		if got := ShippingFee(tt.input); got != tt.want {
			t.Errorf("ShippingFee(%d) = %d, want %d", tt.input, got, tt.want)
		}
	}
}
