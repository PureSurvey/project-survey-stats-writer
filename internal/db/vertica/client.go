package vertica

import (
	"context"
	"database/sql"
	"fmt"
	vertigo "github.com/vertica/vertica-sql-go"
	"io"
	"log"
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

func (cl *Client) BulkInsert(inputStream io.Reader, tableName string) error {
	vCtx := vertigo.NewVerticaContext(context.Background())
	err := vCtx.SetCopyInputStream(inputStream)
	if err != nil {
		return err
	}

	_, err = cl.db.ExecContext(vCtx, fmt.Sprintf(`COPY %v FROM STDIN DELIMITER ',' ENCLOSED '''';`, tableName))
	return err
}

func (cl *Client) CloseConnection() {
	if cl.db != nil {
		cl.db.Close()
	}
}
