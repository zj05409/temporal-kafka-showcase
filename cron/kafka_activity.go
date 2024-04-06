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

	conn := InitKafka()
	// w := &kafka.Writer{
	// 	// Addr:     kafka.TCP("localhost:9092", "localhost:9093", "localhost:9094"),
	// 	Addr:     kafka.TCP("kafka-controller-0.kafka-controller-headless.default.svc.cluster.local:9092","kafka-controller-1.kafka-controller-headless.default.svc.cluster.local:9092",
	// 		"kafka-controller-2.kafka-controller-headless.default.svc.cluster.local:9092"),
	// 	Topic:   KAFKA_TOPIC,
	// 	// Balancer: &kafka.LeastBytes{},
	// 	// Transport: &kafka.Transport{
	// 	// 	SASL:mechanism,
	// 	// },
	// 	BatchTimeout: 1 * time.Second,
	// 	BatchSize: 1,
	// }
	// defer w.Close()
	// _, err := w.WriteMessages(context.Background(),

	// defer conn.Close()
	_, err := conn.WriteMessages(
		kafka.Message{
			Key:   []byte("Key-A"),
			Value: []byte("Hello World!"),
		},
		kafka.Message{
			Key:   []byte("Key-B"),
			Value: []byte("One!"),
		},
		kafka.Message{
			Key:   []byte("Key-C"),
			Value: []byte("Two!"),
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
