package value

// Atoi converts a string into an int according to the Knight spec.
//
// This function is used instead of something like fmt.Sscanf() because it
// requires no heap allocation and is small enough to be inlined.
func Atoi(s string) int {
	num, mul := 0, 0

	for i := 0; i < len(s); i++ {
		b := s[i]

		// If mul is 0 then it means we haven't started parsing the actual
		// number part of the string yet
		//
		// If mul isn't 0 then it should fall through and break out of the loop
		// when it's caught by another condition
		if mul == 0 {
			// Skip over any leading whitespace
			if b == ' ' {
				continue
			}

			// If we're seeing a minus sign as the first character then we
			// set mul to -1, which will signify the start of parsing an
			// actual number from the string
			if b == '-' {
				mul = -1

				continue
			}
		}

		// Stop parsing when we find something that isn't a digit
		if b < '0' || b > '9' {
			break
		}

		// If we've started parsing an actual number then we should make
		// sure the mul value is set
		//
		// If it's not been set yet then it means it's a positive number, so
		// it should be set to 1
		if mul == 0 {
			mul = 1
		}

		num *= 10
		num += int(b - '0')
	}

	return num * mul
}
