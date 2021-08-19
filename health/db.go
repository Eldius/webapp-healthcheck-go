package health

import (
	"context"
	"database/sql"
	"time"
)

/*
DBServiceChecker is the checker config type for database
*/
type DBServiceChecker struct {
	name    string
	db      *sql.DB
	timeout time.Duration
	query   string
}

/*
Name returns the test/service name
*/
func (cfg *DBServiceChecker) Name() string {
	return cfg.name
}

/*
Type returns the test/service type (CheckerTypeDB)
*/
func (cfg *DBServiceChecker) Type() CheckerType {
	return CheckerTypeDB
}

/*
DB returns the DB
*/
func (cfg *DBServiceChecker) DB() *sql.DB {
	return cfg.db
}

/*
Timeout returns the test/service TCP test timeout
*/
func (cfg *DBServiceChecker) Timeout() time.Duration {
	return cfg.timeout
}

/*
Test returns the test/service status
*/
func (cfg *DBServiceChecker) Test() ServiceStatus {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, cfg.Timeout())
	defer cancel()
	start := time.Now()
	res := cfg.DB().QueryRowContext(ctx, "select 1")
	n := -1
	if err := res.Scan(&n); err != nil {
		return ServiceStatus{
			Name:   cfg.Name(),
			Status: CheckerStatusNOK,
			Details: map[string]string{
				"error": err.Error(),
			},
		}
	}
	return ServiceStatus{
		Name:   cfg.Name(),
		Status: CheckerStatusOK,
		Details: map[string]string{
			"time": time.Since(start).String(),
		},
	}
}

/*
NewDBChecker returns a new DB checker using the default query (`select 1`)
*/
func NewDBChecker(name string, db *sql.DB, timeout time.Duration) ServiceChecker {
	return NewDBCheckerCustomQuery(name, db, timeout, "select 1")
}

/*
NewDBCheckerCustomQuery returns a new DB checker using a query passed by parameter
*/
func NewDBCheckerCustomQuery(name string, db *sql.DB, timeout time.Duration, query string) ServiceChecker {
	return &DBServiceChecker{
		name:    name,
		db:      db,
		timeout: timeout,
		query:   query,
	}
}
