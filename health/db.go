package health

import (
	"context"
	"database/sql"
	"time"
)

/*
TCPServiceConfig is the checker config type for TCP checks
*/
type DBServiceConfig struct {
	name    string
	db      *sql.DB
	timeout time.Duration
	query   string
}

/*
Name returns the test/service name
*/
func (cfg *DBServiceConfig) Name() string {
	return cfg.name
}

/*
Type returns the test/service type (ServiceTypeTCP)
*/
func (cfg *DBServiceConfig) Type() ServiceType {
	return ServiceTypeDB
}

/*
Endpoint returns the test/service endpoint
*/
func (cfg *DBServiceConfig) DB() *sql.DB {
	return cfg.db
}

/*
Timeout returns the test/service TCP test timeout
*/
func (cfg *DBServiceConfig) Timeout() time.Duration {
	return cfg.timeout
}

/*
Test returns the test/service status
*/
func (cfg *DBServiceConfig) Test() Status {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, cfg.Timeout())
	defer cancel()
	start := time.Now()
	res := cfg.db.QueryRowContext(ctx, "select 1")
	n := -1
	if err := res.Scan(&n); err != nil {
		return Status{
			Name:   cfg.Name(),
			Status: ServiceStatusNOK,
			Details: map[string]string{
				"error": err.Error(),
			},
		}
	}
	return Status{
		Name:   cfg.Name(),
		Status: ServiceStatusNOK,
		Details: map[string]string{
			"time": time.Since(start).String(),
		},
	}
}

func NewDBChecker(name string, db *sql.DB, timeout time.Duration) ServiceConfig {
	return NewDBCheckerCustomQuery(name, db, timeout, "select 1")
}

func NewDBCheckerCustomQuery(name string, db *sql.DB, timeout time.Duration, query string) ServiceConfig {
	return &DBServiceConfig{
		name:    name,
		db:      db,
		timeout: timeout,
		query:   query,
	}
}
