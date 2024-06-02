package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"project-survey-stats-writer/internal/configuration"
	"project-survey-stats-writer/internal/db/vertica"
	"project-survey-stats-writer/internal/events/kafka"
	"project-survey-stats-writer/internal/messages"
	"sync"
	"syscall"
)

func main() {
	parser := configuration.NewParser()
	config, err := parser.Parse("appsettings.json")
	if err != nil {
		log.Fatalf(err.Error())
		return
	}

	db := vertica.NewClient(config.DbConfiguration)
	err = db.Init()
	defer db.CloseConnection()
	if err != nil {
		log.Fatalf(err.Error())
		return
	}

	consumer := kafka.NewConsumer(config.EventsConfiguration)
	err = consumer.Init()
	defer consumer.CloseConnection()
	if err != nil {
		log.Fatalf(err.Error())
		return
	}

	proceeder := messages.NewProceeder(db, config.EventsConfiguration, consumer.Messages)
	go proceeder.Proceed()

	ctx, cancelCtx := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = consumer.Consume(ctx)
		if err != nil {
			log.Println(err.Error())
		}
	}()

	<-consumer.Ready
	if err != nil {
		log.Println(err.Error())
		return
	}

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	keepRunning := true
	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
		}
	}
	cancelCtx()
	wg.Wait()

	//statsWriter := stats.GetStatsWriter(db)

	//sigchan := make(chan os.Signal, 1)
	//signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	//
	//run := true
	//for run {
	//	select {
	//	case sig := <-sigchan:
	//		fmt.Printf("Caught signal %v: terminating\n", sig)
	//		run = false
	//	default:
	//		message, err := consumer.ConsumeMessage()
	//		if err != nil {
	//			log.Printf("Error when consuming message: %v\n", err)
	//			run = false
	//		}
	//
	//		if message != nil && len(message) != 0 {
	//			statsWriter.AddMessageToWrite(string(message))
	//		}
	//	}
	//}
	//
	//statsWriter.Finalise()
	//
	//fmt.Println("Stats writer terminated")
}
