package cron

import (
	"context"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

// CronResult is used to return data from one cron run to the next
type CronResult struct {
	RunTime time.Time
}

var (
    producerCount      = getEnv("PRODUCER_COUNT", "20")
    consumerCount      = getEnv("CONSUMER_COUNT", "2")
    kafkaTopic         = getEnv("KAFKA_TOPIC", "preonline")
    kafkaProducers     = getEnv("KAFKA_PRODUCERS", "kafka.kafka.svc.cluster.local:9092")
    kafkaConsumers     = getEnv("KAFKA_CONSUMERS", "kafka.kafka.svc.cluster.local:9092")
    kafkaConsumerGroup = getEnv("KAFKA_CONSUMER_GROUP", "consumer-group-id")
    temporalFrontend   = getEnv("TEMPORAL_FRONTEND", "temporaltest-frontend.default.svc.cluster.local:7233")
    producerOnce       sync.Once
    kafkaWriter        *kafka.Conn
	consumerOnce sync.Once
	kafkaReader  *kafka.Reader
)

func init() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
}

func getEnv(key, defaultValue string) string {
    value, exists := os.LookupEnv(key)
    if !exists {
        value = defaultValue
    }
    return value
}

func InitKafka() *kafka.Conn {
	producerOnce.Do(func() {
		var err error
		kafkaWriter, err = kafka.DialLeader(context.Background(), "tcp", kafkaProducers, kafkaTopic, 0)
		if err != nil {
			log.Fatal("failed to dial leader:", err)
		}
	})
	// kafkaWriter.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return kafkaWriter
}
func InitKafkaConsumer() *kafka.Reader {
	consumerOnce.Do(func() {

		dialer := &kafka.Dialer{
			Timeout: 120 * time.Second,
		}

		kafkaReader = kafka.NewReader(kafka.ReaderConfig{
			Brokers:          strings.Split(kafkaConsumers, ","),
			GroupID:          kafkaConsumerGroup,
			Topic:            kafkaTopic,
			MaxBytes:         10e6, // 10MB
			MaxWait:          1 * time.Second,
			ReadBatchTimeout: 1 * time.Second,
			Dialer:           dialer,
		})

	})
	return kafkaReader
}
