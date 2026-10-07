package slug

func IsValid(value string, maxLength int) bool {
	if len(value) == 0 || len(value) > maxLength || !isLowerAlphaNumeric(value[0]) {
		return false
	}
	for i := 1; i < len(value); i++ {
		if !isLowerAlphaNumeric(value[i]) && value[i] != '-' {
			return false
		}
		if value[i] == '-' && value[i-1] == '-' {
			return false
		}
	}
	return value[len(value)-1] != '-'
}

func isLowerAlphaNumeric(char byte) bool {
	return char >= 'a' && char <= 'z' || char >= '0' && char <= '9'
}
