package stats

import (
	"log"
	"project-survey-stats-writer/internal/db/contracts"
	"project-survey-stats-writer/internal/events"
	"time"
)

var (
	batchSize           = 100
	maxBatchesPerSecond = 1000
	batchRateLimiter    = time.Tick(time.Second / time.Duration(maxBatchesPerSecond))
)

type StatsWriter struct {
	dbClient   contracts.IDbClient
	messageBuf []string
}

func GetStatsWriter(dbClient contracts.IDbClient) StatsWriter {
	return StatsWriter{dbClient: dbClient, messageBuf: make([]string, 0)}
}

func (sw *StatsWriter) AddMessageToWrite(message string) {
	sw.messageBuf = append(sw.messageBuf, message)

	if len(sw.messageBuf) >= batchSize {
		batch := make([]string, len(sw.messageBuf))
		copy(batch, sw.messageBuf)
		sw.messageBuf = nil

		<-batchRateLimiter

		go sw.processBatch(batch)
	}
}

func (sw *StatsWriter) Finalise() {
	if len(sw.messageBuf) > 0 {
		sw.processBatch(sw.messageBuf)
	}
}

func (sw *StatsWriter) processBatch(batch []string) {
	query := events.GetEventsInsertStatement(batch)

	err := sw.dbClient.BulkInsertFromCsv(query)
	if err != nil {
		log.Printf("Error executing statement: %v\n", err)
		return
	}
}
