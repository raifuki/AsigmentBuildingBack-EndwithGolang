package config

import "testing"

func TestLoadUsesRenderEnvironmentValues(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://render-db")
	t.Setenv("REDIS_URL", "redis://render-redis:6379")
	t.Setenv("JWT_SECRET", "render-secret")
	t.Setenv("APP_PORT", "9090")
	t.Setenv("PORT", "8080")

	cfg := Load()

	if cfg.DatabaseURL != "postgres://render-db" {
		t.Fatalf("expected DATABASE_URL override, got %q", cfg.DatabaseURL)
	}
	if cfg.RedisURL != "redis://render-redis:6379" {
		t.Fatalf("expected REDIS_URL override, got %q", cfg.RedisURL)
	}
	if cfg.JWTSecret != "render-secret" {
		t.Fatalf("expected JWT_SECRET override, got %q", cfg.JWTSecret)
	}
	if cfg.AppPort != "9090" {
		t.Fatalf("expected APP_PORT override, got %q", cfg.AppPort)
	}
}
