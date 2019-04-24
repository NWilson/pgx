package pool

import (
	"github.com/jackc/pgx/v4"
)

type errRows struct {
	err error
}

func (errRows) Close()                                    {}
func (e errRows) Err() error                              { return e.err }
func (errRows) FieldDescriptions() []pgx.FieldDescription { return nil }
func (errRows) Next() bool                                { return false }
func (e errRows) Scan(dest ...interface{}) error          { return e.err }
func (e errRows) Values() ([]interface{}, error)          { return nil, e.err }

type errRow struct {
	err error
}

func (e errRow) Scan(dest ...interface{}) error { return e.err }

type poolRows struct {
	r   pgx.Rows
	c   *Conn
	err error
}

func (rows *poolRows) Close() {
	rows.r.Close()
	if rows.c != nil {
		rows.c.Release()
		rows.c = nil
	}
}

func (rows *poolRows) Err() error {
	if rows.err != nil {
		return rows.err
	}
	return rows.r.Err()
}

func (rows *poolRows) FieldDescriptions() []pgx.FieldDescription {
	return rows.r.FieldDescriptions()
}

func (rows *poolRows) Next() bool {
	if rows.err != nil {
		return false
	}

	n := rows.r.Next()
	if !n {
		rows.Close()
	}
	return n
}

func (rows *poolRows) Scan(dest ...interface{}) error {
	err := rows.r.Scan(dest...)
	if err != nil {
		rows.Close()
	}
	return err
}

func (rows *poolRows) Values() ([]interface{}, error) {
	values, err := rows.r.Values()
	if err != nil {
		rows.Close()
	}
	return values, err
}

type poolRow struct {
	r   pgx.Row
	c   *Conn
	err error
}

func (row *poolRow) Scan(dest ...interface{}) error {
	if row.err != nil {
		return row.err
	}

	err := row.r.Scan(dest...)
	if row.c != nil {
		row.c.Release()
	}
	return err
}