package config_source

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	var configSource ConfigSource
	// Read the file and unmarshal the data
	t.Run("TestLoad", func(t *testing.T) {
		var err error
		configSource, err = LoadYaml("yaml_config_source_test.yaml")
		assert.Nil(t, err)
	})

	var actions []ConfigTask
	t.Run("TestGetActions", func(t *testing.T) {
		// Get the actions
		actions = configSource.GetActions()
		assert.Equal(t, 2, len(actions))
	})

	var condictions = []struct {
		action       ConfigTask
		context      map[string]interface{}
		name         string
		command      string
		filterResult bool
		storeKey     *string
	}{
		{
			action:       actions[0],
			context:      map[string]interface{}{},
			name:         "Test Task 1",
			command:      "echo \"Hello World\"",
			filterResult: true,
			storeKey:     nil,
		},
		{
			action:       actions[1],
			context:      map[string]interface{}{},
			name:         "Test Task 2",
			command:      "echo \"Hello World 2\"",
			filterResult: false,
			storeKey:     nil,
		},
		{
			action:       actions[1],
			context:      map[string]interface{}{"some_store": "hi"},
			name:         "Test Task 2",
			command:      "echo \"Hello World 2\"",
			filterResult: true,
			storeKey:     nil,
		},
	}
	for _, condiction := range condictions {
		t.Run(fmt.Sprintf("TestActionDetail: %v", condiction), func(t *testing.T) {
			assert.Equal(t, condiction.name, condiction.action.Name)
			assert.Equal(t, condiction.command, condiction.action.Command)
			assert.Equal(t, condiction.filterResult, condiction.action.MatchFilter(condiction.context))
			assert.Equal(t, condiction.storeKey, condiction.action.StoreKey)
		})
	}
}
