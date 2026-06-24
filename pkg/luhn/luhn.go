package luhn

 
func Valid(s string) bool {
	 
	if len(s) < 2 {
		return false
	}

	sum := 0
	n := len(s)

	for i := range n {
		 
		if s[i] < '0' || s[i] > '9' {
			return false
		}

		digit := int(s[i] - '0')

	 
		if (n-1-i)%2 == 1 {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
	}

	return sum%10 == 0
}
