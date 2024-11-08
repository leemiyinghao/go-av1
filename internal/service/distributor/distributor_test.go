package distributor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/leemiyinghao/go-av1/internal/runner"
	task_mocks "github.com/leemiyinghao/go-av1/mocks/internal_/models/task"
	runner_mocks "github.com/leemiyinghao/go-av1/mocks/internal_/runner"
)

func TestNewDistributor(t *testing.T) {
	runners := []runner.Runner{
		&runner_mocks.MockRunner{},
	}
	d := NewDistributor(runners, 1024)
	assert.NotNil(t, d)
	assert.Equal(t, cap(d.input_queue), 1024)
	assert.Equal(t, cap(d.output_queue), 1024)
	assert.Equal(t, runners, d.runners)
}

func TestDistributor_AddTask(t *testing.T) {
	mock_runner := runner_mocks.MockRunner{}
	mock_runner.On("SetSource", mock.Anything)
	mock_runner.On("Start")
	var runners []runner.Runner = []runner.Runner{
		&mock_runner,
	}
	task := &task_mocks.MockTask{}
	taskTemplate := &task_mocks.MockTaskTemplate{}

	taskTemplate.On("GenerateNext").Return(task).Once()
	taskTemplate.On("GenerateNext").Return(nil)
	task.On("GetTemplate").Return(taskTemplate).Once()

	d := NewDistributor(runners, 1024)
	d.Start()
	d.AddTaskFlow(taskTemplate)
	d.output_queue <- task
	d.Wait()

	task.AssertCalled(t, "GetTemplate")
	mock_runner.AssertCalled(t, "SetSource", d.input_queue)
	mock_runner.AssertCalled(t, "Start")
	taskTemplate.AssertCalled(t, "GenerateNext")
}
