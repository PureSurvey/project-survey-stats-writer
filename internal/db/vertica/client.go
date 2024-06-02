package vertica

import (
	"context"
	"database/sql"
	vertigo "github.com/vertica/vertica-sql-go"
	"log"
	"os"
	"project-survey-stats-writer/internal/configuration"
	"time"
)

type Client struct {
	db     *sql.DB
	config *configuration.DbConfiguration
}

func NewClient(config *configuration.DbConfiguration) *Client {
	return &Client{config: config}
}

func (cl *Client) Init() error {
	var err error
	cl.db, err = sql.Open("vertica", cl.config.ConnectionString)
	if err != nil {
		log.Fatal(err)
		return err
	}

	cl.db.SetMaxOpenConns(10)
	cl.db.SetMaxIdleConns(5)

	for i := 0; i < cl.config.ConnectionRetryCount; i++ {
		err = cl.db.Ping()
		if err != nil {
			log.Println("Error when connecting to DB:", err.Error(), "Try", i+1, "of", cl.config.ConnectionRetryCount)
			time.Sleep(time.Duration(cl.config.ConnectionRetryTimeout) * time.Second)
		}
	}
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
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
