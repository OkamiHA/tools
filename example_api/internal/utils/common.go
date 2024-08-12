package utils

func IsNumeric(v interface{}) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	default:
		return false
	}
}

func IsSliceContainsElement(values []string, element string) bool {
	for _, value := range values {
		if value == element {
			return true
		}
	}
	return false
}
