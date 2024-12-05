package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

func ConnectWorker(brokersUrl []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	conn, err := sarama.NewConsumer(brokersUrl, config)
	if err != nil {
		return nil, err
	}
	return conn, nil

}

func main() {
	brokersUrl := []string{"localhost:29092"}
	worker, err := ConnectWorker(brokersUrl)
	if err != nil {
		panic(err)
	}
	topic := "comments"
	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}
	fmt.Println("Consumer Started")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	msgCount := 0
	donechn := make(chan struct{})

	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)

			case msg := <-consumer.Messages():
				msgCount++
				fmt.Printf("Received message Count: %d, Topic: %s, Message %s \n", msgCount, string(msg.Topic), string(msg.Value))

			case <-sigchan:
				fmt.Println("Interruption Detected")
				donechn <- struct{}{}
				return
			}
		}

	}()

	<-donechn
	fmt.Println("Processed", msgCount, "messages")

}
