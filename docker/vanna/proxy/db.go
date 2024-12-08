package proxy

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type Queryable interface {
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
}

type Scanable interface {
	Scan(dest ...interface{}) error
}

type Scanner[T any] func(Scanable) (T, error)

func ScanSlice[T any](scanner Scanner[T], rows *sql.Rows) ([]T, error) {
	results := make([]T, 0)
	for rows.Next() {
		x, err := scanner(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, x)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

type PGVersion struct {
	Version string
}

func scanPGVersion(s Scanable) (*PGVersion, error) {
	var e PGVersion
	if err := s.Scan(&e.Version); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &e, nil
}

func RetrievePGVersion(ctx context.Context, q Queryable) (*PGVersion, error) {
	return scanPGVersion(q.QueryRowContext(ctx, `
	SELECT version()
	`))
}

type PGDatabase struct {
	Datname string
}

func scanPGDatabase(s Scanable) (*PGDatabase, error) {
	var e PGDatabase
	if err := s.Scan(&e.Datname); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &e, nil
}

func RetrieveDatabases(ctx context.Context, q Queryable) ([]*PGDatabase, error) {
	const query = `
	SELECT datname FROM pg_database WHERE datistemplate = false AND datallowconn = true
	`
	rows, err := q.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*PGDatabase, 0)
	for rows.Next() {
		e, err := scanPGDatabase(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, e)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return results, nil
}

func Connect(db *DBInfo) (*sql.DB, error) {
	return sql.Open(db.DBName, db.DSN())
}

func Ping(db *DBInfo) error {
	pg, err := Connect(db)
	if err != nil {
		return err
	}
	defer pg.Close()
	return pg.Ping()
}
