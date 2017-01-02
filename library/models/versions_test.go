package models

import "testing"

var errorMsg string = "Expected %s, Got %s"

func TestStringifyVersion(t *testing.T) {
	v := Version{
		Name:  "test",
		Major: 0,
		Minor: 0,
		Patch: 1,
	}
	expected := "test (0.0.1)"
	if v.String() != expected {
		t.Errorf(errorMsg, expected, v.String())
	}
}

func TestParse(t *testing.T) {
	v := Version{
		Name:  "test",
		Major: 0,
		Minor: 0,
		Patch: 1,
	}
	newV := Version{}
	newV.Parse("test (0.0.1)")
	if v.String() != newV.String() {
		t.Errorf(errorMsg, v.String(), newV.String())
	}
}
func TestNonMatchingVersion(t *testing.T) {
	v := Version{
		Name:  "test",
		Major: 0,
		Minor: 0,
		Patch: 1,
	}
	newV := Version{
		Name:  "test",
		Major: 0,
		Minor: 2,
		Patch: 1,
	}
	if newV.Match(v) == true {
		t.Errorf(errorMsg, "false", "true")
	}
}
func TestMatchingVersion(t *testing.T) {
	v := Version{
		Name:  "test",
		Major: 0,
		Minor: 0,
		Patch: 1,
	}
	newV := Version{
		Name:  "test",
		Major: 0,
		Minor: 0,
		Patch: 1,
	}
	if newV.Match(v) == false {
		t.Errorf(errorMsg, "true", "false")
	}
}
