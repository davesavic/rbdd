package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestNewAPITest(t *testing.T) {
	baseURL := "https://example.com"
	apiTest := NewAPITest(baseURL)

	if apiTest.baseURL != baseURL {
		t.Errorf("Expected baseURL to be %s, got %s", baseURL, apiTest.baseURL)
	}

	if apiTest.client == nil {
		t.Error("Expected client to be initialized, got nil")
	}

	if len(apiTest.headers) != 1 || apiTest.headers["Content-Type"] != "application/json" {
		t.Errorf("Expected default Content-Type header, got %v", apiTest.headers)
	}

	if apiTest.store == nil {
		t.Error("Expected store to be initialized, got nil")
	}
}

func TestReplaceVars(t *testing.T) {
	apiTest := NewAPITest("https://example.com")
	apiTest.store["name"] = "John"
	apiTest.store["age"] = 30

	tests := []struct {
		input    string
		expected string
	}{
		{"Hello ${name}", "Hello John"},
		{"${name} is ${age} years old", "John is 30 years old"},
		{"No variables here", "No variables here"},
		{"Missing ${unknown} variable", "Missing ${unknown} variable"},
	}

	for _, test := range tests {
		result := apiTest.replaceVars(test.input)
		if result != test.expected {
			t.Errorf("For %q, expected %q, got %q", test.input, test.expected, result)
		}
	}
}

func createTestServer(_ *testing.T, status int, response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		fmt.Fprintln(w, response)
	}))
}

func TestGenerateFromTag(t *testing.T) {
	tests := []string{
		"{name}",
		"{email}",
		"{ipv4address}",
		"{color}",
	}

	for _, test := range tests {
		result, err := generateFromTag(test)
		if err != nil {
			t.Errorf("For %q, expected no error, got %v", test, err)
		}
		if result == "" {
			t.Errorf("For %q, got empty result", test)
		}
		if result == test {
			t.Errorf("For %q, expected generated value, got same string back", test)
		}
	}

	// Should return same string, not error
	invalidTag := "{invalidtag}"
	result, err := generateFromTag(invalidTag)
	if err != nil {
		t.Errorf("For invalid tag, expected no error, got %v", err)
	}
	if result != invalidTag {
		t.Errorf("For invalid tag, expected original string %q, got %q", invalidTag, result)
	}
}

func TestSplitByCommaOutsideBrackets(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"a,b,c", []string{"a", "b", "c"}},
		{"a{1,2},b,c", []string{"a{1,2}", "b", "c"}},
		{"a[1,2,3],b{x,y}", []string{"a[1,2,3]", "b{x,y}"}},
		{"", nil},
	}

	for _, test := range tests {
		result := splitByCommaOutsideBrackets(test.input)
		if (result == nil) != (test.expected == nil) {
			t.Errorf("For %q, nil check failed: expected nil? %v, got nil? %v",
				test.input, test.expected == nil, result == nil)
		} else if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("For %q, expected %v, got %v", test.input, test.expected, result)
		}
	}
}

func TestContainsSubset(t *testing.T) {
	data := map[string]any{
		"name": "John",
		"age":  30,
		"address": map[string]any{
			"city":    "Brisbane",
			"country": "Australia",
		},
	}

	err := containsSubset(data, data)
	if err != nil {
		t.Errorf("Expected no error for full match, got %v", err)
	}

	subset := map[string]any{
		"name": "John",
		"address": map[string]any{
			"city": "Brisbane",
		},
	}
	err = containsSubset(data, subset)
	if err != nil {
		t.Errorf("Expected no error for partial match, got %v", err)
	}

	missingKey := map[string]any{
		"missing": "value",
	}
	err = containsSubset(data, missingKey)
	if err == nil {
		t.Error("Expected error for missing key, got nil")
	}

	nestedMissingKey := map[string]any{
		"address": map[string]any{
			"missing": "value",
		},
	}
	err = containsSubset(data, nestedMissingKey)
	if err == nil {
		t.Error("Expected error for nested missing key, got nil")
	}

	valueMismatch := map[string]any{
		"name": "Jane",
	}
	err = containsSubset(data, valueMismatch)
	if err == nil {
		t.Error("Expected error for value mismatch, got nil")
	}
}

func TestMakeValidJSON(t *testing.T) {
	apiTest := NewAPITest("https://example.com")
	apiTest.store["name"] = "John"
	apiTest.store["age"] = 30
	apiTest.store["active"] = true

	template := `{"name": "${name}"}`
	jsonString := apiTest.replaceVars(template)

	var data map[string]any
	err := json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}
	if data["name"] != "John" {
		t.Errorf("Expected name to be 'John', got %v", data["name"])
	}

	template = `{"name": "${name}", "age": ${age}, "active": ${active}}`
	jsonString = apiTest.replaceVars(template)

	// For non-string values, we need to handle JSON parsing separately
	// This simulates what makeValidJSON is trying to do
	parsedTemplate := strings.Replace(jsonString, "\"${age}\"", "30", 1)
	parsedTemplate = strings.Replace(parsedTemplate, "\"${active}\"", "true", 1)

	err = json.Unmarshal([]byte(parsedTemplate), &data)
	if err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}

	if data["age"] != float64(30) {
		t.Errorf("Expected age to be 30, got %v", data["age"])
	}
	if data["active"] != true {
		t.Errorf("Expected active to be true, got %v", data["active"])
	}
}
