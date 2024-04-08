package cron

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

// CronResult is used to return data from one cron run to the next
type CronResult struct {
	RunTime time.Time
}

var (
	producerOnce       sync.Once
    kafkaWriter        *kafka.Conn
	consumerOnce sync.Once
	kafkaReader  *kafka.Reader
	Configs Config
)
type Config struct {
	ProducerCount         int         	 `env:"PRODUCER_COUNT,required"`
	ConsumerCount         int            `env:"CONSUMER_COUNT,required"`
	KafkaTopic            string         `env:"KAFKA_TOPIC,required"`
	KafkaProducers        string         `env:"KAFKA_PRODUCERS"`
	KafkaConsumers        string         `env:"KAFKA_CONSUMERS,required"`
	KafkaConsumerGroup    string         `env:"KAFKA_CONSUMER_GROUP,required"`
	TemporalFrontend      string         `env:"TEMPORAL_FRONTEND,required"`
}

func init() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
	
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	fmt.Printf("%+v\n", cfg)
	log.Fatal("Error loading .env file")
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
		kafkaWriter, err = kafka.DialLeader(context.Background(), "tcp", Configs.KafkaProducers, Configs.KafkaTopic, 0)
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
			Brokers:          strings.Split(Configs.KafkaConsumers, ","),
			GroupID:          Configs.KafkaConsumerGroup,
			Topic:            Configs.KafkaTopic,
			MaxBytes:         10e6, // 10MB
			MaxWait:          1 * time.Second,
			ReadBatchTimeout: 1 * time.Second,
			Dialer:           dialer,
		})

	})
	return kafkaReader
}
