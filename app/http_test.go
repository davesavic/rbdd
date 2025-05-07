package app

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

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
		"address": {"city": "Brisbane"},
		"empty": "",
		"null": null
	}`

	err := apiTest.theResponsePropertyShouldBe("not.exists", "empty")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = apiTest.theResponsePropertyShouldBe("empty", "")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = apiTest.theResponsePropertyShouldBe("null", "empty")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = apiTest.theResponsePropertyShouldBe("name", "John")
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
