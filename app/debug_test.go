package app

import "testing"

func TestIStartDebugging(t *testing.T) {
	test := NewAPITest("http://localhost:8080")
	test.iStartDebugging()
	if !test.debug {
		t.Errorf("Expected debug to be true, got false")
	}
}

func TestIStopDebugging(t *testing.T) {
	test := NewAPITest("http://localhost:8080")
	test.iStartDebugging()
	test.iStopDebugging()
	if test.debug {
		t.Errorf("Expected debug to be false, got true")
	}
}
