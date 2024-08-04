var stov = map[rune]int{
	'I': 1,
	'V': 5,
	'X': 10,
	'L': 50,
	'C': 100,
	'D': 500,
	'M': 1000,
}

var validSigil = map[rune]rune{
	'V': 'I',
	'X': 'I',
	'L': 'X',
	'C': 'X',
	'D': 'C',
	'M': 'C',
}

func romanToInt(s string) int {
	var val int
	var prev rune
	for _, curr := range s {
		if prev != 0 && prev == validSigil[curr] {
			val -= stov[prev]
		} else {
			val += stov[curr]
		}
	}

	return val
}
