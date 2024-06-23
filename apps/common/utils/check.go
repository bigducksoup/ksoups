package utils

func UnexpectThenDefault[T comparable](originalVal T, unexpectValue T, defaultValue T) T {
	if originalVal == unexpectValue {
		return defaultValue
	}

	return originalVal
}
