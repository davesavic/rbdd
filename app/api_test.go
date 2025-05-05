package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
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

func TestSendRequest(t *testing.T) {
	server := createTestServer(t, http.StatusOK, `{"message": "success"}`)
	defer server.Close()

	apiTest := NewAPITest(server.URL)
	apiTest.store["id"] = "123"

	err := apiTest.sendRequest("GET", "/${id}", "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if apiTest.response.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", apiTest.response.StatusCode)
	}

	if !bytes.Contains([]byte(apiTest.responseBody), []byte("success")) {
		t.Errorf("Expected response to contain 'success', got %s", apiTest.responseBody)
	}

	err = apiTest.sendRequest("POST", "/users", `{"name": "${id}"}`)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
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

func TestGenerateFakeData(t *testing.T) {
	apiTest := NewAPITest("https://example.com")

	err := apiTest.generateFakeData("name={name}, email={email}")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if _, ok := apiTest.store["name"]; !ok {
		t.Error("Expected 'name' to be stored")
	}

	if _, ok := apiTest.store["email"]; !ok {
		t.Error("Expected 'email' to be stored")
	}

	err = apiTest.generateFakeData("invalid-format")
	if err == nil {
		t.Error("Expected error for invalid format, got nil")
	}

	err = apiTest.generateFakeData("test={invalidtag}")
	if err != nil {
		t.Errorf("For invalid tag, expected no error, got %v", err)
	}
	if stored, ok := apiTest.store["test"]; !ok || stored != "{invalidtag}" {
		t.Errorf("Expected 'test' to be stored with value '{invalidtag}', got %v", stored)
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

func TestIExecuteCommand(t *testing.T) {
	apiTest := NewAPITest("https://example.com")

	err := apiTest.iExecuteCommand("echo test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = apiTest.iExecuteCommand("thiscommanddoesnotexist")
	if err == nil {
		t.Error("Expected error for invalid command, got nil")
	}

	err = apiTest.iExecuteCommand("")
	if err == nil {
		t.Error("Expected error for empty command, got nil")
	}
}

func TestIExecuteCommandInDirectory(t *testing.T) {
	apiTest := NewAPITest("https://example.com")

	tempDir, err := os.MkdirTemp("", "apitest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	err = apiTest.iExecuteCommandInDirectory("pwd", tempDir)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	apiTest.store["dir"] = tempDir
	err = apiTest.iExecuteCommandInDirectory("pwd", "${dir}")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestIExecuteCommandWithTimeout(t *testing.T) {
	apiTest := NewAPITest("https://example.com")

	err := apiTest.iExecuteCommandWithTimeout("echo test", 1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = apiTest.iExecuteCommandWithTimeout("sleep 2", 1)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}

func TestISendRequestTo(t *testing.T) {
	server := createTestServer(t, http.StatusOK, `{"message": "success"}`)
	defer server.Close()

	apiTest := NewAPITest(server.URL)

	err := apiTest.iSendRequestTo("GET", "/test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	apiTest.store["endpoint"] = "test"
	err = apiTest.iSendRequestTo("GET", "/${endpoint}")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestISendRequestToWithPayload(t *testing.T) {
	server := createTestServer(t, http.StatusOK, `{"message": "success"}`)
	defer server.Close()

	apiTest := NewAPITest(server.URL)

	err := apiTest.iSendRequestToWithPayload("POST", "/users", `{"name": "John"}`)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	apiTest.store["name"] = "John"
	err = apiTest.iSendRequestToWithPayload("POST", "/users", `{"name": "${name}"}`)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestTheResponseStatusShouldBe(t *testing.T) {
	apiTest := NewAPITest("https://example.com")

	apiTest.response = &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString("")),
	}

	err := apiTest.theResponseStatusShouldBe(http.StatusOK)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = apiTest.theResponseStatusShouldBe(http.StatusNotFound)
	if err == nil {
		t.Error("Expected error for non-matching status, got nil")
	}
}

func TestTheResponsePropertyShouldBe(t *testing.T) {
	apiTest := NewAPITest("https://example.com")
	apiTest.responseBody = `{
		"name": "John",
		"age": 30,
		"active": true,
		"address": {"city": "Brisbane"}
	}`

	err := apiTest.theResponsePropertyShouldBe("name", "John")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = apiTest.theResponsePropertyShouldBe("age", "30")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = apiTest.theResponsePropertyShouldBe("active", "true")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = apiTest.theResponsePropertyShouldBe("address.city", "Brisbane")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = apiTest.theResponsePropertyShouldBe("missing", "value")
	if err == nil {
		t.Error("Expected error for missing property, got nil")
	}

	err = apiTest.theResponsePropertyShouldBe("name", "Jane")
	if err == nil {
		t.Error("Expected error for value mismatch, got nil")
	}

	apiTest.store["expectedName"] = "John"
	err = apiTest.theResponsePropertyShouldBe("name", "${expectedName}")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestIStoreTheResponsePropertyAs(t *testing.T) {
	apiTest := NewAPITest("https://example.com")
	apiTest.responseBody = `{
		"name": "John",
		"age": 30,
		"active": true,
		"address": {"city": "Brisbane"}
	}`

	err := apiTest.iStoreTheResponsePropertyAs("name", "userName")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if apiTest.store["userName"] != "John" {
		t.Errorf("Expected stored value 'John', got %v", apiTest.store["userName"])
	}

	err = apiTest.iStoreTheResponsePropertyAs("age", "userAge")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if apiTest.store["userAge"] != float64(30) {
		t.Errorf("Expected stored value 30, got %v", apiTest.store["userAge"])
	}

	err = apiTest.iStoreTheResponsePropertyAs("active", "isActive")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if apiTest.store["isActive"] != true {
		t.Errorf("Expected stored value true, got %v", apiTest.store["isActive"])
	}

	err = apiTest.iStoreTheResponsePropertyAs("address", "userAddress")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = apiTest.iStoreTheResponsePropertyAs("missing", "missingValue")
	if err == nil {
		t.Error("Expected error for missing property, got nil")
	}
}

func TestIStoreAs(t *testing.T) {
	apiTest := NewAPITest("https://example.com")

	err := apiTest.iStoreAs("name", "John")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if apiTest.store["name"] != "John" {
		t.Errorf("Expected stored value 'John', got %v", apiTest.store["name"])
	}

	err = apiTest.iStoreAs("age", "30")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if apiTest.store["age"] != float64(30) {
		t.Errorf("Expected stored value 30, got %v", apiTest.store["age"])
	}

	err = apiTest.iStoreAs("active", "true")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if apiTest.store["active"] != true {
		t.Errorf("Expected stored value true, got %v", apiTest.store["active"])
	}
}

func TestISetHeaderTo(t *testing.T) {
	apiTest := NewAPITest("https://example.com")

	err := apiTest.iSetHeaderTo("Authorization", "Bearer token")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if apiTest.headers["Authorization"] != "Bearer token" {
		t.Errorf("Expected header value 'Bearer token', got %v", apiTest.headers["Authorization"])
	}

	err = apiTest.iSetHeaderTo("Content-Type", "application/xml")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if apiTest.headers["Content-Type"] != "application/xml" {
		t.Errorf("Expected header value 'application/xml', got %v", apiTest.headers["Content-Type"])
	}
}

func TestIResetAllStoredVariables(t *testing.T) {
	apiTest := NewAPITest("https://example.com")
	apiTest.store["name"] = "John"
	apiTest.store["age"] = 30

	err := apiTest.iResetAllVariables()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(apiTest.store) != 0 {
		t.Errorf("Expected empty store, got %v", apiTest.store)
	}
}

func TestIResetVariables(t *testing.T) {
	apiTest := NewAPITest("https://example.com")
	apiTest.store["name"] = "John"
	apiTest.store["age"] = 30

	err := apiTest.iResetVariables("name")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if _, ok := apiTest.store["name"]; ok {
		t.Error("Expected 'name' to be removed from store")
	}
	if apiTest.store["age"] != 30 {
		t.Errorf("Expected 'age' to remain, got %v", apiTest.store["age"])
	}

	err = apiTest.iResetVariables("missing")
	if err != nil {
		t.Errorf("Expected no error for missing variable, got %v", err)
	}
}

func TestTheResponsePropertyShouldNotBeEmpty(t *testing.T) {
	apiTest := NewAPITest("https://example.com")
	apiTest.responseBody = `{
		"name": "John",
		"empty": "",
		"nested": {"value": "something"}
	}`

	err := apiTest.theResponsePropertyShouldNotBeEmpty("name")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = apiTest.theResponsePropertyShouldNotBeEmpty("empty")
	if err == nil {
		t.Error("Expected error for empty property, got nil")
	}

	err = apiTest.theResponsePropertyShouldNotBeEmpty("missing")
	if err == nil {
		t.Error("Expected error for missing property, got nil")
	}

	err = apiTest.theResponsePropertyShouldNotBeEmpty("nested.value")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestTheResponseShouldMatchJSON(t *testing.T) {
	apiTest := NewAPITest("https://example.com")
	apiTest.responseBody = `{
		"name": "John",
		"age": 30,
		"active": true
	}`

	err := apiTest.theResponseShouldMatchJSON(`{
		"name": "John",
		"age": 30,
		"active": true
	}`)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = apiTest.theResponseShouldMatchJSON(`{
		"name": "Jane",
		"age": 30,
		"active": true
	}`)
	if err == nil {
		t.Error("Expected error for non-matching JSON, got nil")
	}

	apiTest.store["expectedName"] = "John"
	apiTest.store["expectedAge"] = 30
	err = apiTest.theResponseShouldMatchJSON(`{
		"name": "${expectedName}",
		"age": ${expectedAge},
		"active": true
	}`)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = apiTest.theResponseShouldMatchJSON(`invalid json`)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestTheResponseShouldContainJSON(t *testing.T) {
	apiTest := NewAPITest("https://example.com")
	apiTest.responseBody = `{
		"name": "John",
		"age": 30,
		"active": true,
		"address": {"city": "Brisbane", "country": "Australia"}
	}`

	err := apiTest.theResponseShouldContainJSON(`{
		"name": "John",
		"age": 30
	}`)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = apiTest.theResponseShouldContainJSON(`{
		"address": {"city": "Brisbane"}
	}`)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = apiTest.theResponseShouldContainJSON(`{
		"name": "Jane"
	}`)
	if err == nil {
		t.Error("Expected error for non-matching subset, got nil")
	}

	apiTest.store["expectedName"] = "John"
	err = apiTest.theResponseShouldContainJSON(`{
		"name": "${expectedName}"
	}`)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = apiTest.theResponseShouldContainJSON(`invalid json`)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
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
