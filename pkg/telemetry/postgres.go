package telemetry

import (
	"context"
	"fmt"
	"strings"

	"github.com/xo/dburl"
	"go.opencensus.io/trace"
)

type PostgresSpan struct {
	dbSystem string
	dbName   string
	dbSchema string
	dbUser   string
	dbAddr   string
	dbPort   string
}

func InitPostgresSpan(dbSchema string, url string) (*PostgresSpan, error) {
	u, err := dburl.Parse(url)
	if err != nil {
		return nil, err
	}
	return &PostgresSpan{
		dbSystem: "postgresql",
		dbName:   strings.TrimPrefix(u.EscapedPath(), "/"),
		dbSchema: dbSchema,
		dbUser:   u.User.Username(),
		dbAddr:   u.Hostname(),
		dbPort:   u.Port(),
	}, err
}

func (d PostgresSpan) StartSpan(ctx context.Context, op string, tableName string, spanAttributes map[string]string) (context.Context, *trace.Span) {
	// Refer https://github.com/open-telemetry/opentelemetry-specification/blob/master/specification/trace/semantic_conventions/database.md
	ctx, span := trace.StartSpan(ctx, fmt.Sprintf("%s %s.%s.%s", op, d.dbName, d.dbSchema, tableName), trace.WithSpanKind(trace.SpanKindClient))

	traceAttributes := []trace.Attribute{
		trace.StringAttribute("db.system", d.dbSystem),
		trace.StringAttribute("db.user", d.dbUser),
		trace.StringAttribute("net.sock.peer.addr", d.dbAddr),
		trace.StringAttribute("net.peer.port", d.dbPort),
		trace.StringAttribute("db.name", d.dbName),
		trace.StringAttribute("db.operation", op),
		trace.StringAttribute("db.sql.table", tableName),
	}

	for k, v := range spanAttributes {
		traceAttributes = append(traceAttributes, trace.StringAttribute(k, v))
	}

	span.AddAttributes(
		traceAttributes...,
	)

	return ctx, span
}
