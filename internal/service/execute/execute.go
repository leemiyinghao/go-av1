package execute

import (
	"github.com/leemiyinghao/go-av1/internal/domain/config_domain"
	"github.com/leemiyinghao/go-av1/internal/domain/coordinator_domain"
	"github.com/leemiyinghao/go-av1/internal/domain/coordinator_domain/runner/local_runner"
	"github.com/leemiyinghao/go-av1/internal/domain/file_scan_domain"
	"github.com/leemiyinghao/go-av1/internal/entity/task_flow"
)

const (
	CPURunnerCount, GPURunnerCount = 1, 2
)

type WalkingExecutor struct {
	coordinator      *coordinator_domain.Coordinator
	files            []string
	taskFlowTemplate task_flow.TaskFlowTemplate
}

func NewWalkingExecutor(rootPath string, config_repository config_domain.ConfigRepository) *WalkingExecutor {
	coordinator := coordinator_domain.NewCoordinator()
	coordinator.AddCPURunner(local_runner.NewLocalRunner(), 1)
	coordinator.AddGPURunner(local_runner.NewLocalRunner(), 2)

	taskTemplates, err := config_repository.GetTaskFlowTemplate(rootPath)
	if err != nil {
		panic(err)
	}
	return &WalkingExecutor{
		coordinator:      coordinator,
		files:            file_scan_domain.ScanFiles(rootPath),
		taskFlowTemplate: taskTemplates,
	}
}

func (w *WalkingExecutor) Execute() {
	taskFlows := make([]*task_flow.TaskFlow, 0, len(w.files))
	for _, file := range w.files {
		taskFlows = append(taskFlows, task_flow.NewTaskFlow(w.taskFlowTemplate, file))
	}
	w.coordinator.ExecuteTaskFlows(taskFlows)
}
