package qumulo

func interfaceSliceToStringSlice(interfaceSlice []interface{}) []string {
	stringSlice := make([]string, len(interfaceSlice))
	for i, element := range interfaceSlice {
		stringSlice[i] = element.(string)
	}

	return stringSlice
}
