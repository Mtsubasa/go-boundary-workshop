package shipping

const freeShippingBoundary = 4000 + 1000

func ShippingFee(total int) int {
	if total < freeShippingBoundary {
		return 500
	}
	return 0
}
