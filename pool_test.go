package mysql_test

import (
	"context"
	"errors"
	"testing"

	"github.com/flohansen/go-mysql-pool"
	"github.com/flohansen/go-mysql-pool/mocks"
	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type sqlResult struct {
}

func (r sqlResult) LastInsertId() (int64, error) {
	return 0, nil
}

func (r sqlResult) RowsAffected() (int64, error) {
	return 0, nil
}

//go:generate mockgen -destination=mocks/conn.go -package=mocks github.com/flohansen/go-mysql-pool Conn
//go:generate mockgen -destination=mocks/driver.go -package=mocks github.com/flohansen/go-mysql-pool Driver

func TestPool_ExecContext(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Run("should return error immediately if error is not a connection error", func(t *testing.T) {
		// given
		ctx := context.TODO()

		driver := mocks.NewMockDriver(ctrl)
		conn := mocks.NewMockConn(ctrl)

		driver.EXPECT().
			CreateConnection().
			Return(conn, nil)

		pool, err := mysql.NewPool(driver)
		if err != nil {
			t.Fatal(err)
		}

		conn.EXPECT().
			ExecContext(ctx, "query").
			Return(nil, &mysqldriver.MySQLError{
				Number:  0,
				Message: "any other error",
			})

		// when
		result, err := pool.ExecContext(ctx, "query")

		// then
		assert.Error(t, err, "any other error")
		assert.Nil(t, result)
	})

	t.Run("should retry query execution if connection has problems", func(t *testing.T) {
		// given
		ctx := context.TODO()

		driver := mocks.NewMockDriver(ctrl)
		conn := mocks.NewMockConn(ctrl)

		driver.EXPECT().
			CreateConnection().
			Return(conn, nil)

		pool, err := mysql.NewPool(driver)
		if err != nil {
			t.Fatal(err)
		}

		driver.EXPECT().
			CreateConnection().
			Return(conn, nil)

		conn.EXPECT().
			ExecContext(ctx, "query").
			Return(nil, &mysqldriver.MySQLError{
				Number:  mysql.ER_ACCESS_DENIED_ERROR,
				Message: "connection error",
			})

		conn.EXPECT().
			ExecContext(ctx, "query").
			Return(&sqlResult{}, nil)

		conn.EXPECT().
			Close().
			Return(nil)

		// when
		result, err := pool.ExecContext(ctx, "query")

		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("should return error if connection has problems and connection cannot be recreated", func(t *testing.T) {
		// given
		ctx := context.TODO()

		driver := mocks.NewMockDriver(ctrl)
		conn := mocks.NewMockConn(ctrl)

		driver.EXPECT().
			CreateConnection().
			Return(conn, nil)

		pool, err := mysql.NewPool(driver)
		if err != nil {
			t.Fatal(err)
		}

		conn.EXPECT().
			ExecContext(ctx, "query").
			Return(nil, &mysqldriver.MySQLError{
				Number:  mysql.ER_ACCESS_DENIED_ERROR,
				Message: "connection error",
			})

		driver.EXPECT().
			CreateConnection().
			Return(nil, errors.New("error"))

		// when
		result, err := pool.ExecContext(ctx, "query")

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("should retry query execution if connection has problems and return error if every retry failed", func(t *testing.T) {
		// given
		ctx := context.TODO()

		driver := mocks.NewMockDriver(ctrl)
		conn := mocks.NewMockConn(ctrl)

		driver.EXPECT().
			CreateConnection().
			Return(conn, nil)

		pool, err := mysql.NewPool(driver)
		if err != nil {
			t.Fatal(err)
		}

		driver.EXPECT().
			CreateConnection().
			Return(conn, nil).
			Times(3)

		conn.EXPECT().
			ExecContext(ctx, "query").
			Return(nil, &mysqldriver.MySQLError{
				Number:  mysql.ER_ACCESS_DENIED_ERROR,
				Message: "connection error",
			}).
			Times(3)

		conn.EXPECT().
			Close().
			Return(nil).
			Times(3)

		// when
		result, err := pool.ExecContext(ctx, "query")

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("should execute query with context", func(t *testing.T) {
		// given
		ctx := context.TODO()

		driver := mocks.NewMockDriver(ctrl)
		conn := mocks.NewMockConn(ctrl)

		driver.EXPECT().
			CreateConnection().
			Return(conn, nil)

		pool, err := mysql.NewPool(driver)
		if err != nil {
			t.Fatal(err)
		}

		conn.EXPECT().
			ExecContext(ctx, "query").
			Return(&sqlResult{}, nil)

		// when
		result, err := pool.ExecContext(ctx, "query")

		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}
