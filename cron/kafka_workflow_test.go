package cron

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func TestProducerAndConsumerWorkflows(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}

	env := s.NewTestWorkflowEnvironment()

	env.RegisterWorkflow(ProducerWorkflow)
	env.RegisterActivity(ProducerActivity)

	env.OnActivity(ProducerActivity, mock.Anything).Return(nil)

	env.ExecuteWorkflow(ProducerWorkflow)

	assert.True(t, env.IsWorkflowCompleted())
	// result := interface{}(nil)
	// env.GetWorkflowResult(&result)
	assert.True(t, env.IsWorkflowCompleted())
	error := env.GetWorkflowError()

	require.NoError(t, error)

	env = s.NewTestWorkflowEnvironment()

	env.RegisterWorkflow(ConsumerWorkflow)
	env.RegisterActivity(ConsumerActivity)
	env.OnActivity(ConsumerActivity, mock.Anything).Return(nil)

	env.ExecuteWorkflow(ConsumerWorkflow)
	assert.True(t, env.IsWorkflowCompleted())

	error = env.GetWorkflowError()

	require.NoError(t, error)
}
