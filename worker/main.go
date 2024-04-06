package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"zj05409.github.com/temporal_demo/cron"
)

func main() {

	connection := cron.InitKafka()
	defer connection.Close()

	reader := cron.InitKafkaConsumer()
	defer reader.Close()

	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{
		// HostPort: client.DefaultHostPort,
		HostPort: cron.TEMPORAL_FRONTEND,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "cron", worker.Options{})

	// w.RegisterWorkflow(cron.SampleCronWorkflow)
	// w.RegisterActivity(cron.DoSomething)
	w.RegisterActivity(cron.ProducerActivity)
	w.RegisterActivity(cron.ConsumerActivity)

	w.RegisterWorkflow(cron.ProducerWorkflow)
	w.RegisterWorkflow(cron.ConsumerWorkflow)

	w.RegisterWorkflow(cron.ProducersWorkflow)
	w.RegisterWorkflow(cron.ConsumersWorkflow)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
	log.Println("Done.")
}
