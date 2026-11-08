package utils

func SetIfPresent[T any](target *T, source *T) {
	if source != nil {
		*target = *source
	}
}
