package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"project-survey-stats-writer/kafka"
	"project-survey-stats-writer/stats"
	"project-survey-stats-writer/vertica"
	"project-survey-stats-writer/internal/configuration"
	"syscall"
)

func main() {
	db := vertica.InitClient()
	parser := configuration.NewParser()
	config, err := parser.Parse("appsettings.json")
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
	defer db.CloseConnection()

	consumer, err := kafka.InitConsumer("localhost:9092")
	if err != nil {
		log.Fatalf("Error creating message consumer: %v", err)
		os.Exit(1)
	}
	defer consumer.CloseConnection()

	statsWriter := stats.GetStatsWriter(db)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			message, err := consumer.ConsumeMessage()
			if err != nil {
				log.Printf("Error when consuming message: %v\n", err)
				run = false
			}

			if message != nil && len(message) != 0 {
				statsWriter.AddMessageToWrite(string(message))
			}
		}
	}

	statsWriter.Finalise()

	fmt.Println("Stats writer terminated")
}
