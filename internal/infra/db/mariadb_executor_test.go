package db

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"spend-api/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockReadFile to replace os.ReadFile during tests
var mockReadFile = func(filename string) ([]byte, error) {
	return []byte("mocked certificate content"), nil
}

// mockLoadX509KeyPair to replace tls.LoadX509KeyPair during tests
var mockLoadX509KeyPair = func(certFile, keyFile string) (tls.Certificate, error) {
	return tls.Certificate{}, nil
}

// mockNewCertPool to replace x509.NewCertPool during tests
var mockNewCertPool = func() *x509.CertPool {
	return &x509.CertPool{}
}

// mockAppendCertsFromPEM to replace x509.CertPool.AppendCertsFromPEM during tests
var mockAppendCertsFromPEM = func(pool *x509.CertPool, pemCerts []byte) bool {
	return true
}

// mockLoadX509KeyPairError to simulate tls.LoadX509KeyPair failure
var mockLoadX509KeyPairError = func(certFile, keyFile string) (tls.Certificate, error) {
	return tls.Certificate{}, errors.New("failed to load client cert and key")
}

// mockAppendCertsFromPEMError to simulate x509.CertPool.AppendCertsFromPEM failure
var mockAppendCertsFromPEMError = func(pool *x509.CertPool, pemCerts []byte) bool {
	return false
}

// mockRegisterTLSConfigError to simulate TLS Registration Error
var mockRegisterTLSConfigError = func(name string, config *tls.Config) error {
	return errors.New("failed to register TLS config")
}

func TestCreateDSN(t *testing.T) {

	cfg := &config.Config{
		DbUser:     "testuser",
		DbPassword: "testpassword",
		DbHost:     "localhost",
		DbPort:     "3306",
		DbName:     "testdb",
	}

	expectedDSN := "testuser:testpassword@tcp(localhost:3306)/testdb"

	dsn := createDSN(cfg)

	assert.Equal(t, expectedDSN, dsn)
}

func TestSetupTLSConfig(t *testing.T) {
	// Define table of test cases
	tests := []struct {
		name             string
		mockReadFile     func(filename string) ([]byte, error)
		mockAppendCerts  func(pool *x509.CertPool, pemCerts []byte) bool
		mockLoadX509Key  func(certFile, keyFile string) (tls.Certificate, error)
		mockTLSConfigErr func(name string, config *tls.Config) error
		expectPanic      bool
		expectedPanicMsg string
	}{
		{
			name:             "ReadFile fails",
			mockReadFile:     readFile,
			mockAppendCerts:  appendCertsFromPEM,
			mockLoadX509Key:  loadX509KeyPair,
			expectPanic:      true,
			expectedPanicMsg: "failed to read cert file",
		},
		{
			name:             "AppendCertsFromPEM fails",
			mockReadFile:     mockReadFile,
			mockAppendCerts:  mockAppendCertsFromPEMError,
			mockLoadX509Key:  loadX509KeyPair,
			expectPanic:      true,
			expectedPanicMsg: "failed to append CA cert",
		},
		{
			name:             "LoadX509KeyPair fails",
			mockReadFile:     mockReadFile,
			mockAppendCerts:  mockAppendCertsFromPEM,
			mockLoadX509Key:  mockLoadX509KeyPairError,
			expectPanic:      true,
			expectedPanicMsg: "failed to load client cert and key",
		},
		{
			name:             "Registration fails",
			mockReadFile:     mockReadFile,
			mockAppendCerts:  mockAppendCertsFromPEM,
			mockLoadX509Key:  mockLoadX509KeyPair,
			mockTLSConfigErr: mockRegisterTLSConfigError,
			expectPanic:      true,
			expectedPanicMsg: "failed to register TLS config",
		},
		{
			name:             "All pass",
			mockReadFile:     mockReadFile,
			mockAppendCerts:  mockAppendCertsFromPEM,
			mockLoadX509Key:  mockLoadX509KeyPair,
			mockTLSConfigErr: registerTLSConfig,
			expectPanic:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalReadFile := readFile
			readFile = tt.mockReadFile
			defer func() { readFile = originalReadFile }()

			originalAppendCertsFromPEM := appendCertsFromPEM
			appendCertsFromPEM = tt.mockAppendCerts
			defer func() { appendCertsFromPEM = originalAppendCertsFromPEM }()

			originalLoadX509KeyPair := loadX509KeyPair
			loadX509KeyPair = tt.mockLoadX509Key
			defer func() { loadX509KeyPair = originalLoadX509KeyPair }()

			originalRegisterTLSConfig := registerTLSConfig
			registerTLSConfig = tt.mockTLSConfigErr
			defer func() { registerTLSConfig = originalRegisterTLSConfig }()

			cfg := &config.Config{
				CACertPath:     "/mock/path/ca-cert.pem",
				ClientCertPath: "/mock/path/client-cert.pem",
				ClientKeyPath:  "/mock/path/client-key.pem",
			}

			dsn := "testuser:testpassword@tcp(localhost:3306)/testdb"

			if tt.expectPanic {
				assert.Panics(t, func() {
					setupTLSConfig(dsn, cfg)
				})
			} else {
				assert.NotPanics(t, func() {
					setupTLSConfig(dsn, cfg)
				})
			}
		})
	}
}

func TestMariaDbExecutor_Exec_Success(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	mock.ExpectExec("INSERT INTO accounts").
		WithArgs("John Doe").
		WillReturnResult(sqlmock.NewResult(1, 1))

	executor := &MariaDbExecutor{mockDB}

	result, err := executor.Exec("INSERT INTO accounts (name) VALUES (?)", "John Doe")
	assert.NoError(t, err)

	rowsAffected, err := result.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMariaDbExecutor_Exec_Failure(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	mock.ExpectExec("INSERT INTO accounts").
		WithArgs("Account1").
		WillReturnError(sql.ErrConnDone)

	executor := &MariaDbExecutor{mockDB}

	_, err = executor.Exec("INSERT INTO accounts (name) VALUES (?)", "Account1")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, sql.ErrConnDone), "Expected error to be sql.ErrConnDone")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMariaDbExecutor_Close_Success(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	mock.ExpectClose()

	executor := &MariaDbExecutor{mockDB}

	err = executor.Close()
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
