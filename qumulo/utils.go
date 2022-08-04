package qumulo

func InterfaceSliceToStringSlice(interfaceSlice []interface{}) []string {
	stringSlice := make([]string, len(interfaceSlice))
	for i, element := range interfaceSlice {
		stringSlice[i] = element.(string)
	}

	return stringSlice
}
