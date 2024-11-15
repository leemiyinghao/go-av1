package coordinator_domain

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"

	"github.com/leemiyinghao/go-av1/internal/entity/task_flow"
	"github.com/leemiyinghao/go-av1/internal/value_object/task_execution_type"
	runner_mocks "github.com/leemiyinghao/go-av1/mocks/internal_/domain/coordinator_domain/runner"
)

func TestDynamicSemaphore(t *testing.T) {

	var ds *DynamicSemaphore
	t.Run("InitialDynamicSemaphore", func(t *testing.T) {
		ds = NewDynamicSemaphore()
		assert.Equal(t, 0, ds.capacity)
	})

	t.Run("AddCapacity", func(t *testing.T) {

		ds.AddCapacity(10)
		assert.Equal(t, 10, ds.capacity)
	})

	for i := 1; i < 11; i++ {
		t.Run(fmt.Sprintf("TryAcquire: %d", i), func(t *testing.T) {
			assert.True(t, ds.TryAcquire())
		})
	}

	t.Run("TryAcquire: 11", func(t *testing.T) {
		assert.False(t, ds.TryAcquire())
	})

	t.Run("Release", func(t *testing.T) {
		ds.Release()
		assert.Equal(t, 9, ds.current)
	})
}

func TestNewSimpleQueue(t *testing.T) {
	q := NewSimpleQueue()
	assert.NotNil(t, q)
}

func TestSimpleQueuePushPop(t *testing.T) {
	q := NewSimpleQueue()
	q.Push(1)
	q.Push(2)
	popped := q.Pop()
	assert.Equal(t, 1, popped)

	popped = q.Pop()
	assert.Equal(t, 2, popped)

	// Try popping from an empty queue
	popped = q.Pop()
	assert.Nil(t, popped)
}

func TestSimpleQueueLen(t *testing.T) {
	q := NewSimpleQueue()
	assert.Equal(t, 0, q.Len())

	q.Push(1)
	q.Push(2)
	assert.Equal(t, 2, q.Len())
}

func TestCoordinator(t *testing.T) {

	c := NewCoordinator()

	cpuRunner := runner_mocks.NewMockRunner(t)
	gpuRunner := runner_mocks.NewMockRunner(t)

	t.Run("AddCPURunner", func(t *testing.T) {
		c.AddCPURunner(cpuRunner, 10)
		assert.Len(t, c.CPURunnerSets, 1)
	})

	t.Run("AddGPURunner", func(t *testing.T) {
		c.AddGPURunner(gpuRunner, 10)
		assert.Len(t, c.GPURunnerSets, 1)
	})

	t.Run("ExecuteTaskFlows", func(t *testing.T) {
		task1 := task_flow.NewShellTask(
			"task 1",
			task_execution_type.CPU,
			ptr.To("1==1"),
			ptr.To("task_1"),
			"command 1",
		)
		task2 := task_flow.NewShellTask(
			"task 2",
			task_execution_type.GPU,
			ptr.To("2==1"),
			ptr.To("task_2"),
			"command 2",
		)
		tf := task_flow.NewTaskFlow([]task_flow.Task{
			task1,
			task2,
		}, "original_file_name")

		cpuRunner.On("Run", task1, "original_file_name").Return(ptr.To("0"), ptr.To("new_file_name"), nil)

		c.ExecuteTaskFlows([]*task_flow.TaskFlow{tf})

		cpuRunner.AssertCalled(t, "Run", task1, "original_file_name")
		cpuRunner.AssertNotCalled(t, "Run", task2, "new_file_name")
		gpuRunner.AssertNotCalled(t, "Run", task1, "original_file_name")
		gpuRunner.AssertNotCalled(t, "Run", task2, "new_file_name")
		assert.Equal(t, 0, tf.GetRemainingTaskCount())
		assert.Equal(t, "new_file_name", tf.GetFile())
	})
}
