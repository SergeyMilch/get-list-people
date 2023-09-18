package db

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Database interface {
	Select(dest interface{}, query string, args ...interface{}) error
	Get(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	NamedExec(query string, arg interface{}) (sql.Result, error)
}

type DB struct {
	db *sqlx.DB
}

func NewDatabase(db *sqlx.DB) Database {
	return &DB{db: db}
}

func (c *DB) Select(dest interface{}, query string, args ...interface{}) error {
	return c.db.Select(dest, query, args...)
}

func (c *DB) Get(dest interface{}, query string, args ...interface{}) error {
	return c.db.Get(dest, query, args...)
}

func (c *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return c.db.Exec(query, args...)
}

func (c *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	return c.db.QueryRow(query, args...)
}

func (c *DB) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return c.db.NamedExec(query, arg)
}
