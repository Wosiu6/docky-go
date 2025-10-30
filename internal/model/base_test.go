package model

import "testing"

func TestParseEnvBasic(t *testing.T) {
	env := []string{"A=1", "B=two", "C=three=four", "INVALID", "A=override"}
	m := ParseEnv(env)
	if m["A"] != "override" {
		t.Errorf("expected override for A, got %q", m["A"])
	}
	if m["B"] != "two" {
		t.Errorf("expected B=two, got %q", m["B"])
	}
	if m["C"] != "three=four" {
		t.Errorf("expected full value after first '=', got %q", m["C"])
	}
	if _, ok := m["INVALID"]; ok {
		t.Errorf("INVALID should be skipped")
	}
}

func TestParseEnvEmpty(t *testing.T) {
	m := ParseEnv([]string{})
	if len(m) != 0 {
		t.Errorf("expected empty map, got %v", m)
	}
}
