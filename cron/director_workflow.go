package cron

import (
	"errors"
	"fmt"

	"go.temporal.io/sdk/workflow"
)

func ProducersWorkflow(ctx workflow.Context) error {
	return batchExecuteChildWorkflowGroup(ctx, ProducerWorkflow, "producer", Configs.ProducerCount)

}
func ConsumersWorkflow(ctx workflow.Context) error {
	// GetChildWorkflowExecution()
	return batchExecuteChildWorkflowGroup(ctx, ConsumerWorkflow, "consumer", Configs.ConsumerCount)
}

func batchExecuteChildWorkflowGroup(ctx workflow.Context, childWorkflow interface{}, childName string, childCount int) error {
	logger := workflow.GetLogger(ctx)
	logger.Info(fmt.Sprintf("%s workflow group started", childName), "StartTime", workflow.Now(ctx))

	ao := workflow.ChildWorkflowOptions{}
	ctx1 := workflow.WithChildOptions(ctx, ao)

	var futures = make([]workflow.ChildWorkflowFuture, childCount)
	for i := 0; i < childCount; i++ {
		futures[i] = workflow.ExecuteChildWorkflow(ctx1, childWorkflow)

	}
	workflow.Await(ctx, func() bool {
		for _, future := range futures {
			if !future.GetChildWorkflowExecution().IsReady() {
				return false
			}
		}
		return true
	})
	for _, future := range futures {
		err := future.GetChildWorkflowExecution().Get(ctx, nil)
		if err != nil {
			return errors.Join(errors.New("ChildWorkflowStartError"), err)
		}
	}
	workflow.Await(ctx, func() bool {
		for _, future := range futures {
			if !future.IsReady() {
				return false
			}
		}
		return true
	})

	errorCount := 0
	for i, future := range futures {
		err := future.Get(ctx, nil)
		if err != nil {
			logger.Error("Child %s workflow failed.", "Error", err, "Index", i)
			errorCount++
		}
	}

	if errorCount == 0 {
		logger.Info(fmt.Sprintf("Child %s workflows all succeeded.", childName))
		return nil
	}
	if errorCount < childCount {
		logger.Warn(fmt.Sprintf("Child %s workflows partially failed. (%d/%d)", childName, errorCount, childCount), "errorCount", errorCount)
		return nil
	}
	err := errors.New("ConsumerSubWorkflowsAllFailedError")
	logger.Error(
		fmt.Sprintf("Child %s workflows all failed , total count is %d", childName,
			errorCount), err)

	return err
}
