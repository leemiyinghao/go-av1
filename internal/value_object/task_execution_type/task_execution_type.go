package task_execution_type

type TaskExecutionType int

const (
	CPU TaskExecutionType = iota
	GPU
)
