package models

import (
	"database/sql"
	"os"
	"testing"
    "context"
    "time"

	"github.com/redis/go-redis/v9"
)

func newTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("mysql",
  "test_web:pass@/test_snippetbox?parseTime=true&multiStatements=true&time_zone=%27%2B00%3A00%27",
)
    if err != nil {
        t.Fatal(err)
    }

    script, err := os.ReadFile("./testdata/setup.sql")
    if err != nil {
        t.Fatal(err)
    }
    _, err = db.Exec(string(script))
    if err != nil {
        t.Fatal(err)
    }

    t.Cleanup(func() {
        script, err := os.ReadFile("./testdata/teardown.sql")
        if err != nil {
            t.Fatal(err)
        }
        _, err = db.Exec(string(script))
        if err != nil {
            t.Fatal(err)
        }

        db.Close()
    })

    return db
}

func newTestRedis(t *testing.T) *redis.Client {
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
        DB: 1,
    })
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Fatal(err)
	}
    t.Cleanup(func() {
        rdb.Close()
    })
    return rdb
}