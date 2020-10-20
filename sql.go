package otsql

/**
 * Open telemetry Wrapper for database/sql package. Does not implement all functionality yet
 */

import (
	"context"
	"database/sql"
	"database/sql/driver"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
)

// DB internal struct
type DB struct {
	sqldb *sql.DB
}

// Open opens a database connection
func Open(ctx context.Context, driverName, dataSourceName string) (*DB, error) {
	_, span := startSpan(ctx, "sql.Open", nil)
	defer span.End()

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{sqldb: db}, nil
}

// OpenDB  open a db connection us the given string
func OpenDB(ctx context.Context, c driver.Connector) *DB {
	_, span := startSpan(ctx, "sql.OpenDB", nil)
	defer span.End()

	db := sql.OpenDB(c)
	return &DB{sqldb: db}
}

// Close closes the connction
func (db *DB) Close(ctx context.Context) error {
	_, span := startSpan(ctx, "sql.DB.Close", nil)
	defer span.End()

	return db.sqldb.Close()
}

// Query execute a sql
func (db *DB) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	attribs := map[string]string{
		"sql": query,
	}
	_, span := startSpan(ctx, "sql.DB.Query", attribs)
	defer span.End()

	return db.sqldb.Query(query, args...)
}

// Prepare prepare a query
func (db *DB) Prepare(ctx context.Context, query string) (*sql.Stmt, error) {
	attribs := map[string]string{
		"sql": query,
	}
	_, span := startSpan(ctx, "sql.DB.Prepare", attribs)
	defer span.End()

	return db.sqldb.Prepare(query)
}

func startSpan(ctx context.Context, name string, attribs map[string]string) (context.Context, trace.Span) {
	tropts := []trace.SpanOption{}
	// tropts = append(tropts, opts...)
	tropts = append(tropts, trace.WithSpanKind(trace.SpanKindInternal))
	tr := global.Tracer("otsql-wrapper")
	ctx, span := tr.Start(ctx, name, tropts...)

	// add headers from request as span attributes
	spanAttribs := []label.KeyValue{}
	for k, v := range attribs {
		spanAttribs = append(spanAttribs, label.String(k, v))
	}
	span.SetAttributes(spanAttribs...)

	return ctx, span
}
