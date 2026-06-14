// package postgresinfra

// import (
// 	"log"

// 	"github.com/jmoiron/sqlx"
// 	_ "github.com/lib/pq"
// )

// func NewDB(url string) *sqlx.DB {
// 	db, err := sqlx.Connect("postgres", url)
// 	if err != nil {
// 		log.Fatalf("Postgres connection failed: %v", err)
// 	}

// 	return db
// }

package postgresinfra

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDB(url string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("postgres connection failed: %w", err)
	}
	return db, nil
}