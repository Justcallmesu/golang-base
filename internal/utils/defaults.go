package utils

func GetDefaultValue[T comparable](value T, defaultValue T) T {

	var zero T

	if value == zero {
		return defaultValue
	}

	return value

}
