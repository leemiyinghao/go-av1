package config_source

import (
	"log"

	"github.com/maja42/goval"
)

type ConfigTask struct {
	Name     string
	Command  string
	Filter   *string `default:"nil"`
	StoreKey *string `default:"nil"`
}

func (c *ConfigTask) MatchFilter(context map[string]interface{}) bool {
	if c.Filter == nil {
		return true
	}
	
	evaluator := goval.NewEvaluator()
	result, err := evaluator.Evaluate(*c.Filter, context, nil)
	if err != nil {
		log.Printf("Error evaluating filter: %v", err)
		return false
	}
	log.Printf("Filter result: %v", result)
	return result == true
}
