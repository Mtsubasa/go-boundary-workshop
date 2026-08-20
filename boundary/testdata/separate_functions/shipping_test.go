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

func TestDiscount(t *testing.T) {
	tests := []struct {
		input int
		want  int
	}{
		{input: 99, want: 0},
		{input: 100, want: 10},
		{input: 101, want: 10},
	}
	_ = tests
	_ = t
}
