package app

import "testing"

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

	err := apiTest.iStoreAs("John", "name")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if apiTest.store["name"] != "John" {
		t.Errorf("Expected stored value 'John', got %v", apiTest.store["name"])
	}

	err = apiTest.iStoreAs("30", "age")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if apiTest.store["age"] != float64(30) {
		t.Errorf("Expected stored value 30, got %v", apiTest.store["age"])
	}

	err = apiTest.iStoreAs("true", "active")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if apiTest.store["active"] != true {
		t.Errorf("Expected stored value true, got %v", apiTest.store["active"])
	}

	err = apiTest.iStoreAs("${age}", "otherAge")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if apiTest.store["otherAge"] != float64(30) {
		t.Errorf("Expected stored value 30, got %v", apiTest.store["otherAge"])
	}
}

func TestIStoreCommandOutputAs(t *testing.T) {
	apiTest := NewAPITest("https://example.com")

	err := apiTest.iExecuteCommand("echo 'Hello World'")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = apiTest.iStoreTheCommandOutputAs("output")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if apiTest.store["output"] != "Hello World" {
		t.Errorf("Expected stored value 'Hello World', got %v", apiTest.store["output"])
	}

	err = apiTest.iExecuteCommand("thiscommanddoesnotexist")
	if err == nil {
		t.Error("Expected error for invalid command, got nil")
	}
	err = apiTest.iStoreTheCommandOutputAs("invalid")
	if err == nil {
		t.Error("Expected error for invalid command output, got nil")
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
