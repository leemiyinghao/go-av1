package filterable_task_template

import (
	"log"

	"github.com/maja42/goval"
)

type FilterableTaskTemplate struct {
	Filter *string `default:"nil"`
}

func NewFilterableTaskTemplate(filter *string) FilterableTaskTemplate {
	return FilterableTaskTemplate{Filter: filter}
}

func (c FilterableTaskTemplate) MatchFilter(context map[string]interface{}) bool {
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
