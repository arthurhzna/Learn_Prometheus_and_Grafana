package postgres

import (
	"github.com/imrenagicom/demo-app/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

func NewSQLx(c config.SQL) *sqlx.DB {
	db, err := otelsqlx.Open("postgres", c.DataSourceName(),
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL))
	if err != nil {
		panic(err)
	}

	// Report DB stats metrics - INI YANG PENTING!
	otelsql.ReportDBStatsMetrics(db.DB,
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
	)

	db.SetMaxOpenConns(c.MaxOpenConn)
	db.SetMaxIdleConns(c.MaxIdleConn)

	return db
}
