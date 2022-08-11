package qumulo

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type StringOrInt struct {
	isString  bool
	isInt     bool
	stringVal string
	intVal    int
}

func (si *StringOrInt) UnmarshalJSON(data []byte) error {
	var i interface{}
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	if intVal, err := i.(int); !err {
		si.isString = false
		si.isInt = true
		si.intVal = intVal
	} else if stringVal, err := i.(string); !err {
		si.isString = true
		si.isInt = false
		si.stringVal = stringVal
	} else {
		return fmt.Errorf("unknown input for StringOrInt: got value %q", i)
	}

	return nil
}

func SetStringOrIntValue(d *schema.ResourceData, field string, si StringOrInt) error {
	if si.isString {
		d.Set(field, si.stringVal)
	} else if si.isInt {
		d.Set(field, fmt.Sprintf("%q", si.intVal))
	} else {
		return fmt.Errorf("unknown value/type for StringOrInt, not setting resource data: got %+v", si)
	}

	return nil
}

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
		return nil, fmt.Errorf("local group member ID is malformed. Expected ID of the form {group_id}:{user_id}, got %q", id)
	}
	if _, err := strconv.Atoi(ids[0]); err != nil {
		return nil, fmt.Errorf("group ID is not a number. Got %q", ids[0])
	}
	if _, err := strconv.Atoi(ids[1]); err != nil {
		return nil, fmt.Errorf("user ID is not a number. Got %q", ids[1])
	}

	return ids, nil
}

func FormLocalGroupMemberId(ids []string) (string, error) {
	if len(ids) != 2 {
		return "", fmt.Errorf("expected exactly 2 IDs, got: %q", ids)
	}
	if _, err := strconv.Atoi(ids[0]); err != nil {
		return "", fmt.Errorf("group ID is not a number. Got %q", ids[0])
	}
	if _, err := strconv.Atoi(ids[1]); err != nil {
		return "", fmt.Errorf("user ID is not a number. Got %q", ids[1])
	}
	return fmt.Sprintf("%s:%s", ids[0], ids[1]), nil

}

func ParseRoleMemberId(id string) ([]string, error) {
	ids := strings.Split(id, ":")
	if len(ids) != 2 {
		return nil, fmt.Errorf("role member ID is malformed. Expected ID of the form {role_name}:{user_id}, got %q", id)
	}
	if _, err := strconv.Atoi(ids[1]); err != nil {
		return nil, fmt.Errorf("user ID is not a number. Got %q", ids[1])
	}

	return ids, nil
}

func FormRoleMemberId(ids []string) (string, error) {
	if len(ids) != 2 {
		return "", fmt.Errorf("expected exactly 2 IDs, got: %q", ids)
	}
	if _, err := strconv.Atoi(ids[1]); err != nil {
		return "", fmt.Errorf("user ID is not a number. Got %q", ids[1])
	}
	return fmt.Sprintf("%s:%s", ids[0], ids[1]), nil

}
