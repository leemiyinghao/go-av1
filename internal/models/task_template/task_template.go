package task_template

type TaskTemplate interface {
	MatchFilter(context map[string]interface{}) bool
	GetName() string
	GetStoreKey() *string
}
