package main

import (
	"context"
	"log"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"

	"zj05409.github.com/temporal_demo/cron"
	// "github.com/temporalio/samples-go/cron"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{
		// HostPort: client.DefaultHostPort,
		HostPort: cron.TEMPORAL_FRONTEND,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	// This workflow ID can be user business logic identifier as well.
	workflowOptions := client.StartWorkflowOptions{
		ID:                    "start_200_producer_per_minute",
		TaskQueue:             "cron",
		CronSchedule:          "*/1 * * * *",
		WorkflowIDReusePolicy: enums.WORKFLOW_ID_REUSE_POLICY_TERMINATE_IF_RUNNING,
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, cron.ProducersWorkflow)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	// setup 2 different cron workflows, Terminate old cron workflows if existing.
	workflowOptions = client.StartWorkflowOptions{
		ID:                    "start_2_consumer_every_5_minute",
		TaskQueue:             "cron",
		CronSchedule:          "*/5 * * * *",
		WorkflowIDReusePolicy: enums.WORKFLOW_ID_REUSE_POLICY_TERMINATE_IF_RUNNING,
	}
	we, err = c.ExecuteWorkflow(context.Background(), workflowOptions, cron.ConsumersWorkflow)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

}
