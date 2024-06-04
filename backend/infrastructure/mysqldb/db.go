package mysqldb

import (
	"database/sql"

	"github.com/XSAM/otelsql"
	"github.com/go-sql-driver/mysql"
	semconv "go.opentelemetry.io/otel/semconv/v1.5.0"
)

func Initialize(user, pass, addr, name string) (*sql.DB, error) {
	c := mysql.Config{
		User:      user,
		Passwd:    pass,
		Net:       "tcp",
		Addr:      addr,
		DBName:    name,
		ParseTime: true,
	}
	db, err := otelsql.Open("mysql", c.FormatDSN(), otelsql.WithAttributes(
		semconv.DBSystemMySQL,
	), otelsql.WithSpanOptions(
		otelsql.SpanOptions{
			OmitConnQuery: true,
		},
	),
	)
	// db, err := sql.Open("mysql", c.FormatDSN())
	if err != nil {
		return nil, err
	}
	err = otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(
		semconv.DBSystemMySQL,
	))
	if err != nil {
		return nil, err
	}
	return db, nil
}
