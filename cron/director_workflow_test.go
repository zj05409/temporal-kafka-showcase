package cron

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func TestWholeWorkflows(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}

	env := s.NewTestWorkflowEnvironment()

	env.RegisterWorkflow(ProducersWorkflow)
	env.RegisterWorkflow(ProducerWorkflow)
	env.RegisterActivity(ProducerActivity)

	env.OnActivity(ProducerActivity, mock.Anything).Return(nil)

	env.ExecuteWorkflow(ProducersWorkflow)

	assert.True(t, env.IsWorkflowCompleted())
	error := env.GetWorkflowError()

	require.NoError(t, error)

	env = s.NewTestWorkflowEnvironment()

	env.RegisterWorkflow(ConsumersWorkflow)
	env.RegisterWorkflow(ConsumerWorkflow)
	env.RegisterActivity(ConsumerActivity)
	env.OnActivity(ConsumerActivity, mock.Anything).Return(nil)

	env.ExecuteWorkflow(ConsumersWorkflow)
	assert.True(t, env.IsWorkflowCompleted())

	error = env.GetWorkflowError()

	require.NoError(t, error)
}
