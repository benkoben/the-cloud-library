package db

import (
	"database/sql"
    "context"
)

func IsAlive(ctx context.Context, conn *sql.DB) bool {
    if err := conn.PingContext(ctx); err != nil {
        return false
    } 
    return true
}
