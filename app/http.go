package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

func (a *APITest) sendRequest(method, endpoint, payload string) error {
	endpoint = a.replaceVars(endpoint)
	payload = a.replaceVars(payload)

	if a.debug {
		fmt.Printf("Sending %s request to %s with payload: %s", method, a.baseURL+endpoint, payload)
	}

	var req *http.Request
	var err error

	if payload != "" {
		req, err = http.NewRequest(method, a.baseURL+endpoint, bytes.NewBufferString(payload))
	} else {
		req, err = http.NewRequest(method, a.baseURL+endpoint, nil)
	}

	if err != nil {
		return err
	}

	for k, v := range a.headers {
		req.Header.Set(k, a.replaceVars(v))
	}

	a.response, err = a.client.Do(req)
	if err != nil {
		return err
	}

	bodyBytes, err := io.ReadAll(a.response.Body)
	if err != nil {
		return err
	}
	a.responseBody = string(bodyBytes)
	a.response.Body.Close()

	if a.debug {
		fmt.Printf("Response status: %d", a.response.StatusCode)
		fmt.Printf("Response body: %s", a.responseBody)
	}

	return nil
}

func (a *APITest) iSendRequestTo(method, endpoint string) error {
	return a.sendRequest(method, endpoint, "")
}

func (a *APITest) iSendRequestToWithPayload(method, endpoint, payload string) error {
	return a.sendRequest(method, endpoint, payload)
}

func (a *APITest) theResponsePropertyShouldNotBeEmpty(property string) error {
	value := gjson.Get(a.responseBody, property)
	if !value.Exists() || value.String() == "" {
		return fmt.Errorf("property %s is empty or not found", property)
	}

	if a.debug {
		fmt.Printf("Property %s is not empty: %s", property, value.String())
	}

	return nil
}

func (a *APITest) theResponseShouldMatchJSON(expected string) error {
	templated := a.replaceVars(expected)

	var expectedObj any
	var actualObj any

	if err := json.Unmarshal([]byte(templated), &expectedObj); err != nil {
		validJSON, jsonErr := a.makeValidJSON(expected)
		if jsonErr != nil {
			return fmt.Errorf("invalid expected JSON: %w", err)
		}
		if err := json.Unmarshal([]byte(validJSON), &expectedObj); err != nil {
			return fmt.Errorf("still invalid JSON after fixing: %w", err)
		}
	}

	if a.debug {
		fmt.Printf("Expected JSON: %v", expectedObj)
		fmt.Printf("Actual JSON: %v", a.responseBody)
	}

	if err := json.Unmarshal([]byte(a.responseBody), &actualObj); err != nil {
		return fmt.Errorf("invalid response JSON: %w", err)
	}

	if !reflect.DeepEqual(expectedObj, actualObj) {
		return fmt.Errorf("JSON mismatch\nExpected: %v\nActual: %v",
			expectedObj, actualObj)
	}

	if a.debug {
		fmt.Printf("JSON match successful")
	}

	return nil
}

func (a *APITest) theResponseShouldContainJSON(expected string) error {
	templated := a.replaceVars(expected)

	var expectedMap map[string]any
	var actualMap map[string]any

	if err := json.Unmarshal([]byte(templated), &expectedMap); err != nil {
		validJSON, jsonErr := a.makeValidJSON(expected)
		if jsonErr != nil {
			return fmt.Errorf("invalid expected JSON: %w", err)
		}
		if err := json.Unmarshal([]byte(validJSON), &expectedMap); err != nil {
			return fmt.Errorf("still invalid JSON after fixing: %w", err)
		}
	}

	if a.debug {
		fmt.Printf("Expected JSON: %v", expectedMap)
		fmt.Printf("Actual JSON: %v", a.responseBody)
	}

	if err := json.Unmarshal([]byte(a.responseBody), &actualMap); err != nil {
		return fmt.Errorf("invalid response JSON: %w", err)
	}

	if err := containsSubset(actualMap, expectedMap); err != nil {
		return fmt.Errorf("JSON subset mismatch: %w", err)
	}

	if a.debug {
		fmt.Printf("JSON subset match successful")
	}

	return nil
}

func (a *APITest) theResponseStatusShouldBe(status int) error {
	if a.response.StatusCode != status {
		return fmt.Errorf("expected status %d but got %d with body %s", status, a.response.StatusCode, a.responseBody)
	}

	if a.debug {
		fmt.Printf("Response status %d matches expected %d", a.response.StatusCode, status)
	}

	return nil
}

func (a *APITest) theResponsePropertyShouldBe(property, expectedValue string) error {
	value := gjson.Get(a.responseBody, property)
	expected := a.replaceVars(expectedValue)

	if property == "empty" || property == "not.exists" {
		fmt.Printf("Type is %s", value.Type)
	}

	if expectedValue == "empty" {
		switch value.Type {
		case gjson.Null:
			return nil
		case gjson.String:
			if value.String() != "" {
				return fmt.Errorf("property %s should be empty but got %s", property, value.String())
			}
		default:
			return fmt.Errorf("property %s should be empty but got %s", property, value.String())
		}
		return nil
	}

	if !value.Exists() {
		return fmt.Errorf("property %s not found in response %s", property, a.responseBody)
	}

	if a.debug {
		fmt.Printf("Property %s: expected %s, got %s", property, expected, value.String())
	}

	switch {
	case expected == "true" || expected == "false":
		if value.Bool() != (expected == "true") {
			return fmt.Errorf("expected %s to be %s but got %v", property, expected, value.Bool())
		}
	case strings.HasPrefix(expected, "\"") && strings.HasSuffix(expected, "\""):
		expString := expected[1 : len(expected)-1]
		if value.String() != expString {
			return fmt.Errorf("expected %s to be %s but got %s", property, expString, value.String())
		}
	default:
		if expNum, err := strconv.ParseFloat(expected, 64); err == nil {
			if value.Float() != expNum {
				return fmt.Errorf("expected %s to be %f but got %f", property, expNum, value.Float())
			}
		} else if value.String() != expected {
			return fmt.Errorf("expected %s to be %s but got %s", property, expected, value.String())
		}
	}

	if a.debug {
		fmt.Printf("Property %s matches expected value %s", property, expected)
	}

	return nil
}
