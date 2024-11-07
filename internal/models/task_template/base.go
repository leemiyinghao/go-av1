package task_template

import (
	"log"

	"github.com/maja42/goval"

	"github.com/leemiyinghao/go-av1/internal/models/execution_type"
)

type BaseConfigTaskTemplate struct {
	Name          string
	ExecutionType execution_type.ExecutionType
	Filter        *string `default:"nil"`
	StoreKey      *string `default:"nil"`
}

func (c BaseConfigTaskTemplate) MatchFilter(context map[string]interface{}) bool {
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
