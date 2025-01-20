package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
)

// https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.html
const (
	ER_ACCESS_DENIED_ERROR = 1045
)

type Conn interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	Close() error
}

type Driver interface {
	CreateConnection() (Conn, error)
}

type Pool struct {
	db     Conn
	driver Driver
	mu     sync.RWMutex
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

func (p *Pool) reconnect() error {
	if !p.mu.TryLock() {
		return nil
	}
	defer p.mu.Unlock()

	db, err := p.driver.CreateConnection()
	if err != nil {
		return fmt.Errorf("could not create new connection: %w", err)
	}

	if p.db != nil {
		if err := p.db.Close(); err != nil {
			return fmt.Errorf("could not close connections: %w", err)
		}
	}

	p.db = db
	return nil
}

func (p *Pool) runWithRetry(op func() error) error {
	var errs []string
	timeout := time.Second
	maxNumberRetries := 3

	for i := 0; i < maxNumberRetries; i++ {
		p.mu.RLock()

		if err := op(); err != nil {
			if isConnectionError(err) {
				p.mu.RUnlock()
				if err := p.reconnect(); err != nil {
					return err
				}

				time.Sleep(timeout)
				timeout *= 2
				continue
			}

			p.mu.RUnlock()
			return fmt.Errorf("exec error: %w", err)
		}

		p.mu.RUnlock()
		return nil
	}

	errString := strings.Join(errs, ": ")
	return errors.New(errString)
}

func (p *Pool) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	var result sql.Result

	err := p.runWithRetry(func() error {
		var err error
		result, err = p.db.ExecContext(ctx, query, args...)
		return err
	})

	return result, err
}

func (p *Pool) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	var result *sql.Rows

	err := p.runWithRetry(func() error {
		var err error
		result, err = p.db.QueryContext(ctx, query, args...)
		return err
	})

	return result, err
}

func isConnectionError(err error) bool {
	if err, ok := err.(*mysql.MySQLError); ok {
		if err.Number == ER_ACCESS_DENIED_ERROR {
			return true
		}
	}

	return false
}
