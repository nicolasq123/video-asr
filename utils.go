package videoasr

import (
	"errors"
	"net/url"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
)

var ErrFileNotExisted = errors.New("file not existed error")
var ErrFileExisted = errors.New("file existed error")

var ErrDupMapKey = errors.New("duplicate map key error")

func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

func Open(dsn string) (*sqlx.DB, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "postgres" || u.Scheme == "postgresql" {
		db, err := sqlx.Connect("postgres", dsn)
		if err != nil {
			return nil, err
		}
		db = db.Unsafe()
		return db, nil
	}

	if u.Scheme == "mysql" {
		dsn = strings.TrimPrefix(dsn, "mysql://")

		db, err := sqlx.Connect("mysql", dsn)
		if err != nil {
			return nil, err
		}
		db = db.Unsafe()
		return db, nil
	}

	if u.Scheme == "clickhouse" {
		q := u.Query()
		if database := strings.Trim(u.Path, "/"); database != "" {
			q.Add("database", database)
		}
		dsn = "tcp://" + u.Host + "?" + q.Encode()
		db, err := sqlx.Connect("clickhouse", dsn)
		if err != nil {
			return nil, err
		}
		db = db.Unsafe()
		return db, nil
	}

	return nil, errors.New("unsurpport Scheme")
}

func isFileExisted(f string) bool {
	return fileStat(f) == nil
}

// nil就是存在
func fileStat(f string) error {
	_, err := os.Stat(f)
	return err
}

func min(x, y int) int {
	if x <= y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x <= y {
		return y
	}
	return x
}
