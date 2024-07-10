package common

func Capitalize(s string) string {
	if len(s) <= 0 {
		return ""
	}
	tbr := []rune(s)
	for i := 0; i < len(tbr); i++ {
		if tbr[i] >= 'A' && tbr[i] <= 'Z' || tbr[i] >= 'a' && tbr[i] <= 'z' || tbr[i] >= '0' && tbr[i] <= '9' {
			if tbr[i] >= 'a' && tbr[i] <= 'z' {
				tbr[i] -= 32
			}
			i++
			for i < len(tbr) && (tbr[i] >= 'a' && tbr[i] <= 'z' || tbr[i] >= 'A' && tbr[i] <= 'Z' || tbr[i] >= '0' && tbr[i] <= '9') {
				if tbr[i] >= 'A' && tbr[i] <= 'Z' {
					tbr[i] += 32
				}
				i++
			}
		}
	}
	return string(tbr)
}
