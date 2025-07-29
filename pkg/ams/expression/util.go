package expression

func StringArrayFrom(input []string) StringArray {
	result := make(StringArray, len(input))
	for i, v := range input {
		result[i] = String(v)
	}
	return result
}
