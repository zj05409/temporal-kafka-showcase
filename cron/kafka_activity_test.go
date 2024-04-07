package cron

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func TestKafkaProducerConsumer(t *testing.T) {
	testSuite := testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestActivityEnvironment().SetTestTimeout(120 * time.Second)
	env.RegisterActivity(ProducerActivity)
	env.RegisterActivity(ConsumerActivity)

	_, err := env.ExecuteActivity(ProducerActivity)
	require.NoError(t, err)
	_, err = env.ExecuteActivity(ConsumerActivity)
	require.NoError(t, err)

}

func TestKafkaConsumerProducer(t *testing.T) {
	testSuite := testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestActivityEnvironment()

	env.RegisterActivity(ProducerActivity)
	env.RegisterActivity(ConsumerActivity)

	_, err := env.ExecuteActivity(ConsumerActivity)
	require.NoError(t, err)
	_, err = env.ExecuteActivity(ProducerActivity)
	require.NoError(t, err)
}
