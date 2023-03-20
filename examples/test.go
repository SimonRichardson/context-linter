package examples

import (
	c "context"
	"database/sql"
)

func Example() {
	var db *sql.DB
	db.Query("SELECT * FROM users")

	db.QueryContext(c.Background(), "SELECT * FROM users")
}
