package util

func CoalesceString(val, def string) string {
	if val == "" {
		return def
	}
	return val
}

func CoalesceInt(val, def int) int {
	if val == 0 {
		return def
	}
	return val
}
