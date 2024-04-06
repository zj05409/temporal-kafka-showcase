// @@@SNIPSTART samples-go-cron-workflow
package cron

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

// SampleCronWorkflow executes on the given schedule
// The schedule is provided when starting the Workflow
func ProducerWorkflow(ctx workflow.Context) (error) {

	workflow.GetLogger(ctx).Info("Producer workflow started.", "StartTime", workflow.Now(ctx))

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Minute,
	}
	ctx1 := workflow.WithActivityOptions(ctx, ao)

	err := workflow.ExecuteActivity(ctx1, ProducerActivity).Get(ctx, nil)
	if err != nil {
		workflow.GetLogger(ctx).Error("Producer workflow failed.", "Error", err)
	}

    return err
}
func ConsumerWorkflow(ctx workflow.Context) (error) {

	workflow.GetLogger(ctx).Info("Consumer workflow started.", "StartTime", workflow.Now(ctx))

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 1 * time.Minute,
	}
	ctx1 := workflow.WithActivityOptions(ctx, ao)

	err := workflow.ExecuteActivity(ctx1, ConsumerActivity).Get(ctx, nil)
	if err != nil {
		workflow.GetLogger(ctx).Error("Consumer workflow failed.", "Error", err)
	}

    return err
}

// @@@SNIPEND