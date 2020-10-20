package main

import (
	"context"
	"log"

	"k-shaikh/golang/otsql"

	_ "github.com/go-sql-driver/mysql"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/label"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func initTracer() func() {
	// Create and install Jaeger export pipeline
	flush, err := jaeger.InstallNewPipeline(
		jaeger.WithCollectorEndpoint("http://localhost:14268/api/traces"),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: "ot-sql-demo",
			Tags: []label.KeyValue{
				label.String("exporter", "jaeger"),
				label.Float64("float", 312.23),
			},
		}),
		jaeger.WithSDK(&sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
	)
	if err != nil {
		log.Fatal(err)
	}

	return func() {
		flush()
	}
}

func main() {
	fn := initTracer()
	defer fn()

	ctx := context.Background()
	db, err := otsql.Open(ctx, "mysql",
		"<usename>:<pwd>@tcp(127.0.0.1:3306)/<dbname>")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(ctx)

	retrieve(db)
}

func retrieve(db *otsql.DB) {
	tr := global.Tracer("sample-go-server")
	ctx := context.Background()
	ctx, span := tr.Start(ctx, "Retrieve")
	defer span.End()
	var (
		currentTime string
	)
	rows, err := db.Query(ctx, "select now() as currentTime from dual Where 1 = ? ", 1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&currentTime)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(currentTime)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
