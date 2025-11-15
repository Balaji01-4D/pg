package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"pg/internal/pg"
)

// Example program demonstrating the usage of the pg package
// FOR DEMONSTRATION PURPOSES ONLY

func main() {
	ctx := context.Background()

	// 1. Connect to Postgres using your Connect() function
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	if user == "" || pass == "" {
		log.Fatal("DB_USER and DB_PASSWORD must be set in the environment")
	}
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:5432/uber?sslmode=disable", user, pass)
	pool, err := pg.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("cannot connect: %v", err)
	}
	defer pool.Close()

	// 2. Create Executor
	executor := pg.NewExecutor(pool)

	// 3. Run examples
	runExecExample(ctx, *executor)
	insertRow(ctx, *executor)
	runQueryExample(ctx, *executor)

}

func runExecExample(ctx context.Context, e pg.Excutor) {
	sql := `CREATE TABLE IF NOT EXISTS demo (
        id SERIAL PRIMARY KEY,
        name TEXT
    )`

	result, err := e.Exec(ctx, sql)
	if err != nil {
		log.Fatalf("exec error: %v", err)
	}

	fmt.Println("Exec status:", result.Status)
	fmt.Println("Rows affected:", result.RowsAffected)
	fmt.Println("Duration:", result.Duration)
}

func insertRow(ctx context.Context, e pg.Excutor) {
	result, err := e.Exec(ctx, "INSERT INTO demo(name) VALUES($1)", "Balaji")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted rows:", result.RowsAffected)
}
func runQueryExample(ctx context.Context, e pg.Excutor) {
	rs, err := e.Query(ctx, "select * from users")
	if err != nil {
		log.Fatalf("query error: %v", err)
	}
	defer rs.Close()

	fmt.Println("Columns:", rs.Columns(), "\nDuration:", rs.Duration())
	for {
		row, err := rs.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("row error: %v", err)
		}
		fmt.Println("Row:", row)
	}
}
