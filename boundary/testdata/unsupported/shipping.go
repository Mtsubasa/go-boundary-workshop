package shipping

var boundary = 5000

func NamedConstant(total int) int {
	if total < boundary {
		return 500
	}
	return 0
}

func LessThanOrEqual(total int) int {
	if total <= 5000 {
		return 500
	}
	return 0
}

func Reversed(total int) int {
	if 5000 > total {
		return 500
	}
	return 0
}

func FloatInput(total float64) int {
	if total < 5000 {
		return 500
	}
	return 0
}

func MultipleInputs(total, memberLevel int) int {
	if total < 5000 {
		return memberLevel
	}
	return 0
}

type service struct{}

func (service) ShippingFee(total int) int {
	if total < 5000 {
		return 500
	}
	return 0
}
