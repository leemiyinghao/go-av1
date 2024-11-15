package coordinator_domain

import (
	"log"
	"sync"
	"time"

	"github.com/leemiyinghao/go-av1/internal/domain/coordinator_domain/runner"
	"github.com/leemiyinghao/go-av1/internal/entity/task_flow"
	"github.com/leemiyinghao/go-av1/internal/value_object/task_execution_type"
)

type DynamicSemaphore struct {
	operation_lock sync.Mutex
	capacity       int
	current        int
}

func NewDynamicSemaphore() *DynamicSemaphore {
	return &DynamicSemaphore{
		operation_lock: sync.Mutex{},
		capacity:       0,
		current:        0,
	}
}

func (s *DynamicSemaphore) AddCapacity(capacity int) {
	s.operation_lock.Lock()
	defer s.operation_lock.Unlock()
	s.capacity += capacity
}

func (s *DynamicSemaphore) TryAcquire() bool {
	s.operation_lock.Lock()
	defer s.operation_lock.Unlock()
	if s.current < s.capacity {
		s.current++
		return true
	}
	return false
}

func (s *DynamicSemaphore) Release() {
	s.operation_lock.Lock()
	defer s.operation_lock.Unlock()
	s.current--
}

type Coordinator struct {
	CPURunnerSets []*RunnerSet
	GPURunnerSets []*RunnerSet
}

func NewCoordinator() *Coordinator {
	return &Coordinator{
		CPURunnerSets: []*RunnerSet{},
		GPURunnerSets: []*RunnerSet{},
	}
}

func (c *Coordinator) AddCPURunner(runner runner.Runner, capacity int) {
	runnerSet := RunnerSet{runnerLock: NewDynamicSemaphore(), runner: runner}
	runnerSet.runnerLock.AddCapacity(capacity)
	c.CPURunnerSets = append(c.CPURunnerSets, &runnerSet)
}
func (c *Coordinator) AddGPURunner(runner runner.Runner, capacity int) {
	runnerSet := RunnerSet{runnerLock: NewDynamicSemaphore(), runner: runner}
	runnerSet.runnerLock.AddCapacity(capacity)
	c.GPURunnerSets = append(c.GPURunnerSets, &runnerSet)
}

func (c *Coordinator) ExecuteTaskFlows(tfs []*task_flow.TaskFlow) {
	tfQueue := NewSimpleQueue()
	for _, tf := range tfs {
		tfQueue.Push(tf)
	}
	for tf := tfQueue.Pop(); tf != nil; tf = tfQueue.Pop() {
		tf := tf.(*task_flow.TaskFlow)
		if tf.IsFinish() {
			continue
		}
		tfQueue.Push(tf)
		if !tf.TryLock(){
			continue
		}

		var runnerSets []*RunnerSet
		switch tf.GetCurrentTask().GetExecutionType() {
		case task_execution_type.CPU:
			runnerSets = c.CPURunnerSets
		case task_execution_type.GPU:
			runnerSets = c.GPURunnerSets
		}

		for _, runnerSet := range runnerSets {
			if runnerSet.runnerLock.TryAcquire() {
				go executeTaskFlow(runnerSet, tf)
				break
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func executeTaskFlow(runnerSet *RunnerSet, taskFlow *task_flow.TaskFlow) {
	defer runnerSet.runnerLock.Release()
	defer taskFlow.Unlock()
	defer taskFlow.Next()
	task := taskFlow.GetCurrentTask()
	if filter := task.GetFilter(); filter != nil && !eval(*filter, *taskFlow.GetContext()) {
		return
	}
	output, new_file, err := runnerSet.runner.Run(task, taskFlow.GetFile())
	if err != nil {
		log.Printf("Runtime Error: %v, %v", err, task)
	} else if new_file != nil {
		taskFlow.SetFile(*new_file)
	}
	if key := task.GetStoreKey(); output != nil && key != nil {
		taskFlow.SetContext(*key, *output)
	}
}

type RunnerSet struct {
	runnerLock *DynamicSemaphore
	runner     runner.Runner
}
