package qumulo

import (
	"fmt"
	"strconv"
	"strings"
)

func InterfaceSliceToStringSlice(interfaceSlice []interface{}) []string {
	stringSlice := make([]string, len(interfaceSlice))
	for i, element := range interfaceSlice {
		stringSlice[i] = element.(string)
	}

	return stringSlice
}

func ParseLocalGroupMemberId(id string) ([]string, error) {
	ids := strings.Split(id, ":")
	if len(ids) != 2 {
		return nil, fmt.Errorf("Local group member ID is malformed. Expected ID of the form {group_id}:{user_id}, got %q", id)
	}
	if _, err := strconv.Atoi(ids[0]); err == nil {
		return nil, fmt.Errorf("Group ID is not a number. Got %q", ids[0])
	}
	if _, err := strconv.Atoi(ids[1]); err == nil {
		return nil, fmt.Errorf("User ID is not a number. Got %q", ids[1])
	}

	return ids, nil
}

func FormLocalGroupMemberId(ids []string) (string, error) {
	if len(ids) != 2 {
		return "", fmt.Errorf("Expected exactly 2 IDs, got: %q", ids)
	}
	return fmt.Sprintf("%s:%s", ids[0], ids[1]), nil

}
