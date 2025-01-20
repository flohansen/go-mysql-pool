# MySQL Pool

A pool implementation for MySQL database connections which can be used for
failure resistant requirements like password rotations of the database.

## Usage

```go
package main

import (
	"database/sql"
	"fmt"
	"os"

	mysql "github.com/flohansen/go-mysql-pool"
)

type sqlDriver struct {
}

func (s *sqlDriver) CreateConnection() (mysql.Conn, error) {
    // Implement your logic on how to create a new connection here. It will be
    // called inside the retry logic of the pool whenever a query returns a error
    // regarding to the connection.

	return sql.Open("mysql", "username:password@tcp(host:port)/dbname")
}

func main() {
	driver := &sqlDriver{}
	pool, err := mysql.NewPool(driver)
	if err != nil {
		fmt.Printf("could not create pool: %s", err)
		os.Exit(1)
	}

	// use the pool like sql.DB
	result, err := pool.ExecContext(ctx, "SELECT * FROM users WHERE username = ?", username)
}
```
