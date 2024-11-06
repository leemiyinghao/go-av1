package distributor

import (
	"sync"

	"github.com/leemiyinghao/go-av1/internal/models/task"
	"github.com/leemiyinghao/go-av1/internal/runner"
)

type Distributor struct {
	runners      []runner.Runner
	running      sync.WaitGroup
	input_queue  chan task.Task
	output_queue chan task.Task
}

func NewDistributor(runners []runner.Runner, size int) *Distributor {
	distributor := Distributor{
		runners:      runners,
		running:      sync.WaitGroup{},
		input_queue:  make(chan task.Task, size),
		output_queue: make(chan task.Task, size),
	}
	return &distributor
}

func (d *Distributor) Start() {
	for _, r := range d.runners {
		r.SetSource(d.input_queue)
		r.Start()
	}
	go d.rotateResult()
}

func (d *Distributor) addTask(t task.Task) {
	d.running.Add(1)
	d.input_queue <- t
}

func (d *Distributor) receiveResult() (t task.Task) {
	defer d.running.Done()
	t = <-d.output_queue
	return
}

func (d *Distributor) rotateResult() {
	for {
		current := d.receiveResult()
		var template task.TaskTemplate
		var next task.Task
		if template = current.GetTemplate(); template == nil {
			return
		}
		if next = template.GenerateNext(); next == nil {
			return
		}
		d.addTask(next)
	}
}

func (d *Distributor) AddTaskTemplate(tt task.TaskTemplate) {
	task := tt.GenerateNext()
	if task != nil {
		d.addTask(task)
	}
}

func (d *Distributor) Wait() {
	d.running.Wait()
}
