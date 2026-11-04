package utils

func GetDefaultValue[T comparable](value T, defaultValue T) T {
	if value == *new(T) {
		return defaultValue
	}
	return value
}
