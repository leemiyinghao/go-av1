package task_template

type FFmpegTaskTemplate struct {
	BaseConfigTaskTemplate
	InputKwargs  map[string]string
	OutputKwargs map[string]string
}

func (c *FFmpegTaskTemplate) GetName() string {
	return c.Name
}

func (c *FFmpegTaskTemplate) GetStoreKey() *string {
	return c.StoreKey
}
