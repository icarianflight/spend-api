package db

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"net/url"
	"os"
	"spend-api/internal/config"
)

var readFile = os.ReadFile
var loadX509KeyPair = tls.LoadX509KeyPair
var appendCertsFromPEM = func(pool *x509.CertPool, pemCerts []byte) bool {
	return pool.AppendCertsFromPEM(pemCerts)
}
var registerTLSConfig = mysql.RegisterTLSConfig

// MariaDbExecutor is a concrete implementation of Executor for MariaDB
type MariaDbExecutor struct {
	db *sql.DB
}

// NewMariaDbExecutor creates a new MariaDB executor
func NewMariaDbExecutor(cfg *config.Config) (*MariaDbExecutor, error) {

	dsn := createDSN(cfg)
	dsn = setupTLSConfig(dsn, cfg)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	return &MariaDbExecutor{db: db}, nil

}

func createDSN(cfg *config.Config) string {
	escapedPassword := url.QueryEscape(cfg.DbPassword)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cfg.DbUser, escapedPassword, cfg.DbHost, cfg.DbPort, cfg.DbName)
	return dsn
}

// setupTLSConfig configures TLS using the specified certificates and registers it with the MySQL driver
func setupTLSConfig(dsn string, cfg *config.Config) string {

	if cfg.CACertPath == "" || cfg.ClientCertPath == "" || cfg.ClientKeyPath == "" {
		log.Println("certificate variables not set, bypassing tls setup")
		return dsn
	}

	// Load CA cert
	rootCertPool := x509.NewCertPool()
	pem, err := readFile(cfg.CACertPath)
	if err != nil {
		log.Panicf("failed to read CA cert file: %v", err)
	}
	if ok := appendCertsFromPEM(rootCertPool, pem); !ok {
		log.Panicf("failed to append CA cert")
	}

	// Load client cert
	clientCertPair, err := loadX509KeyPair(cfg.ClientCertPath, cfg.ClientKeyPath)
	if err != nil {
		log.Panicf("failed to load client cert and key: %v", err)
	}

	// Register a custom TLS config with the driver
	tlsConfig := &tls.Config{
		RootCAs:      rootCertPool,
		Certificates: []tls.Certificate{clientCertPair},
	}

	// Register and return the TLS config name
	const tlsConfigName = "dbTLS"
	if err := registerTLSConfig(tlsConfigName, tlsConfig); err != nil {
		log.Panicf("failed to register TLS config: %v", err)
	}

	// Append the TLS config to the DSN
	dsn += fmt.Sprintf("?tls=%s", tlsConfigName)
	return dsn
}

// Exec executes a query with the given arguments
func (e *MariaDbExecutor) Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := e.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	return result, nil
}

// Close closes the database connection
func (e *MariaDbExecutor) Close() error {
	return e.db.Close()
}
