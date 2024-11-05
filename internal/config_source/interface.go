package config_source

type ConfigSource interface {
	GetActions() []ConfigTask
}
