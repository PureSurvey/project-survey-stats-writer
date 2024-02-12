package vertica

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/vertica/vertica-sql-go"
	"log"
)

type VerticaClient struct {
	dbClient *sql.DB
}

func InitClient() *VerticaClient {
	db, err := sql.Open("vertica", "vertica://dbadmin:@localhost:5433?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to Vertica database")
	return &VerticaClient{db}
}

func (cl *VerticaClient) InsertRows(query string) error {
	vCtx := vertigo.NewVerticaContext(context.Background())

	_, err := cl.dbClient.ExecContext(vCtx, query)
	return err
}

func (cl *VerticaClient) CloseConnection() {
	cl.dbClient.Close()
}
