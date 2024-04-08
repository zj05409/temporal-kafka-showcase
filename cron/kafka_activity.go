package cron

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"go.temporal.io/sdk/activity"
)

// DoSomething is an Activity
func ProducerActivity(ctx context.Context) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Producing Message")

	conn := GetKafkaWriter()
	defer ReturnKafkaWriter(conn)

	_, err := conn.WriteMessages(
		kafka.Message{
			Key:   []byte("Key-A"),
			Value: []byte("Hello World!"),
		},
	)
	if err != nil {
		logger.Error("failed to write messages:", err)
	} else {
		logger.Info("write messages succeed!:")
	}
	return err
}

func ConsumerActivity(ctx context.Context) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Consuming Message")
	r := InitKafkaConsumer()
	defer r.Close()

	var m kafka.Message
	var err error
	for {
		c := make(chan bool)
		go func() {
			m, err = r.FetchMessage(ctx)
			c <- true
		}()
		select {
		case <-c:
			if err != nil {
				logger.Error("failed to fetch message:", err)
				return err
			}
			logger.Debug(fmt.Sprintf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value)))
			if err := r.CommitMessages(ctx, m); err != nil {
				logger.Error("failed to commit messages:", err)
				return err
			}
			logger.Info("commit messages succeed!:")
		case <-time.After(10 * time.Second):
			logger.Debug("timed out waiting for fetches")
			return nil
		}
	}

}
