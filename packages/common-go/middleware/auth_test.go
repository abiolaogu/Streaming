package middleware

import (
	"reflect"
	"testing"
)

func TestExtractRolesFromStringSlice(t *testing.T) {
	got := extractRoles([]string{"admin", "user"})
	want := []string{"admin", "user"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected roles: got %v want %v", got, want)
	}
}

func TestExtractRolesFromInterfaceSlice(t *testing.T) {
	got := extractRoles([]interface{}{"admin", 12, " user "})
	want := []string{"admin", "user"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected roles: got %v want %v", got, want)
	}
}

func TestExtractRolesFromString(t *testing.T) {
	got := extractRoles("admin")
	want := []string{"admin"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected roles: got %v want %v", got, want)
	}
}

func TestExtractRolesFromUnsupportedType(t *testing.T) {
	if got := extractRoles(42); got != nil {
		t.Fatalf("expected nil roles for unsupported type, got %v", got)
	}
}
