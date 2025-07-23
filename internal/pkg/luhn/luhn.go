package luhn

func IsValid(number string) bool {
	if number == "" {
		return false
	}

	sum := 0
	nDigits := len(number)
	parity := nDigits % 2

	for i := 0; i < nDigits; i++ {
		digit := int(number[i] - '0')
		if digit < 0 || digit > 9 {
			return false
		}
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	return sum%10 == 0
}
