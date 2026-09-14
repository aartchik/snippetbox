package models

import (
	"context"
	"database/sql"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func newTestDB(t *testing.T) *sql.DB {
	dsn := os.Getenv("SNIPPETBOX_TEST_DSN")
	if dsn == "" {
		dsn = "test_web:pass@tcp(localhost:3307)/test_snippetbox?parseTime=true&multiStatements=true&time_zone=%27%2B00%3A00%27"
	}

	db, err := sql.Open("mysql", dsn)
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
	addr := os.Getenv("SNIPPETBOX_TEST_REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6380"
	}

	db := 1
	if value := os.Getenv("SNIPPETBOX_TEST_REDIS_DB"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			t.Fatal(err)
		}
		db = parsed
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Fatal(err)
	}
	if err := rdb.FlushDB(ctx).Err(); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		if err := rdb.FlushDB(cleanupCtx).Err(); err != nil {
			t.Fatal(err)
		}
		rdb.Close()
	})
	return rdb
}
