package cron

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

// CronResult is used to return data from one cron run to the next
type CronResult struct {
	RunTime time.Time
}

const PRODUCER_COUNT = 200
const CONSUMER_COUNT = 2

// const KAFKA_TOPIC = "kafka_activity_test"
const KAFKA_TOPIC = "preonline"

// const KAFKA_PROCUCERS = "localhost:9092"
// const KAFKA_CONSUMERS = "localhost:9092"

const KAFKA_PROCUCERS = "kafka-controller-0.kafka-controller-headless.default.svc.cluster.local:9092"
const KAFKA_CONSUMERS = "kafka.default.svc.cluster.local:9092"
const KAFKA_CONSUMER_GROUP = "consumer-group-id"

// const KAFKA_USERNAME = "user1"
// const KAFKA_PASSWORD = "nyo0WH8B78"

const TEMPORAL_FRONTEND = "temporaltest-frontend.default.svc.cluster.local:7233"

var (
	procucerOnce sync.Once
	kafkaWriter  *kafka.Conn

	consumerOnce sync.Once
	kafkaReader  *kafka.Reader
)

func InitKafka() *kafka.Conn {
	procucerOnce.Do(func() {
		var err error
		kafkaWriter, err = kafka.DialLeader(context.Background(), "tcp", KAFKA_PROCUCERS, KAFKA_TOPIC, 0)
		if err != nil {
			log.Fatal("failed to dial leader:", err)
		}
	})
	kafkaWriter.SetWriteDeadline(time.Now().Add(3 * time.Second))
	return kafkaWriter
}
func InitKafkaConsumer() *kafka.Reader {
	consumerOnce.Do(func() {

		dialer := &kafka.Dialer{
			Timeout: 120 * time.Second,
		}

		kafkaReader = kafka.NewReader(kafka.ReaderConfig{
			Brokers:          strings.Split(KAFKA_CONSUMERS, ","),
			GroupID:          KAFKA_CONSUMER_GROUP,
			Topic:            KAFKA_TOPIC,
			MaxBytes:         10e6, // 10MB
			MaxWait:          1 * time.Second,
			ReadBatchTimeout: 1 * time.Second,
			Dialer:           dialer,
		})

	})
	return kafkaReader
}
