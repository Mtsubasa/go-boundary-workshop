package shipping

func ShippingFee(total int) int {
	if total < 5000 {
		return 500
	}
	return 0
}

func Discount(total int) int {
	if total < 100 {
		return 0
	}
	return 10
}
