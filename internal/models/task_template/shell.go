package task_template

type ShellTaskTemplate struct {
	BaseConfigTaskTemplate
	Command string
}

func (c *ShellTaskTemplate) GetName() string {
	return c.Name
}

func (c *ShellTaskTemplate) GetStoreKey() *string {
	return c.StoreKey
}
