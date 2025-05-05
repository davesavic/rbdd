package app

import "testing"

func TestIGenerateFakeData(t *testing.T) {
	apiTest := NewAPITest("https://example.com")

	err := apiTest.iGenerateFakeData("name={name}, email={email}")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if _, ok := apiTest.store["name"]; !ok {
		t.Error("Expected 'name' to be stored")
	}

	if _, ok := apiTest.store["email"]; !ok {
		t.Error("Expected 'email' to be stored")
	}

	err = apiTest.iGenerateFakeData("invalid-format")
	if err == nil {
		t.Error("Expected error for invalid format, got nil")
	}

	err = apiTest.iGenerateFakeData("test={invalidtag}")
	if err != nil {
		t.Errorf("For invalid tag, expected no error, got %v", err)
	}
	if stored, ok := apiTest.store["test"]; !ok || stored != "{invalidtag}" {
		t.Errorf("Expected 'test' to be stored with value '{invalidtag}', got %v", stored)
	}
}
