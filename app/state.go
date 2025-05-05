package app

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

func (a *APITest) iStoreTheResponsePropertyAs(property, variable string) error {
	value := gjson.Get(a.responseBody, property)
	if !value.Exists() {
		return fmt.Errorf("property %s not found in response", property)
	}

	switch value.Type {
	case gjson.String:
		a.store[variable] = value.String()
	case gjson.Number:
		a.store[variable] = value.Float()
	case gjson.True, gjson.False:
		a.store[variable] = value.Bool()
	case gjson.JSON:
		a.store[variable] = value.Raw
	default:
		a.store[variable] = value.Raw
	}

	return nil
}

func (a *APITest) iStoreTheCommandOutputAs(variable string) error {
	if a.commandOutput == "" {
		return fmt.Errorf("command output is empty")
	}
	a.store[variable] = a.commandOutput
	return nil
}

func (a *APITest) iStoreAs(variable, value string) error {
	value = a.replaceVars(value)
	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		value = value[1 : len(value)-1]
	}
	if strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}") {
		var jsonObj map[string]any
		if err := json.Unmarshal([]byte(value), &jsonObj); err != nil {
			return fmt.Errorf("invalid JSON format: %w", err)
		}
		a.store[variable] = jsonObj
	} else if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		var jsonArray []any
		if err := json.Unmarshal([]byte(value), &jsonArray); err != nil {
			return fmt.Errorf("invalid JSON array format: %w", err)
		}
		a.store[variable] = jsonArray
	} else if num, err := strconv.ParseFloat(value, 64); err == nil {
		a.store[variable] = num
	} else if boolVal, err := strconv.ParseBool(value); err == nil {
		a.store[variable] = boolVal
	} else {
		a.store[variable] = value
	}
	return nil
}

func (a *APITest) iSetHeaderTo(header, value string) error {
	a.headers[header] = value
	return nil
}

func (a *APITest) iResetAllVariables() error {
	a.store = map[string]any{}
	return nil
}

func (a *APITest) iResetVariables(variables string) error {
	vars := strings.SplitSeq(variables, ",")
	for variable := range vars {
		varName := strings.TrimSpace(variable)
		delete(a.store, varName)
	}
	return nil
}
