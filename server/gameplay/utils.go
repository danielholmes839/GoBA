package gameplay

func direction(a int) int {
	if a == 0 {
		return 0
	} else if a < 0 {
		return -1
	} else {
		return 1
	}
}