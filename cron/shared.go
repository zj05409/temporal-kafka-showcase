package cron

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/go-baa/pool"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

// CronResult is used to return data from one cron run to the next
type CronResult struct {
	RunTime time.Time
}

var (
	producerOnce       sync.Once
	kafkaReader  *kafka.Reader
	Configs Config
)
type Config struct {
	ProducerCount         int         	 `env:"PRODUCER_COUNT,required"`
	ConsumerCount         int            `env:"CONSUMER_COUNT,required"`
	KafkaTopic            string         `env:"KAFKA_TOPIC,required"`
	KafkaProducers        string         `env:"KAFKA_PRODUCERS,required"`
	KafkaConsumers        string         `env:"KAFKA_CONSUMERS,required"`
	KafkaConsumerGroup    string         `env:"KAFKA_CONSUMER_GROUP,required"`
	TemporalFrontend      string         `env:"TEMPORAL_FRONTEND,required"`
}

func init() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file", err)
    }
	
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("Error parsing .env file", err)
	} else {
		Configs = cfg
	}

	fmt.Printf("%+v\n", cfg)
}

var (
	kafkaWriterPool *pool.Pool
	err error;
)

func InitKafka() *pool.Pool {
	producerOnce.Do(func() {
		kafkaWriterPool, err = pool.New(10, 100, func() interface{} {
			kafkaWriter, err := kafka.DialLeader(context.Background(), "tcp", Configs.KafkaProducers, Configs.KafkaTopic, 0)
			if err != nil {
				log.Fatal("failed to dial leader:", err)
			}
			return kafkaWriter
		})
		if err != nil {
			log.Fatalf("create pool error: %v\n", err)
		}
		// len
		log.Printf("total connections: %d\n", kafkaWriterPool.Len())

		kafkaWriterPool.Ping = func(conn interface{}) bool {
			// check connection status
			brokers, error := conn.(*kafka.Conn).Brokers()
			return error != nil && brokers != nil && len(brokers) > 0
		}
	
		kafkaWriterPool.Close = func(conn interface{}) {
			// close connection
			conn.(*kafka.Conn).Close()
		}
	})
	return kafkaWriterPool
}

func CloseKafka() {
	// len
	log.Printf("total connections: %d\n", kafkaWriterPool.Len())

	// destroy, close all connections
	kafkaWriterPool.Destroy()
}

func GetKafkaWriter() *kafka.Conn {
	kafkaConn, err := kafkaWriterPool.Get()
	if (err != nil) {
		log.Fatalf("get connection error: %v\n", err)
	}
	return kafkaConn.(*kafka.Conn)
}

func ReturnKafkaWriter(conn *kafka.Conn) {
	kafkaWriterPool.Put(conn)
}
func InitKafkaConsumer() *kafka.Reader {

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
	return kafkaReader
}
