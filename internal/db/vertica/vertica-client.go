package vertica

import (
	"context"
	"database/sql"
	vertigo "github.com/vertica/vertica-sql-go"
	"log"
	"os"
	"project-survey-stats-writer/internal/configuration"
)

type Client struct {
	db     *sql.DB
	config *configuration.DbConfiguration
}

func NewClient(config *configuration.DbConfiguration) *Client {
	return &Client{config: config}
}

func (cl *Client) Init() {
	db, err := sql.Open("vertica", cl.config.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	cl.db = db
}

func (cl *Client) BulkInsertFromCsv(fileName string, tableName string) error {
	fp, err := os.OpenFile(fileName, os.O_RDONLY, 0600)

	vCtx := vertigo.NewVerticaContext(context.Background())
	err = vCtx.SetCopyInputStream(fp)
	if err != nil {
		return err
	}

	_, err = cl.db.ExecContext(vCtx, "COPY @tableName FROM STDIN DELIMITER ','", sql.Named("tableName", tableName))
	return err
}

func (cl *Client) CloseConnection() {
	if cl.db != nil {
		cl.db.Close()
	}
}
