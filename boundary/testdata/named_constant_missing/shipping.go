package shipping

const freeShippingBoundary = 5000

func ShippingFee(total int) int {
	if total < freeShippingBoundary {
		return 500
	}
	return 0
}
