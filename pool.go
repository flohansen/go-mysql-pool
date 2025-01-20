package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

// https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.html
const (
	ER_ACCESS_DENIED_ERROR = 1045
)

type Conn interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type Driver interface {
	CreateConnection() (Conn, error)
}

type Pool struct {
	db     Conn
	driver Driver
}

func NewPool(driver Driver) (*Pool, error) {
	db, err := driver.CreateConnection()
	if err != nil {
		return nil, fmt.Errorf("could not open sql: %w", err)
	}

	return &Pool{
		db:     db,
		driver: driver,
	}, nil
}

func (p *Pool) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	var errs []string
	timeout := time.Second
	maxNumberRetries := 3

	for i := 0; i < maxNumberRetries; i++ {
		result, err := p.db.ExecContext(ctx, query, args...)
		if err != nil {
			if isConnectionError(err) {
				errs = append(errs, err.Error())
				time.Sleep(timeout)
				timeout *= 2
				continue
			}

			return nil, fmt.Errorf("exec error: %w", err)
		}

		return result, nil
	}

	errString := strings.Join(errs, ": ")
	return nil, errors.New(errString)
}

func isConnectionError(err error) bool {
	if err, ok := err.(*mysql.MySQLError); ok {
		if err.Number == ER_ACCESS_DENIED_ERROR {
			return true
		}
	}

	return false
}
