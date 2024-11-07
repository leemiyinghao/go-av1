package task_template

type TaskTemplate interface {
	MatchFilter(context map[string]interface{}) bool
}
