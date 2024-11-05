package config_source

type ConfigSource interface {
	GetTasks() []ConfigTask
}
