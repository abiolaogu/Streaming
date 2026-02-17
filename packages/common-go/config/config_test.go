package config

import "testing"

func assertPanics(t *testing.T, fn func(), msg string) {
	t.Helper()

	defer func() {
		if recover() == nil {
			t.Fatalf("expected panic: %s", msg)
		}
	}()

	fn()
}

func assertDoesNotPanic(t *testing.T, fn func(), msg string) {
	t.Helper()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("did not expect panic (%s), got: %v", msg, r)
		}
	}()

	fn()
}

func TestLoadPanicsInProductionWhenSecretMissing(t *testing.T) {
	t.Setenv("ENVIRONMENT", "production")
	t.Setenv("APP_ENV", "")
	t.Setenv("JWT_SECRET_KEY", "")

	assertPanics(t, func() {
		Load()
	}, "production missing JWT secret should fail fast")
}

func TestLoadPanicsInProductionWhenSecretTooShort(t *testing.T) {
	t.Setenv("ENVIRONMENT", "production")
	t.Setenv("APP_ENV", "")
	t.Setenv("JWT_SECRET_KEY", "short-secret")

	assertPanics(t, func() {
		Load()
	}, "production short JWT secret should fail fast")
}

func TestLoadAllowsDevelopmentDefaults(t *testing.T) {
	t.Setenv("ENVIRONMENT", "development")
	t.Setenv("APP_ENV", "")
	t.Setenv("JWT_SECRET_KEY", "")

	assertDoesNotPanic(t, func() {
		Load()
	}, "development should allow fallback secret")
}

func TestLoadAllowsStrongProductionSecret(t *testing.T) {
	t.Setenv("ENVIRONMENT", "production")
	t.Setenv("APP_ENV", "")
	t.Setenv("JWT_SECRET_KEY", "12345678901234567890123456789012")

	assertDoesNotPanic(t, func() {
		Load()
	}, "production should allow 32+ character secret")
}
