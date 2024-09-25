package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_WithEnvVariables(t *testing.T) {
	_ = os.Setenv("DB_USER", "testuser")
	_ = os.Setenv("DB_PASSWORD", "testpassword")
	_ = os.Setenv("DB_HOST", "testhost")
	_ = os.Setenv("DB_PORT", "1234")
	_ = os.Setenv("DB_NAME", "testdb")
	_ = os.Setenv("CACERT_PATH", "testcert")
	_ = os.Setenv("CLIENT_CERT_PATH", "testcert")

	defer func() {
		_ = os.Unsetenv("DB_USER")
		_ = os.Unsetenv("DB_PASSWORD")
		_ = os.Unsetenv("DB_HOST")
		_ = os.Unsetenv("DB_PORT")
		_ = os.Unsetenv("DB_NAME")
		_ = os.Unsetenv("CACERT_PATH")
		_ = os.Unsetenv("CLIENT_CERT_PATH")
	}()

	cfg := LoadConfig()

	assert.Equal(t, "testuser", cfg.DbUser)
	assert.Equal(t, "testpassword", cfg.DbPassword)
	assert.Equal(t, "testhost", cfg.DbHost)
	assert.Equal(t, "1234", cfg.DbPort)
	assert.Equal(t, "testdb", cfg.DbName)

}

func TestLoadConfig_WithoutEnvVariables(t *testing.T) {
	_ = os.Unsetenv("DB_USER")
	_ = os.Unsetenv("DB_PASSWORD")
	_ = os.Unsetenv("DB_HOST")
	_ = os.Unsetenv("DB_PORT")
	_ = os.Unsetenv("DB_NAME")

	assert.Panics(t, func() { LoadConfig() }, "load config panicked")

}
